package util

import "testing"

func TestStripServerTlsFields(t *testing.T) {
	tests := []struct {
		name    string
		in      map[string]interface{}
		changed bool
		gone    []string
		kept    []string
	}{
		{
			name: "top-level server fields removed",
			in: map[string]interface{}{
				"enabled":          true,
				"server_name":      "example.com",
				"alpn":             []string{"h2"},
				"certificate_path": "/root/cert/example.com/fullchain.pem",
				"key":              []string{"-----BEGIN PRIVATE KEY-----"},
				"key_path":         "/root/cert/example.com/private.key",
				"acme":             map[string]interface{}{"domain": []string{"example.com"}},
			},
			changed: true,
			gone:    []string{"certificate_path", "key", "key_path", "acme"},
			kept:    []string{"enabled", "server_name", "alpn"},
		},
		{
			name: "ech private key removed, client ech fields kept",
			in: map[string]interface{}{
				"enabled": true,
				"ech": map[string]interface{}{
					"enabled":     true,
					"key":         []string{"ech-key"},
					"key_path":    "/root/ech.key",
					"config":      []string{"ech-config"},
					"config_path": "/etc/ech.cfg",
				},
			},
			changed: true,
			kept:    []string{"enabled", "ech"},
		},
		{
			name: "clean map untouched",
			in: map[string]interface{}{
				"enabled":     true,
				"server_name": "example.com",
				"certificate": []string{"-----BEGIN CERTIFICATE-----"},
				"utls":        map[string]interface{}{"enabled": true, "fingerprint": "safari"},
			},
			changed: false,
			kept:    []string{"enabled", "server_name", "certificate", "utls"},
		},
		{
			name:    "nil map is a no-op",
			in:      nil,
			changed: false,
		},
		{
			name: "non-map ech value ignored",
			in: map[string]interface{}{
				"ech": "not-a-map",
			},
			changed: false,
			kept:    []string{"ech"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripServerTlsFields(tt.in); got != tt.changed {
				t.Fatalf("changed = %v, want %v", got, tt.changed)
			}
			for _, k := range tt.gone {
				if _, ok := tt.in[k]; ok {
					t.Errorf("field %q should have been removed", k)
				}
			}
			for _, k := range tt.kept {
				if _, ok := tt.in[k]; !ok {
					t.Errorf("field %q should have been kept", k)
				}
			}
		})
	}

	t.Run("nested ech key removed but ech object survives", func(t *testing.T) {
		in := map[string]interface{}{
			"ech": map[string]interface{}{
				"enabled": true,
				"key":     "secret",
				"config":  []string{"cfg"},
			},
		}
		if !StripServerTlsFields(in) {
			t.Fatal("expected changed = true")
		}
		ech := in["ech"].(map[string]interface{})
		if _, ok := ech["key"]; ok {
			t.Error("ech.key should have been removed")
		}
		if _, ok := ech["config"]; !ok {
			t.Error("ech.config should have been kept")
		}
	})

	t.Run("idempotent", func(t *testing.T) {
		in := map[string]interface{}{"certificate_path": "/x"}
		if !StripServerTlsFields(in) {
			t.Fatal("first call should report a change")
		}
		if StripServerTlsFields(in) {
			t.Fatal("second call should be a no-op")
		}
	})
}
