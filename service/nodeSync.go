package service

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shenaba/2s-ui/database"
	"github.com/shenaba/2s-ui/database/model"
	"github.com/shenaba/2s-ui/logger"
	"github.com/shenaba/2s-ui/util/common"

	"gorm.io/gorm"
)

// NodeSyncService owns everything that writes TO a node: adopting its inbounds
// as read-only replicas, and pushing/reconciling the master's clients onto it.
// Reconcile is the single write channel — first push, offline catch-up, drift
// repair and the manual button all funnel through it.
type NodeSyncService struct {
	NodeService
}

const (
	// clusterGroup marks clients this master pushed to a node. Reconciliation
	// is scoped to this group so it never touches a node's own local users and
	// needs no tombstones for deletes.
	clusterGroup = "@cluster"

	nodePushTimeout = 15 * time.Second // node hot-restarts inbounds on save
	reconcileBackoff = 30 * time.Second
)

// per-node single-flight + backoff so overlapping triggers (heartbeat, fanout,
// manual button) don't stampede the same node.
var (
	reconcileMu   sync.Mutex
	reconcileBusy = map[uint]bool{}
	reconcileLast = map[uint]time.Time{}
)

// ---------- remote calls ----------

// nodePost sends an x-www-form-urlencoded apiv2 request (apiv2 save reads
// c.Request.FormValue, not a JSON body) and unwraps the {success,msg,obj} envelope.
func (s *NodeSyncService) nodePost(n *model.Node, client *http.Client, action string, form url.Values) (json.RawMessage, error) {
	u := nodeApiURL(n, action)
	req, err := http.NewRequest("POST", u, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Token", n.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, common.NewErrorf("HTTP %d from node panel", resp.StatusCode)
	}
	var msg struct {
		Success bool            `json:"success"`
		Msg     string          `json:"msg"`
		Obj     json.RawMessage `json:"obj"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, nodeMaxResponseSize)).Decode(&msg); err != nil {
		return nil, common.NewError("unexpected response from node")
	}
	if !msg.Success {
		if msg.Msg == "" {
			msg.Msg = "node rejected the request"
		}
		return nil, common.NewError(msg.Msg)
	}
	return msg.Obj, nil
}

// pushClient runs one clients save (new|edit|del) against a node.
func (s *NodeSyncService) pushClient(n *model.Node, client *http.Client, act string, payload interface{}) (json.RawMessage, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Set("object", "clients")
	form.Set("action", act)
	form.Set("data", string(data))
	return s.nodePost(n, client, "save", form)
}

// ---------- inbound adoption ----------

type remoteInbound struct {
	Id      uint   `json:"id"`
	Type    string `json:"type"`
	Tag     string `json:"tag"`
	Adopted bool   `json:"adopted"`
}

// nodeClient builds a short-lived HTTP client honouring the node's TLS mode,
// with the longer push timeout.
func nodePushClient(n *model.Node) *http.Client {
	c := buildNodeHTTPClient(n)
	c.Timeout = nodePushTimeout
	return c
}

// FetchNodeInbounds lists a node's inbounds, flagging which tags this panel has
// already adopted as replicas.
func (s *NodeSyncService) FetchNodeInbounds(nodeId uint) ([]remoteInbound, error) {
	node, err := s.getNodeById(nodeId)
	if err != nil {
		return nil, err
	}
	obj, err := s.nodeGet(node, nodePushClient(node), "inbounds", nil)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Inbounds []struct {
			Id   uint   `json:"id"`
			Type string `json:"type"`
			Tag  string `json:"tag"`
		} `json:"inbounds"`
	}
	if err := json.Unmarshal(obj, &payload); err != nil {
		return nil, common.NewError("unexpected inbounds payload from node")
	}
	adopted, err := s.adoptedTags(nodeId)
	if err != nil {
		return nil, err
	}
	out := make([]remoteInbound, 0, len(payload.Inbounds))
	for _, ib := range payload.Inbounds {
		out = append(out, remoteInbound{
			Id: ib.Id, Type: ib.Type, Tag: ib.Tag,
			Adopted: adopted[ib.Tag],
		})
	}
	return out, nil
}

func (s *NodeSyncService) adoptedTags(nodeId uint) (map[string]bool, error) {
	var tags []string
	err := database.GetDB().Model(model.Inbound{}).Where("node_id = ?", nodeId).Pluck("tag", &tags).Error
	if err != nil {
		return nil, err
	}
	m := make(map[string]bool, len(tags))
	for _, t := range tags {
		m[t] = true
	}
	return m, nil
}

// AdoptInbounds pulls the full panel-shape inbound for each selected tag and
// stores it as a local replica row (node_id set). tag collisions fail loudly —
// the tag is the reconciliation key, so we never silently rename.
func (s *NodeSyncService) AdoptInbounds(nodeId uint, tags []string, actor string) error {
	if len(tags) == 0 {
		return nil
	}
	node, err := s.getNodeById(nodeId)
	if err != nil {
		return err
	}
	client := nodePushClient(node)
	obj, err := s.nodeGet(node, client, "inbounds", nil)
	if err != nil {
		return err
	}
	var payload struct {
		Inbounds []json.RawMessage `json:"inbounds"`
	}
	if err := json.Unmarshal(obj, &payload); err != nil {
		return common.NewError("unexpected inbounds payload from node")
	}
	wanted := map[string]bool{}
	for _, t := range tags {
		wanted[t] = true
	}

	db := database.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		dt := time.Now().Unix()
		adopted := 0
		for _, raw := range payload.Inbounds {
			var meta struct {
				Tag string `json:"tag"`
			}
			if err := json.Unmarshal(raw, &meta); err != nil {
				return err
			}
			if !wanted[meta.Tag] {
				continue
			}
			// tag must be globally unique (DB constraint + reconciliation key).
			var count int64
			if err := tx.Model(model.Inbound{}).Where("tag = ?", meta.Tag).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return common.NewErrorf("tag %q already exists here — rename it on the node first", meta.Tag)
			}
			replica, err := buildReplicaInbound(raw, nodeId)
			if err != nil {
				return err
			}
			if err := tx.Create(replica).Error; err != nil {
				return err
			}
			adopted++
		}
		if adopted == 0 {
			return common.NewError("no matching inbounds found on the node")
		}
		if err := tx.Create(&model.Changes{
			DateTime: dt, Actor: actor, Key: "inbounds", Action: "adopt",
			Obj: json.RawMessage(mustJSON(map[string]interface{}{"nodeId": nodeId, "tags": tags})),
		}).Error; err != nil {
			return err
		}
		LastUpdate = time.Now().Unix()
		return nil
	})
}

// buildReplicaInbound turns a node's panel-shape inbound into a local replica:
// keep type/tag/addrs/out_json (links already point at the node), drop tls_id
// (TLS terminates on the node), strip panel-only keys from Options.
func buildReplicaInbound(raw json.RawMessage, nodeId uint) (*model.Inbound, error) {
	var full map[string]interface{}
	if err := json.Unmarshal(raw, &full); err != nil {
		return nil, err
	}
	inb := &model.Inbound{
		Type:   asString(full["type"]),
		Tag:    asString(full["tag"]),
		NodeId: &nodeId,
	}
	if addrs, ok := full["addrs"]; ok && addrs != nil {
		inb.Addrs, _ = json.MarshalIndent(addrs, "", "  ")
	}
	if outJson, ok := full["out_json"]; ok && outJson != nil {
		inb.OutJson, _ = json.MarshalIndent(outJson, "", "  ")
	}
	for _, k := range []string{"id", "tls_id", "tls", "addrs", "out_json", "users", "node_id", "type", "tag"} {
		delete(full, k)
	}
	options, err := json.MarshalIndent(full, "", "  ")
	if err != nil {
		return nil, err
	}
	inb.Options = options
	return inb, nil
}

// ---------- reconciliation ----------

type nodeClientState struct {
	Id       uint            `json:"id"`
	Name     string          `json:"name"`
	Enable   bool            `json:"enable"`
	Config   json.RawMessage `json:"config"`
	Inbounds json.RawMessage `json:"inbounds"`
	Expiry   int64           `json:"expiry"`
	Group    string          `json:"group"`
	Links    json.RawMessage `json:"links"`
	Up       int64           `json:"up"`
	Down     int64           `json:"down"`
}

// Reconcile makes a node's @cluster clients match the master's expectation:
// clients that reference any of this node's replica inbounds. It is the ONLY
// path that writes clients to a node.
func (s *NodeSyncService) Reconcile(nodeId uint) error {
	if !s.claimReconcile(nodeId) {
		return nil
	}
	defer s.releaseReconcile(nodeId)

	node, err := s.getNodeById(nodeId)
	if err != nil {
		return err
	}
	if !node.Enable {
		return nil
	}
	client := nodePushClient(node)

	// tag -> node-local inbound id
	tagToId, err := s.nodeInboundTagMap(node, client)
	if err != nil {
		return err
	}

	expected, err := s.expectedClients(nodeId, tagToId)
	if err != nil {
		return err
	}
	actual, err := s.actualClusterClients(node, client)
	if err != nil {
		return err
	}

	changed := false
	linkUpdates := map[string]json.RawMessage{} // master client name -> node links

	// new / edit
	for name, want := range expected {
		cur, exists := actual[name]
		if !exists {
			obj, err := s.pushClient(node, client, "new", want)
			if err != nil {
				logger.Warning("reconcile: push new ", name, " to node ", node.Name, ": ", err)
				return err
			}
			s.collectNodeLinks(obj, name, linkUpdates)
			changed = true
		} else if clientDiffers(want, cur) {
			want["id"] = cur.Id // node-local id for the edit
			obj, err := s.pushClient(node, client, "edit", want)
			if err != nil {
				logger.Warning("reconcile: push edit ", name, " to node ", node.Name, ": ", err)
				return err
			}
			s.collectNodeLinks(obj, name, linkUpdates)
			changed = true
		}
	}
	// del: on the node, in @cluster, but no longer expected
	for name, cur := range actual {
		if _, ok := expected[name]; !ok {
			if _, err := s.pushClient(node, client, "del", cur.Id); err != nil {
				logger.Warning("reconcile: push del ", name, " to node ", node.Name, ": ", err)
				return err
			}
			changed = true
		}
	}

	if len(linkUpdates) > 0 {
		s.mergeNodeLinks(node, linkUpdates)
	}
	_ = changed

	now := time.Now().Unix()
	return database.GetDB().Model(model.Node{}).Where("id = ?", nodeId).
		Updates(map[string]interface{}{"dirty": false, "last_sync": now}).Error
}

// expectedClients builds the desired @cluster client set for a node: every
// master client that references at least one of the node's replica inbounds,
// shaped as a node-local client (Volume=0, Expiry copied, Group=@cluster,
// Inbounds mapped tag->node id, Links omitted so the node generates them).
func (s *NodeSyncService) expectedClients(nodeId uint, tagToId map[string]uint) (map[string]map[string]interface{}, error) {
	db := database.GetDB()

	// replica inbound id -> tag, for this node
	var replicas []model.Inbound
	if err := db.Model(model.Inbound{}).Where("node_id = ?", nodeId).Find(&replicas).Error; err != nil {
		return nil, err
	}
	replicaTagById := map[uint]string{}
	replicaIds := map[uint]bool{}
	for _, r := range replicas {
		replicaTagById[r.Id] = r.Tag
		replicaIds[r.Id] = true
	}
	if len(replicaIds) == 0 {
		return map[string]map[string]interface{}{}, nil
	}

	var clients []model.Client
	if err := db.Model(model.Client{}).Find(&clients).Error; err != nil {
		return nil, err
	}

	expected := map[string]map[string]interface{}{}
	for i := range clients {
		c := &clients[i]
		var ids []uint
		if err := json.Unmarshal(c.Inbounds, &ids); err != nil {
			continue
		}
		var nodeLocalIds []uint
		for _, id := range ids {
			if tag, ok := replicaTagById[id]; ok {
				if localId, ok := tagToId[tag]; ok {
					nodeLocalIds = append(nodeLocalIds, localId)
				}
			}
		}
		if len(nodeLocalIds) == 0 {
			continue
		}
		inboundsJSON, _ := json.Marshal(nodeLocalIds)
		expected[c.Name] = map[string]interface{}{
			"name":     c.Name,
			"enable":   c.Enable,
			"config":   c.Config,
			"inbounds": json.RawMessage(inboundsJSON),
			// links must be present (even empty): the node's link-refresh path
			// unmarshals it and a nil RawMessage is "unexpected end of JSON input".
			// The node regenerates the actual links from the request Host.
			"links":  json.RawMessage("[]"),
			"volume": 0,        // quota is the master's job only
			"expiry": c.Expiry, // absolute; node self-expires consistently
			"group":  clusterGroup,
			"desc":   c.Desc,
		}
	}
	return expected, nil
}

func (s *NodeSyncService) actualClusterClients(node *model.Node, client *http.Client) (map[string]nodeClientState, error) {
	obj, err := s.nodeGet(node, client, "clients", nil)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Clients []nodeClientState `json:"clients"`
	}
	if err := json.Unmarshal(obj, &payload); err != nil {
		return nil, common.NewError("unexpected clients payload from node")
	}
	out := map[string]nodeClientState{}
	for _, c := range payload.Clients {
		if c.Group == clusterGroup {
			out[c.Name] = c
		}
	}
	return out, nil
}

func (s *NodeSyncService) nodeInboundTagMap(node *model.Node, client *http.Client) (map[string]uint, error) {
	obj, err := s.nodeGet(node, client, "inbounds", nil)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Inbounds []struct {
			Id  uint   `json:"id"`
			Tag string `json:"tag"`
		} `json:"inbounds"`
	}
	if err := json.Unmarshal(obj, &payload); err != nil {
		return nil, common.NewError("unexpected inbounds payload from node")
	}
	m := make(map[string]uint, len(payload.Inbounds))
	for _, ib := range payload.Inbounds {
		m[ib.Tag] = ib.Id
	}
	return m, nil
}

// clientDiffers compares the master's desired client against the node's current
// one on the fields we own. Config is compared structurally to avoid whitespace noise.
func clientDiffers(want map[string]interface{}, cur nodeClientState) bool {
	if asBool(want["enable"]) != cur.Enable {
		return true
	}
	if asInt64(want["expiry"]) != cur.Expiry {
		return true
	}
	if !jsonEqual(want["config"], cur.Config) {
		return true
	}
	if !jsonEqual(want["inbounds"], cur.Inbounds) {
		return true
	}
	return false
}

// collectNodeLinks extracts the node-generated "local" links from a save
// response and stashes them under the master client name.
//
// KNOWN LIMITATION (subscription aggregation): a node's save response comes
// from LoadPartialData -> ClientService.Get, whose column projection omits the
// links field, so this currently collects nothing and node routes do not merge
// into the master subscription. The clients themselves are pushed and usable on
// the node; only the aggregated master subscription is incomplete. Fixing it
// without touching the node side means generating the external link locally from
// the replica inbound's stored out_json (the node-side address snapshot) plus
// the client config, rather than relying on the node to hand its links back.
func (s *NodeSyncService) collectNodeLinks(obj json.RawMessage, name string, out map[string]json.RawMessage) {
	if obj == nil {
		return
	}
	var payload struct {
		Clients []struct {
			Name  string          `json:"name"`
			Links json.RawMessage `json:"links"`
		} `json:"clients"`
	}
	if err := json.Unmarshal(obj, &payload); err != nil {
		return
	}
	for _, c := range payload.Clients {
		if c.Name == name && c.Links != nil {
			out[name] = c.Links
		}
	}
}

// mergeNodeLinks folds a node's freshly generated links into the master
// clients' Links as type:"external" entries prefixed "[node] ", replacing any
// previous entries with the same prefix (idempotent). Existing subscription
// output then serves them with zero extra work.
func (s *NodeSyncService) mergeNodeLinks(node *model.Node, updates map[string]json.RawMessage) {
	prefix := "[" + node.Name + "] "
	db := database.GetDB()
	for name, nodeLinksRaw := range updates {
		var nodeLinks []map[string]string
		if err := json.Unmarshal(nodeLinksRaw, &nodeLinks); err != nil {
			continue
		}
		var client model.Client
		if err := db.Model(model.Client{}).Where("name = ?", name).First(&client).Error; err != nil {
			continue
		}
		var existing []map[string]string
		json.Unmarshal(client.Links, &existing)

		merged := make([]map[string]string, 0, len(existing))
		for _, l := range existing {
			if !strings.HasPrefix(l["remark"], prefix) {
				merged = append(merged, l)
			}
		}
		for _, l := range nodeLinks {
			if l["type"] != "local" {
				continue
			}
			merged = append(merged, map[string]string{
				"remark": prefix + l["remark"],
				"type":   "external",
				"uri":    l["uri"],
			})
		}
		links, err := json.MarshalIndent(merged, "", "  ")
		if err != nil {
			continue
		}
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(model.Client{}).Where("id = ?", client.Id).Update("links", links).Error; err != nil {
				return err
			}
			return tx.Create(&model.Changes{
				DateTime: time.Now().Unix(), Actor: "NodeSync", Key: "clients", Action: "edit",
				Obj: json.RawMessage(mustJSON(map[string]interface{}{"name": name, "node": node.Name})),
			}).Error
		}); err != nil {
			logger.Warning("reconcile: merge links for ", name, ": ", err)
			continue
		}
		LastUpdate = time.Now().Unix()
	}
}

// ---------- dirty tracking / triggers ----------

func (s *NodeSyncService) MarkAllDirty() {
	if err := database.GetDB().Model(model.Node{}).Where("enable = ?", true).Update("dirty", true).Error; err != nil {
		logger.Warning("nodes: mark all dirty: ", err)
	}
}

func (s *NodeSyncService) MarkDirty(nodeId uint) {
	database.GetDB().Model(model.Node{}).Where("id = ?", nodeId).Update("dirty", true)
}

// ReconcileDirtyOnline reconciles every enabled node that is online and dirty.
// Called by the heartbeat so offline-period edits converge once a node returns.
func (s *NodeSyncService) ReconcileDirtyOnline() {
	var nodes []model.Node
	if err := database.GetDB().Model(model.Node{}).Where("enable = ? AND dirty = ?", true, true).Find(&nodes).Error; err != nil {
		return
	}
	statuses := s.GetStatuses()
	for i := range nodes {
		st, ok := statuses[nodes[i].Id]
		if !ok || st.State != "online" {
			continue
		}
		id := nodes[i].Id
		go func() {
			if err := s.Reconcile(id); err != nil {
				logger.Warning("nodes: background reconcile failed: ", err)
			}
		}()
	}
}

// ReconcileAllOnline reconciles all online nodes regardless of dirty flag —
// the hourly safety net that repairs silent node-side drift.
func (s *NodeSyncService) ReconcileAllOnline() {
	var nodes []model.Node
	if err := database.GetDB().Model(model.Node{}).Where("enable = ?", true).Find(&nodes).Error; err != nil {
		return
	}
	statuses := s.GetStatuses()
	for i := range nodes {
		st, ok := statuses[nodes[i].Id]
		if !ok || st.State != "online" {
			continue
		}
		if err := s.Reconcile(nodes[i].Id); err != nil {
			logger.Warning("nodes: safety-net reconcile failed: ", err)
		}
	}
}

// ---------- traffic collection ----------

type trafficBaseline struct {
	Up   int64 `json:"up"`
	Down int64 `json:"down"`
}

// CollectTraffic pulls each online node's @cluster client counters and folds
// the delta since the last collection into the master's per-client totals.
// The node's clients.up/down are cumulative; a per-(node,client) baseline turns
// them into deltas, resetting the baseline when the node's counter drops (reset).
func (s *NodeSyncService) CollectTraffic() {
	var nodes []model.Node
	if err := database.GetDB().Model(model.Node{}).Where("enable = ?", true).Find(&nodes).Error; err != nil {
		return
	}
	statuses := s.GetStatuses()
	for i := range nodes {
		node := &nodes[i]
		st, ok := statuses[node.Id]
		if !ok || st.State != "online" {
			continue
		}
		if err := s.collectNodeTraffic(node); err != nil {
			logger.Warning("nodes: collect traffic from ", node.Name, ": ", err)
		}
	}
}

func (s *NodeSyncService) collectNodeTraffic(node *model.Node) error {
	client := nodePushClient(node)
	current, err := s.actualClusterClients(node, client)
	if err != nil {
		return err
	}

	// existing baseline
	baseline := map[string]trafficBaseline{}
	if node.Baselines != nil {
		json.Unmarshal(node.Baselines, &baseline)
	}

	// master client names (only fold traffic for clients we actually own)
	masterNames := map[string]bool{}
	var names []string
	db := database.GetDB()
	db.Model(model.Client{}).Pluck("name", &names)
	for _, n := range names {
		masterNames[n] = true
	}

	newBaseline := map[string]trafficBaseline{}
	type delta struct{ up, down int64 }
	deltas := map[string]delta{}
	for name, cur := range current {
		curUp := clientUp(cur)
		curDown := clientDown(cur)
		newBaseline[name] = trafficBaseline{Up: curUp, Down: curDown}
		if !masterNames[name] {
			continue
		}
		base := baseline[name]
		du := curUp - base.Up
		if du < 0 { // node counter reset — rebase
			du = curUp
		}
		dd := curDown - base.Down
		if dd < 0 {
			dd = curDown
		}
		if du > 0 || dd > 0 {
			deltas[name] = delta{du, dd}
		}
	}

	baselineJSON, err := json.Marshal(newBaseline)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for name, d := range deltas {
			if err := tx.Model(model.Client{}).Where("name = ?", name).
				Updates(map[string]interface{}{
					"up":   gorm.Expr("up + ?", d.up),
					"down": gorm.Expr("down + ?", d.down),
				}).Error; err != nil {
				return err
			}
		}
		return tx.Model(model.Node{}).Where("id = ?", node.Id).Update("baselines", baselineJSON).Error
	})
}

// nodeClientState only carries name/enable/config/inbounds/expiry/group/links;
// up/down come from a second projection of the same clients GET payload.
func clientUp(c nodeClientState) int64   { return c.Up }
func clientDown(c nodeClientState) int64 { return c.Down }

// ---------- helpers ----------

func (s *NodeSyncService) getNodeById(id uint) (*model.Node, error) {
	var node model.Node
	if err := database.GetDB().Model(model.Node{}).Where("id = ?", id).First(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *NodeSyncService) claimReconcile(nodeId uint) bool {
	reconcileMu.Lock()
	defer reconcileMu.Unlock()
	if reconcileBusy[nodeId] {
		return false
	}
	if last, ok := reconcileLast[nodeId]; ok && time.Since(last) < reconcileBackoff {
		return false
	}
	reconcileBusy[nodeId] = true
	return true
}

func (s *NodeSyncService) releaseReconcile(nodeId uint) {
	reconcileMu.Lock()
	defer reconcileMu.Unlock()
	reconcileBusy[nodeId] = false
	reconcileLast[nodeId] = time.Now()
}

func asString(v interface{}) string {
	s, _ := v.(string)
	return s
}

func asBool(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func asInt64(v interface{}) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	case json.Number:
		i, _ := n.Int64()
		return i
	case string:
		i, _ := strconv.ParseInt(n, 10, 64)
		return i
	}
	return 0
}

// jsonEqual compares two JSON values structurally (ignoring key order / whitespace).
func jsonEqual(a interface{}, b json.RawMessage) bool {
	var ar json.RawMessage
	switch v := a.(type) {
	case json.RawMessage:
		ar = v
	case []byte:
		ar = v
	default:
		m, err := json.Marshal(a)
		if err != nil {
			return false
		}
		ar = m
	}
	var av, bv interface{}
	if err := json.Unmarshal(ar, &av); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &bv); err != nil {
		return false
	}
	am, _ := json.Marshal(canonical(av))
	bm, _ := json.Marshal(canonical(bv))
	return string(am) == string(bm)
}

// canonical recursively normalises maps so Marshal is order-stable.
func canonical(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, val := range t {
			out[k] = canonical(val)
		}
		return out
	case []interface{}:
		for i := range t {
			t[i] = canonical(t[i])
		}
		return t
	}
	return v
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
