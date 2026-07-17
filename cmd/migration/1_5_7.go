package migration

import (
	"encoding/json"

	"github.com/shenaba/2s-ui/util"

	"gorm.io/gorm"
)

// to1_5_7 strips server-only TLS fields (certificate_path, key, key_path,
// acme) that older builds let leak into client-facing TLS JSON: the tls
// table's client column, inbound out_json and per-address tls overrides
// (issue #51). Unlike to1_5_1 this cleans every row, not just self-signed
// ones. Idempotent.
func to1_5_7(tx *gorm.DB) error {
	type tlsRow struct {
		Id     uint
		Client []byte
	}
	var tlsRows []tlsRow
	if err := tx.Raw("SELECT id, client FROM tls").Scan(&tlsRows).Error; err != nil {
		return err
	}
	for _, row := range tlsRows {
		if len(row.Client) == 0 {
			continue
		}
		var client map[string]interface{}
		if err := json.Unmarshal(row.Client, &client); err != nil {
			continue
		}
		if !util.StripServerTlsFields(client) {
			continue
		}
		newClient, err := json.MarshalIndent(client, "", "  ")
		if err != nil {
			return err
		}
		if err := tx.Exec("UPDATE tls SET client = ? WHERE id = ?", newClient, row.Id).Error; err != nil {
			return err
		}
	}

	type inboundRow struct {
		Id      uint
		OutJson []byte
		Addrs   []byte
	}
	var inbounds []inboundRow
	if err := tx.Raw("SELECT id, out_json, addrs FROM inbounds").Scan(&inbounds).Error; err != nil {
		return err
	}
	for _, in := range inbounds {
		if len(in.OutJson) > 0 {
			var out map[string]interface{}
			if err := json.Unmarshal(in.OutJson, &out); err == nil {
				if tlsM, ok := out["tls"].(map[string]interface{}); ok && util.StripServerTlsFields(tlsM) {
					newOut, err := json.MarshalIndent(out, "", "  ")
					if err != nil {
						return err
					}
					if err := tx.Exec("UPDATE inbounds SET out_json = ? WHERE id = ?", newOut, in.Id).Error; err != nil {
						return err
					}
				}
			}
		}
		if len(in.Addrs) > 0 {
			var addrs []map[string]interface{}
			if err := json.Unmarshal(in.Addrs, &addrs); err == nil {
				changed := false
				for _, addr := range addrs {
					if addrTls, ok := addr["tls"].(map[string]interface{}); ok && util.StripServerTlsFields(addrTls) {
						changed = true
					}
				}
				if changed {
					newAddrs, err := json.MarshalIndent(addrs, "", "  ")
					if err != nil {
						return err
					}
					if err := tx.Exec("UPDATE inbounds SET addrs = ? WHERE id = ?", newAddrs, in.Id).Error; err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
