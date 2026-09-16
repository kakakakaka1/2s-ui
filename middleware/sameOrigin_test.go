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
	opt Options) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(SameOrigin(opt))
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
	return sameOriginRecord(t, method, host, headers, Options{}).Code
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
//
// But standing down for every proxied panel is the other failure: most vhosts
// do carry `proxy_set_header Host $host;`, and those were checked correctly all
// along. So the stand-down has to be earned -- the Host has to actually name
// the socket this panel bound.
func TestSameOriginBehindAProxy(t *testing.T) {
	// What a vhost without `proxy_set_header Host $host;` leaves behind.
	const rewritten = "127.0.0.1:2095"
	const public = "panel.example.com"
	proxied := func(o Options) Options {
		o.BehindProxy = true
		if o.Port == 0 {
			o.Port = 2095
		}
		return o
	}

	tests := []struct {
		name    string
		host    string
		headers map[string]string
		opt     Options
		want    int
	}{
		{"the proxy forwards the real host", rewritten,
			map[string]string{"Origin": "https://" + public, "X-Forwarded-Host": public},
			proxied(Options{}), http.StatusOK},
		{"and still refuses a foreign origin", rewritten,
			map[string]string{"Origin": "https://evil.example", "X-Forwarded-Host": public},
			proxied(Options{}), http.StatusForbidden},
		// A chain of proxies appends; any host the request actually passed
		// through is a host the browser may have dialled.
		{"a forwarded chain", rewritten,
			map[string]string{"Origin": "https://" + public, "X-Forwarded-Host": public + ", edge.example.com"},
			proxied(Options{}), http.StatusOK},
		// The port belongs to the browser, not to the forwarded name.
		{"forwarded host without the browser's port", rewritten,
			map[string]string{"Origin": "https://" + public + ":8443", "X-Forwarded-Host": public},
			proxied(Options{}), http.StatusOK},

		// A configured domain settles it by itself. (DomainValidator would have
		// aborted a rewritten Host before this runs, so in practice this is the
		// proxy that forwards Host correctly.)
		{"the panel domain stands in", public,
			map[string]string{"Origin": "https://" + public},
			proxied(Options{PanelDomain: public}), http.StatusOK},
		{"the panel domain still refuses a foreign origin", public,
			map[string]string{"Origin": "https://evil.example"},
			proxied(Options{PanelDomain: public}), http.StatusForbidden},

		// The regression this function exists for: a vhost that forwards Host
		// and nothing else, with no domain configured. The panel can see that
		// Host is not its own socket, so the mismatch is real.
		{"proxy forwards Host, no domain, foreign origin", public,
			map[string]string{"Origin": "https://evil.example"},
			proxied(Options{}), http.StatusForbidden},
		{"proxy forwards Host, no domain, own origin", public,
			map[string]string{"Origin": "https://" + public},
			proxied(Options{}), http.StatusOK},

		// Only here does the check stand down: every name the panel has is its
		// own socket, so none of them says where the browser went.
		{"a proxy that forwards nothing", rewritten,
			map[string]string{"Origin": "https://" + public},
			proxied(Options{}), http.StatusOK},
		{"a proxy that forwards nothing, foreign origin", rewritten,
			map[string]string{"Origin": "https://evil.example"},
			proxied(Options{}), http.StatusOK},
		// ...and an operator tunnelling straight to that socket matches it
		// outright, without needing the stand-down at all.
		{"browsing the socket directly", rewritten,
			map[string]string{"Origin": "http://" + rewritten},
			proxied(Options{}), http.StatusOK},

		// The spellings are not interchangeable as strings, and both sides of
		// each of these is correctly configured -- refusing here locks the
		// operator out of a panel that works.
		{"bound to loopback, proxy_pass localhost", "localhost:2095",
			map[string]string{"Origin": "https://" + public},
			proxied(Options{Listen: "127.0.0.1"}), http.StatusOK},
		{"bound to ::1, proxy_pass 127.0.0.1", rewritten,
			map[string]string{"Origin": "https://" + public},
			proxied(Options{Listen: "::1"}), http.StatusOK},
		// A vhost that sets X-Forwarded-Host from $proxy_host by mistake: the
		// header is there but says only what Host already said.
		{"forwarded host is the upstream address too", rewritten,
			map[string]string{"Origin": "https://" + public, "X-Forwarded-Host": rewritten},
			proxied(Options{}), http.StatusOK},
		// One candidate that is not our socket is enough to make the answer
		// real again, even when Host is.
		{"forwarded host names the public name", rewritten,
			map[string]string{"Origin": "https://evil.example", "X-Forwarded-Host": public},
			proxied(Options{}), http.StatusForbidden},

		// X-Forwarded-Host is client input unless something in front
		// overwrote it, which is the same rule getRemoteIp applies to
		// X-Forwarded-For.
		{"forwarded host is ignored without the proxy setting", public,
			map[string]string{"Origin": "https://evil.example", "X-Forwarded-Host": "evil.example"},
			Options{}, http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sameOriginRecord(t, http.MethodPost, tt.host, tt.headers, tt.opt).Code
			if got != tt.want {
				t.Errorf("Host %q, headers %v, opt %+v -> %d, want %d",
					tt.host, tt.headers, tt.opt, got, tt.want)
			}
		})
	}
}

