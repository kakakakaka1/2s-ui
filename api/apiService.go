package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/shenaba/2s-ui/config"
	"github.com/shenaba/2s-ui/database"
	"github.com/shenaba/2s-ui/logger"
	"github.com/shenaba/2s-ui/service"
	"github.com/shenaba/2s-ui/util"
	"github.com/shenaba/2s-ui/util/common"

	"github.com/gin-gonic/gin"
)

type ApiService struct {
	service.SettingService
	service.UserService
	service.ConfigService
	service.ClientService
	service.TlsService
	service.InboundService
	service.OutboundService
	service.EndpointService
	service.ServicesService
	service.PanelService
	service.StatsService
	service.ServerService
	service.UpdateService
	service.NodeService
	service.NodeSyncService
}

func (a *ApiService) UpdateInfo(c *gin.Context) {
	canUpdate, reason := a.UpdateService.CanSelfUpdate()
	jsonObj(c, map[string]interface{}{
		"canSelfUpdate": canUpdate,
		"reason":        reason,
		"current":       config.GetVersion(),
		// Docker updates live in the container's writable layer; the UI warns
		// that recreating the container reverts to the image's version.
		"docker": a.UpdateService.InDocker(),
	}, nil)
}

func (a *ApiService) UpdatePanel(c *gin.Context) {
	err := a.UpdateService.StartUpdate()
	jsonMsg(c, "updatePanel", err)
}

func (a *ApiService) UpdateStatus(c *gin.Context) {
	jsonObj(c, a.UpdateService.GetStatus(), nil)
}

func (a *ApiService) LoadData(c *gin.Context) {
	data, err := a.getData(c)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, data, nil)
}

func (a *ApiService) getData(c *gin.Context) (interface{}, error) {
	data := make(map[string]interface{}, 0)
	lu := c.Query("lu")
	isUpdated, err := a.ConfigService.CheckChanges(lu)
	if err != nil {
		return "", err
	}
	onlines, err := a.StatsService.GetOnlines()

	sysInfo := a.ServerService.GetSingboxInfo()
	if sysInfo["running"] == false {
		logs := a.ServerService.GetLogs("1", "debug")
		if len(logs) > 0 {
			data["lastLog"] = logs[0]
		}
	}

	if err != nil {
		return "", err
	}
	if isUpdated {
		config, err := a.SettingService.GetConfig()
		if err != nil {
			return "", err
		}
		clients, err := a.ClientService.GetAll()
		if err != nil {
			return "", err
		}
		tlsConfigs, err := a.TlsService.GetAll()
		if err != nil {
			return "", err
		}
		inbounds, err := a.InboundService.GetAll()
		if err != nil {
			return "", err
		}
		outbounds, err := a.OutboundService.GetAll()
		if err != nil {
			return "", err
		}
		endpoints, err := a.EndpointService.GetAll()
		if err != nil {
			return "", err
		}
		services, err := a.ServicesService.GetAll()
		if err != nil {
			return "", err
		}
		subURI, err := a.SettingService.GetFinalSubURI(getHostname(c))
		if err != nil {
			return "", err
		}
		trafficAge, err := a.SettingService.GetTrafficAge()
		if err != nil {
			return "", err
		}
		nodes, err := a.NodeService.GetAll()
		if err != nil {
			return "", err
		}
		data["config"] = json.RawMessage(config)
		data["clients"] = clients
		data["tls"] = tlsConfigs
		data["inbounds"] = inbounds
		data["outbounds"] = outbounds
		data["endpoints"] = endpoints
		data["services"] = services
		data["nodes"] = nodes
		data["subURI"] = subURI
		data["enableTraffic"] = trafficAge > 0
		data["onlines"] = onlines
	} else {
		data["onlines"] = onlines
	}

	// Live node status rides along unconditionally (like onlines): it changes
	// every heartbeat, so it must not hide behind the lu gate. Omitted when
	// empty so zero-node setups pay nothing.
	nodesStatus := a.NodeService.GetStatuses()
	if len(nodesStatus) > 0 {
		data["nodesStatus"] = nodesStatus
	}

	return data, nil
}

