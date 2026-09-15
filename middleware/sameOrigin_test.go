package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func sameOriginRecord(t *testing.T, method, host string, headers map[string]string,
	behindProxy bool, panelDomain string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(SameOrigin(behindProxy, panelDomain))
	handle := func(c *gin.Context) { c.Status(http.StatusOK) }
	engine.GET("/x", handle)
	engine.POST("/x", handle)

	req := httptest.NewRequest(method, "http://"+host+"/x", nil)
	req.Host = host
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

// sameOriginStatus is the directly-exposed panel: nothing in front, so the Host
// the request carries is the browser's own statement of where it went.
func sameOriginStatus(t *testing.T, method, host string, headers map[string]string) int {
	t.Helper()
	return sameOriginRecord(t, method, host, headers, false, "").Code
}

func TestSameOrigin(t *testing.T) {
	const host = "panel.example.com"

	tests := []struct {
		name    string
		method  string
		headers map[string]string
		want    int
	}{
		// A read cannot change anything, and the panel is also fetched by
		// clients that send no Origin at all.
		{"GET is never blocked", http.MethodGet, nil, http.StatusOK},
		{"GET from elsewhere is not blocked either", http.MethodGet,
			map[string]string{"Origin": "https://evil.example"}, http.StatusOK},

		{"same origin over https", http.MethodPost,
			map[string]string{"Origin": "https://" + host}, http.StatusOK},
		// The panel runs on plain HTTP too, and behind a TLS-terminating proxy
		// the scheme the browser used is not the one we see. Only the host is
		// compared.
		{"same host over http", http.MethodPost,
			map[string]string{"Origin": "http://" + host}, http.StatusOK},

		{"another origin", http.MethodPost,
			map[string]string{"Origin": "https://evil.example"}, http.StatusForbidden},
		// A subdomain is a different host, and the cookie is not scoped to it.
		{"a subdomain is not the same host", http.MethodPost,
			map[string]string{"Origin": "https://evil." + host}, http.StatusForbidden},
		{"a host that merely starts the same", http.MethodPost,
			map[string]string{"Origin": "https://" + host + ".evil.example"}, http.StatusForbidden},

		// Referer is the fallback when Origin is absent.
		{"referer from the panel", http.MethodPost,
			map[string]string{"Referer": "https://" + host + "/app/"}, http.StatusOK},
		{"referer from elsewhere", http.MethodPost,
			map[string]string{"Referer": "https://evil.example/page"}, http.StatusForbidden},

		// Neither header: a cross-site form post cannot set a custom one, so
		// the panel's own XHR marker tells them apart.
		{"no headers but the XHR marker", http.MethodPost,
			map[string]string{"X-Requested-With": "XMLHttpRequest"}, http.StatusOK},
		{"no headers at all", http.MethodPost, nil, http.StatusForbidden},
		{"an unparseable origin", http.MethodPost,
			map[string]string{"Origin": "::not a url::"}, http.StatusForbidden},
		// "null" is what a sandboxed iframe or a data: document sends.
		{"a null origin", http.MethodPost,
			map[string]string{"Origin": "null"}, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sameOriginStatus(t, tt.method, host, tt.headers); got != tt.want {
				t.Errorf("status = %d, want %d", got, tt.want)
			}
		})
	}
}

// The request Host and the Origin do not always carry the same port, and the
// panel is reachable in all of these shapes:
//
//   - directly on its own port, where the browser dialled exactly what the
//     panel sees;
//   - behind the nginx vhost this panel generates, which forwards $host --
//     no port -- while the browser's Origin carries the one it dialled;
//   - behind a proxy that forwards $http_host, where the two match again.
//
// So the hostname is what is compared. Comparing host:port whole rejected
// every write on a proxied panel reachable on anything but 443, which is a
// configuration that works today -- DomainValidator strips the port before its
// own comparison for the same reason.
func TestSameOriginIgnoresThePort(t *testing.T) {
	tests := []struct {
		name   string
		host   string // what the panel sees as Host
		origin string // what the browser sends
		want   int
	}{
		{"direct, same port on both", "1.2.3.4:2095", "http://1.2.3.4:2095", http.StatusOK},
		{"nginx $host, browser on 8443", "panel.example.com", "https://panel.example.com:8443", http.StatusOK},
		{"nginx $http_host, both carry the port", "panel.example.com:8443", "https://panel.example.com:8443", http.StatusOK},
		{"nginx $host, browser on 443", "panel.example.com", "https://panel.example.com", http.StatusOK},
		{"IPv6 literal", "[::1]:2095", "http://[::1]:2095", http.StatusOK},
		{"IPv6 literal, port only on one side", "[::1]", "http://[::1]:2095", http.StatusOK},

		// Dropping the port must not start accepting a different host.
		{"another host on the same port", "panel.example.com:8443", "https://evil.example:8443", http.StatusForbidden},
		{"a subdomain, ports aside", "panel.example.com", "https://evil.panel.example.com:8443", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sameOriginStatus(t, http.MethodPost, tt.host, map[string]string{"Origin": tt.origin})
			if got != tt.want {
				t.Errorf("Host %q, Origin %q -> %d, want %d", tt.host, tt.origin, got, tt.want)
			}
		})
	}
}