// $proxy_host expands to whatever the vhost dialled, and that is the only Host
// the panel is allowed to call uninformative. Anything else had to come from
// something that knew the public name.
func TestHostIsOwnSocket(t *testing.T) {
	tests := []struct {
		name   string
		host   string
		listen string
		port   int
		want   bool
	}{
		{"loopback on our port", "127.0.0.1:2095", "", 2095, true},
		{"localhost on our port", "localhost:2095", "", 2095, true},
		{"IPv6 loopback on our port", "[::1]:2095", "", 2095, true},
		{"loopback while bound to every interface", "127.0.0.1:2095", "0.0.0.0", 2095, true},
		{"loopback while bound to every v6 interface", "127.0.0.1:2095", "::", 2095, true},
		// Bound to one address: that is the address a proxy elsewhere dials,
		// and loopback may not even be bound.
		{"the bound address", "10.0.0.5:2095", "10.0.0.5", 2095, true},
		// Loopback counts whatever webListen says: the spellings are not
		// interchangeable as strings, and a proxy on the same host dials one of
		// them. Refusing here is what locked the panel out.
		{"loopback while bound elsewhere", "127.0.0.1:2095", "10.0.0.5", 2095, true},
		{"localhost while bound to loopback", "localhost:2095", "127.0.0.1", 2095, true},
		{"loopback while bound to the other family", "127.0.0.1:2095", "::1", 2095, true},

		// A real public name, which is what a forwarded Host looks like.
		{"a hostname", "panel.example.com:2095", "", 2095, false},
		{"a hostname with no port", "panel.example.com", "", 2095, false},
		// nginx's `Host $host` for a browser on 443 carries no port at all.
		{"loopback with no port", "127.0.0.1", "", 2095, false},
		{"loopback on another port", "127.0.0.1:8443", "", 2095, false},
		// Port 0 is the "could not read the setting" default: nothing matches
		// it, so the check never stands down on a bad read.
		{"unknown port", "127.0.0.1:2095", "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hostIsOwnSocket(tt.host, tt.listen, tt.port); got != tt.want {
				t.Errorf("hostIsOwnSocket(%q, %q, %d) = %v, want %v",
					tt.host, tt.listen, tt.port, got, tt.want)
			}
		})
	}
}

// A refusal has to say which two names failed to match. An empty 403 reaches
// the operator as "Request failed with status code 403", which names neither --
// and httputil only reads a body carrying all three of success, msg and obj.
func TestSameOriginRefusalCarriesTheMsgShape(t *testing.T) {
	rec := sameOriginRecord(t, http.MethodPost, "panel.example.com",
		map[string]string{"Origin": "https://evil.example"}, Options{})
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