func (a *ApiService) LoadPartialData(c *gin.Context, objs []string) error {
	data := make(map[string]interface{}, 0)
	id := c.Query("id")

	for _, obj := range objs {
		switch obj {
		case "inbounds":
			inbounds, err := a.InboundService.Get(id)
			if err != nil {
				return err
			}
			data[obj] = inbounds
		case "outbounds":
			outbounds, err := a.OutboundService.GetAll()
			if err != nil {
				return err
			}
			data[obj] = outbounds
		case "endpoints":
			endpoints, err := a.EndpointService.GetAll()
			if err != nil {
				return err
			}
			data[obj] = endpoints
		case "services":
			services, err := a.ServicesService.GetAll()
			if err != nil {
				return err
			}
			data[obj] = services
		case "tls":
			tlsConfigs, err := a.TlsService.GetAll()
			if err != nil {
				return err
			}
			data[obj] = tlsConfigs
		case "clients":
			clients, err := a.ClientService.Get(id)
			if err != nil {
				return err
			}
			data[obj] = clients
		case "config":
			config, err := a.SettingService.GetConfig()
			if err != nil {
				return err
			}
			data[obj] = json.RawMessage(config)
		case "settings":
			settings, err := a.SettingService.GetAllSetting()
			if err != nil {
				return err
			}
			data[obj] = settings
		case "nodes":
			nodes, err := a.NodeService.GetAll()
			if err != nil {
				return err
			}
			data[obj] = nodes
		}
	}

	jsonObj(c, data, nil)
	return nil
}

func (a *ApiService) GetUsers(c *gin.Context) {
	users, err := a.UserService.GetUsers()
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, *users, nil)
}

func (a *ApiService) GetSettings(c *gin.Context) {
	data, err := a.SettingService.GetAllSetting()
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, data, err)
}

func (a *ApiService) GetStats(c *gin.Context) {
	resource := c.Query("resource")
	tag := c.Query("tag")
	period := c.Query("period")
	if period == "" {
		period = "hour"
	}
	data, err := a.StatsService.GetStats(resource, tag, period)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, data, err)
}

func (a *ApiService) GetStatus(c *gin.Context) {
	request := c.Query("r")
	result := a.ServerService.GetStatus(request)
	jsonObj(c, result, nil)
}

func (a *ApiService) GetOnlines(c *gin.Context) {
	onlines, err := a.StatsService.GetOnlines()
	jsonObj(c, onlines, err)
}

func (a *ApiService) GetLogs(c *gin.Context) {
	count := c.Query("c")
	level := c.Query("l")
	logs := a.ServerService.GetLogs(count, level)
	jsonObj(c, logs, nil)
}

func (a *ApiService) CheckChanges(c *gin.Context) {
	actor := c.Query("a")
	chngKey := c.Query("k")
	count := c.Query("c")
	changes := a.ConfigService.GetChanges(actor, chngKey, count)
	jsonObj(c, changes, nil)
}

func (a *ApiService) GetKeypairs(c *gin.Context) {
	kType := c.Query("k")
	options := c.Query("o")
	keypair := a.ServerService.GenKeypair(kType, options)
	jsonObj(c, keypair, nil)
}

func (a *ApiService) GetDb(c *gin.Context) {
	exclude := c.Query("exclude")
	db, err := database.GetDb(exclude)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=s-ui_"+time.Now().Format("20060102-150405")+".db")
	c.Writer.Write(db)
}

func (a *ApiService) postActions(c *gin.Context) (string, json.RawMessage, error) {
	var data map[string]json.RawMessage
	err := c.ShouldBind(&data)
	if err != nil {
		return "", nil, err
	}
	return string(data["action"]), data["data"], nil
}

func (a *ApiService) Login(c *gin.Context) {
	remoteIP := getRemoteIp(c)
	loginUser, err := a.UserService.Login(c.Request.FormValue("user"), c.Request.FormValue("pass"), remoteIP)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}

	sessionMaxAge, err := a.SettingService.GetSessionMaxAge()
	if err != nil {
		logger.Infof("Unable to get session's max age from DB")
	}

	err = SetLoginUser(c, loginUser, sessionMaxAge)
	if err == nil {
		logger.Info("user ", loginUser, " login success")
	} else {
		logger.Warning("login failed: ", err)
	}

	jsonMsg(c, "", nil)
}