// Behind a reverse proxy, c.Request.Host is whatever that proxy chose to send.
// nginx's own default is `proxy_set_header Host $proxy_host`, i.e. the address
// the vhost dials the panel on, which names nothing a browser ever typed --
// comparing an Origin against it refuses every write on a panel that works,
// and the login POST with them, so there is nowhere left to fix it from.
func TestSameOriginBehindAProxy(t *testing.T) {
	// What a vhost without `proxy_set_header Host $host;` leaves behind.
	const rewritten = "127.0.0.1:2095"
	const public = "panel.example.com"

	tests := []struct {
		name        string
		host        string
		headers     map[string]string
		behindProxy bool
		panelDomain string
		want        int
	}{
		{"the proxy forwards the real host", rewritten,
			map[string]string{"Origin": "https://" + public, "X-Forwarded-Host": public},
			true, "", http.StatusOK},
		{"and still refuses a foreign origin", rewritten,
			map[string]string{"Origin": "https://evil.example", "X-Forwarded-Host": public},
			true, "", http.StatusForbidden},
		// A chain of proxies appends; any host the request actually passed
		// through is a host the browser may have dialled.
		{"a forwarded chain", rewritten,
			map[string]string{"Origin": "https://" + public, "X-Forwarded-Host": public + ", edge.example.com"},
			true, "", http.StatusOK},
		// The port belongs to the browser, not to the forwarded name.
		{"forwarded host without the browser's port", rewritten,
			map[string]string{"Origin": "https://" + public + ":8443", "X-Forwarded-Host": public},
			true, "", http.StatusOK},

		// The configured domain answers the same question, and is what an
		// operator whose proxy forwards nothing can set.
		{"the panel domain stands in", rewritten,
			map[string]string{"Origin": "https://" + public}, true, public, http.StatusOK},
		{"the panel domain still refuses a foreign origin", rewritten,
			map[string]string{"Origin": "https://evil.example"}, true, public, http.StatusForbidden},

		// Neither: the panel does not know its own name, so it stands down
		// rather than locking the operator out. SameSite on the cookie is what
		// guards the request in that configuration.
		{"a proxy that forwards nothing", rewritten,
			map[string]string{"Origin": "https://" + public}, true, "", http.StatusOK},

		// X-Forwarded-Host is client input unless something in front
		// overwrote it, which is the same rule getRemoteIp applies to
		// X-Forwarded-For.
		{"forwarded host is ignored without the proxy setting", public,
			map[string]string{"Origin": "https://evil.example", "X-Forwarded-Host": "evil.example"},
			false, "", http.StatusForbidden},

		// A correctly configured proxy sends Host and nothing else; the domain
		// is what says the panel is reachable there.
		{"proxy forwards Host, domain configured", public,
			map[string]string{"Origin": "https://" + public}, true, public, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sameOriginRecord(t, http.MethodPost, tt.host, tt.headers,
				tt.behindProxy, tt.panelDomain).Code
			if got != tt.want {
				t.Errorf("Host %q, headers %v, behindProxy=%v, domain %q -> %d, want %d",
					tt.host, tt.headers, tt.behindProxy, tt.panelDomain, got, tt.want)
			}
		})
	}
}

// A refusal has to say which two names failed to match. An empty 403 reaches
// the operator as "Request failed with status code 403", which names neither --
// and httputil only reads a body carrying all three of success, msg and obj.
func TestSameOriginRefusalCarriesTheMsgShape(t *testing.T) {
	rec := sameOriginRecord(t, http.MethodPost, "panel.example.com",
		map[string]string{"Origin": "https://evil.example"}, false, "")
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body %q is not JSON: %v", rec.Body.String(), err)
	}
	for _, key := range []string{"success", "msg", "obj"} {
		if _, ok := body[key]; !ok {
			t.Errorf("body is missing %q: %v", key, body)
		}
	}
	if body["success"] != false {
		t.Errorf("success = %v, want false", body["success"])
	}
	msg, _ := body["msg"].(string)
	if !strings.Contains(msg, "evil.example") || !strings.Contains(msg, "panel.example.com") {
		t.Errorf("msg = %q, want both the origin and the expected host named", msg)
	}
}
