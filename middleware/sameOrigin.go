package middleware

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/shenaba/2s-ui/logger"

	"github.com/gin-gonic/gin"
)

// SameOrigin rejects state-changing requests that did not come from the panel's
// own pages.
//
// The session is a cookie, so without this any page the operator visits while
// logged in can post to the panel in their name -- add a client, change the
// password, import a database. SameSite on the cookie already covers current
// browsers; this is the half that does not depend on the browser being current,
// and it is what answers a request that arrives with an Origin from somewhere
// else.
//
// Only the hostname is compared -- never a hardcoded scheme, and never the
// port. The panel has to work on plain HTTP as well as HTTPS, and behind a
// TLS-terminating proxy the scheme the browser used and the one the panel sees
// differ anyway, so requiring https here would reject every legitimate request
// in two of the three deployments. The port is dropped for the same kind of
// reason: nginx's $host, which the vhost this panel generates forwards as
// Host, carries no port, while the browser's Origin carries the one it
// actually dialled -- so comparing them whole rejects every write on a proxied
// panel reachable on anything but 443. DomainValidator already strips the port
// before its own comparison, for the same deployment.
//
// behindProxy is the webNginx setting and panelDomain is webDomain; both decide
// what the panel is allowed to believe about its own public name. See
// expectedHosts -- getting that wrong refuses every write on a working install.
//
// Mounted on the cookie-authenticated group only. apiv2 authenticates with a
// Token header, which a cross-site page cannot set without a CORS preflight the
// panel never answers.
func SameOrigin(behindProxy bool, panelDomain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = c.GetHeader("Referer")
		}
		if origin == "" {
			// Neither header. A cross-site form post cannot set a custom one,
			// so the panel's own XHR marker tells the two apart -- and it is
			// already what checkLogin answers on.
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.Next()
				return
			}
			rejectCrossSite(c, "", nil)
			return
		}

		u, err := url.Parse(origin)
		if err != nil || u.Host == "" {
			rejectCrossSite(c, origin, nil)
			return
		}

		expected, known := expectedHosts(c, behindProxy, panelDomain)
		if !known {
			warnProxyHidesHost(c.Request.Host)
			c.Next()
			return
		}
		from := hostname(u.Host)
		for _, host := range expected {
			if strings.EqualFold(from, host) {
				c.Next()
				return
			}
		}
		rejectCrossSite(c, from, expected)
	}
}

// expectedHosts returns the hostnames a browser may legitimately have dialled
// to reach this panel, and whether the answer is worth anything.
//
// Directly exposed, the Host header is the browser's own statement of where it
// went, so it is the whole answer. Behind a reverse proxy it is whatever that
// proxy chose to send, and nginx's default is `proxy_set_header Host
// $proxy_host` -- the panel's own listen address, which names nothing a browser
// ever typed. Comparing an Origin against that refuses every write on a panel
// that works perfectly, with the operator locked out of the only interface they
// have: the login POST is refused too, so there is nowhere to go and fix it
// from.
//
// Two things can still say what the public name is behind a proxy, so both are
// accepted: the configured domain, and X-Forwarded-Host. The header is read
// only with the reverse-proxy setting on -- the same rule getRemoteIp applies to
// X-Forwarded-For, and for the same reason. It is safe here regardless of which
// entry of a chain is taken, because it is not a CORS-safelisted header: a
// cross-site page cannot set one without a preflight this panel never answers,
// so a request carrying it did not come from the attack this guards against.
//
// When neither says anything, known is false and the caller stands down. That
// is the honest answer -- the panel does not know its own name -- and standing
// down leaves SameSite on the cookie, which browsers have enforced for years,
// as the guard rather than bricking the install. It does not cost the panel's
// own deployment anything: EnsureVhost refuses to generate a vhost without a
// domain, so webNginx being on with webDomain empty means a proxy this panel
// did not write and knows nothing about.
func expectedHosts(c *gin.Context, behindProxy bool, panelDomain string) (hosts []string, known bool) {
	if panelDomain != "" {
		hosts = append(hosts, hostname(panelDomain))
	}
	if !behindProxy {
		return append(hosts, hostname(c.Request.Host)), true
	}
	for _, forwarded := range strings.Split(c.GetHeader("X-Forwarded-Host"), ",") {
		if forwarded = strings.TrimSpace(forwarded); forwarded != "" {
			hosts = append(hosts, hostname(forwarded))
		}
	}
	return hosts, len(hosts) > 0
}

// proxyHostWarned keeps the stand-down to one line in the log. It is reached on
// every write of every session on a panel whose proxy forwards no host, and the
// buffer GetLogs reads holds ten thousand lines in total.
var proxyHostWarned sync.Once

func warnProxyHidesHost(host string) {
	proxyHostWarned.Do(func() {
		logger.Warning("the reverse proxy in front of this panel forwards neither",
			" the browser's Host nor X-Forwarded-Host (this request arrived as \"", host,
			"\"), so the same-origin check cannot run. Add",
			" `proxy_set_header Host $host;` to the vhost, or set the panel domain.")
	})
}

// rejectCrossSite answers in the panel's own {success, msg, obj} shape rather
// than with a bare status. The body is what reaches the operator: an empty 403
// surfaces in the UI as "Request failed with status code 403", which says
// nothing about which two names failed to match or what to do about it.
func rejectCrossSite(c *gin.Context, from string, expected []string) {
	msg := "cross-site request refused"
	if from != "" {
		msg += ": this request says it came from \"" + from + "\""
		if len(expected) > 0 {
			msg += ", but the panel is reachable as \"" + strings.Join(expected, "\", \"") + "\""
		}
	}
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"success": false,
		"msg":     msg,
		"obj":     nil,
	})
}

// hostname drops the port from a host:port, leaving anything without one --
// and an IPv6 literal's brackets -- as it is.
func hostname(host string) string {
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return strings.Trim(host, "[]")
}