func (a *ApiService) ChangePass(c *gin.Context) {
	id := c.Request.FormValue("id")
	oldPass := c.Request.FormValue("oldPass")
	newUsername := c.Request.FormValue("newUsername")
	newPass := c.Request.FormValue("newPass")
	err := a.UserService.ChangePass(id, oldPass, newUsername, newPass)
	if err == nil {
		logger.Info("change user credentials success")
		jsonMsg(c, "save", nil)
	} else {
		logger.Warning("change user credentials failed:", err)
		jsonMsg(c, "", err)
	}
}

// Save handles POST api/save. fanout controls whether a successful client /
// inbound change is propagated to managed nodes. v1 (the SPA) passes true; v2
// passes false so a node that has this panel as a master does not bounce the
// pushed change back out (ping-pong between mutual masters).
func (a *ApiService) Save(c *gin.Context, loginUser string, fanout bool) {
	hostname := getHostname(c)
	obj := c.Request.FormValue("object")
	act := c.Request.FormValue("action")
	data := c.Request.FormValue("data")
	initUsers := c.Request.FormValue("initUsers")
	objs, err := a.ConfigService.Save(obj, act, json.RawMessage(data), initUsers, loginUser, hostname)
	if err != nil {
		jsonMsg(c, "save", err)
		return
	}
	if fanout && (obj == "clients" || obj == "inbounds") {
		// Fire-and-forget: network IO must not block the save's DB txn (already
		// committed) or the HTTP response. Unrelated nodes are a cheap no-op diff.
		a.NodeSyncService.MarkAllDirty()
		go a.NodeSyncService.ReconcileDirtyOnline()
	}
	err = a.LoadPartialData(c, objs)
	if err != nil {
		jsonMsg(c, obj, err)
	}
}

func (a *ApiService) RestartApp(c *gin.Context) {
	err := a.PanelService.RestartPanel(3)
	jsonMsg(c, "restartApp", err)
}

func (a *ApiService) RestartSb(c *gin.Context) {
	err := a.ConfigService.RestartCore()
	jsonMsg(c, "restartSb", err)
}

func (a *ApiService) ResetTraffic(c *gin.Context) {
	if err := a.ClientService.ResetAllClientsTraffic(); err != nil {
		jsonMsg(c, "resetTraffic", err)
		return
	}
	err := a.ConfigService.RestartCore()
	jsonMsg(c, "resetTraffic", err)
}

func (a *ApiService) LinkConvert(c *gin.Context) {
	link := c.Request.FormValue("link")
	result, _, err := util.GetOutbound(link, 0)
	jsonObj(c, result, err)
}

func (a *ApiService) SubConvert(c *gin.Context) {
	link := c.Request.FormValue("link")
	result, err := util.GetExternalSub(link)
	jsonObj(c, result, err)
}

func (a *ApiService) ImportDb(c *gin.Context) {
	file, _, err := c.Request.FormFile("db")
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	defer file.Close()
	err = database.ImportDB(file)
	jsonMsg(c, "", err)
}

func (a *ApiService) Logout(c *gin.Context) {
	loginUser := GetLoginUser(c)
	if loginUser != "" {
		logger.Infof("user %s logout", loginUser)
	}
	ClearSession(c)
	jsonMsg(c, "", nil)
}

func (a *ApiService) LoadTokens() ([]byte, error) {
	return a.UserService.LoadTokens()
}

func (a *ApiService) GetTokens(c *gin.Context) {
	loginUser := GetLoginUser(c)
	tokens, err := a.UserService.GetUserTokens(loginUser)
	jsonObj(c, tokens, err)
}

func (a *ApiService) AddToken(c *gin.Context) {
	loginUser := GetLoginUser(c)
	expiry := c.Request.FormValue("expiry")
	expiryInt, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	desc := c.Request.FormValue("desc")
	token, err := a.UserService.AddToken(loginUser, expiryInt, desc)
	jsonObj(c, token, err)
}

func (a *ApiService) DeleteToken(c *gin.Context) {
	tokenId := c.Request.FormValue("id")
	err := a.UserService.DeleteToken(tokenId)
	jsonMsg(c, "", err)
}

func (a *ApiService) GetSingboxConfig(c *gin.Context) {
	rawConfig, err := a.ConfigService.GetConfig("")
	if err != nil {
		c.Status(400)
		c.Writer.WriteString(err.Error())
		return
	}
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=config_"+time.Now().Format("20060102-150405")+".json")
	c.Writer.Write(*rawConfig)
}

func (a *ApiService) TestAcme(c *gin.Context) {
	domain := c.Request.FormValue("domain")
	email := c.Request.FormValue("email")
	if err := a.ConfigService.TestAcme(domain, email); err != nil {
		pureJsonMsg(c, false, err.Error())
		return
	}
	pureJsonMsg(c, true, "")
}

func (a *ApiService) DetectNginx(c *gin.Context) {
	var acme service.AcmeService
	jsonObj(c, acme.DetectNginx(), nil)
}

func (a *ApiService) IssueCert(c *gin.Context) {
	domain := c.Request.FormValue("domain")
	email := c.Request.FormValue("email")
	useNginx := c.Request.FormValue("nginx") == "true"
	force := c.Request.FormValue("force") == "true"
	var acme service.AcmeService
	res, err := acme.IssueWeb(domain, email, useNginx, force)
	if err != nil {
		pureJsonMsg(c, false, err.Error())
		return
	}
	jsonMsgObj(c, "", res, nil)
}

func (a *ApiService) GetCheckOutbound(c *gin.Context) {
	tag := c.Query("tag")
	link := c.Query("link")
	result := a.ConfigService.CheckOutbound(tag, link)
	jsonObj(c, result, nil)
}

func (a *ApiService) GetCertPing(c *gin.Context) {
	domain := c.PostForm("domain")
	port := c.PostForm("port")
	tlsPing, err := util.GetTlsPing(domain, port)
	jsonObj(c, tlsPing, err)
}

func (a *ApiService) TestNode(c *gin.Context) {
	data := c.PostForm("data")
	status, err := a.NodeService.TestNode(json.RawMessage(data))
	jsonObj(c, status, err)
}

func (a *ApiService) GetNodeInbounds(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		jsonMsg(c, "nodeInbounds", common.NewError("invalid node id"))
		return
	}
	inbounds, err := a.NodeSyncService.FetchNodeInbounds(uint(id))
	jsonObj(c, inbounds, err)
}

func (a *ApiService) AdoptInbounds(c *gin.Context, loginUser string) {
	id, err := strconv.ParseUint(c.PostForm("id"), 10, 64)
	if err != nil {
		jsonMsg(c, "adoptInbounds", common.NewError("invalid node id"))
		return
	}
	var tags []string
	if err := json.Unmarshal([]byte(c.PostForm("tags")), &tags); err != nil {
		jsonMsg(c, "adoptInbounds", common.NewError("invalid tags"))
		return
	}
	if err := a.NodeSyncService.AdoptInbounds(uint(id), tags, loginUser); err != nil {
		jsonMsg(c, "adoptInbounds", err)
		return
	}
	// Push the master's clients onto the freshly adopted inbounds right away.
	// ReconcileNow skips the backoff; if it still loses to an in-flight run,
	// the dirty flag set by AdoptInbounds lets the heartbeat converge instead.
	go func() {
		if err := a.NodeSyncService.ReconcileNow(uint(id)); err != nil {
			logger.Warning("adopt: initial reconcile failed: ", err)
		}
	}()
	jsonMsg(c, "adoptInbounds", nil)
}

func (a *ApiService) ReconcileNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.PostForm("id"), 10, 64)
	if err != nil {
		jsonMsg(c, "reconcileNode", common.NewError("invalid node id"))
		return
	}
	err = a.NodeSyncService.ReconcileNow(uint(id))
	jsonMsg(c, "reconcileNode", err)
}
