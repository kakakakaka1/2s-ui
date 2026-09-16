package middleware

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
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
// Options carries what the panel knows about its own public identity. All of it
// is read once when the router is built, the way DomainValidator is: changing
// any of these already needs a panel restart, since that is what rebuilds the
// router.
type Options struct {
	// BehindProxy is the webNginx setting: something in front wrote the Host
	// header, so it is not by itself the browser's statement of where it went.
	BehindProxy bool
	// PanelDomain is webDomain, "" when unset.
	PanelDomain string
	// Listen and Port are webListen and webPort -- the socket the panel bound.
	// They are what tells a forwarded Host from the panel's own address; see
	// hostIsOwnSocket.
	Listen string
	Port   int
}

// Mounted on the cookie-authenticated group only. apiv2 authenticates with a
// Token header, which a cross-site page cannot set without a CORS preflight the
// panel never answers.
func SameOrigin(opt Options) gin.HandlerFunc {
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

		// Match first, and against every name the panel could answer to --
		// c.Request.Host included, whatever wrote it. A proxy that forwards the
		// browser's Host makes this the whole check, and one that rewrote it to
		// the panel's own address can only match an operator browsing that
		// address directly, which is same-origin anyway.
		expected := expectedHosts(c, opt)
		from := hostname(u.Host)
		for _, host := range expected {
			if strings.EqualFold(from, host) {
				c.Next()
				return
			}
		}

		// A mismatch only means something if the names it was compared against
		// mean something.
		if inconclusive(c, opt) {
			warnProxyHidesHost(c.Request.Host)
			c.Next()
			return
		}
		rejectCrossSite(c, from, expected)
	}
}

// expectedHosts returns every hostname a browser may legitimately have dialled
// to reach this panel.
//
// c.Request.Host is always one of them. Directly exposed it is the browser's own
// statement of where it went, and behind a proxy that forwards it, it still is;
// behind one that rewrote it, it names the panel's own socket, which only an
// operator browsing that address directly can match -- and that is same-origin.
//
// X-Forwarded-Host is read only with the reverse-proxy setting on, the same rule
// getRemoteIp applies to X-Forwarded-For. Which entry of a chain is taken does
// not matter, because it is not a CORS-safelisted header: a cross-site page
// cannot set one without a preflight this panel never answers, so a request
// carrying it did not come from the attack this guards against.
//
// The configured domain is accepted too -- though behind a proxy that rewrote
// Host it cannot be reached, because DomainValidator compares that same Host
// against the same domain and aborts first. What it does here is make the answer
// conclusive; see inconclusive.
func expectedHosts(c *gin.Context, opt Options) []string {
	hosts := []string{hostname(c.Request.Host)}
	if opt.PanelDomain != "" {
		hosts = append(hosts, hostname(opt.PanelDomain))
	}
	if opt.BehindProxy {
		for _, forwarded := range strings.Split(c.GetHeader("X-Forwarded-Host"), ",") {
			if forwarded = strings.TrimSpace(forwarded); forwarded != "" {
				hosts = append(hosts, hostname(forwarded))
			}
		}
	}
	return hosts
}

// inconclusive reports whether a mismatch proves nothing, because the only name
// the panel had to compare against was a Host header it can see was written by
// the proxy rather than by the browser.
//
// nginx's default is `proxy_set_header Host $proxy_host` -- the address the
// vhost dials the panel on -- so on a proxied panel that was not told to forward
// the real one, Host names nothing a browser ever typed and refusing on it
// refuses every write on an install that works, the login POST with them: there
// is then nowhere left to go and fix it from.
//
// Proven, not assumed. The panel knows the socket it bound, so a Host naming
// that socket is its own upstream address; any other Host came from somewhere
// that had to know the public name, and a mismatch against it is real. Without
// that half this stood down for every proxied panel with no configured domain,
// including the many whose vhost carries `proxy_set_header Host $host;` and
// were checked correctly before.
//
// A configured domain or a forwarded host settles the question by itself, so
// neither reaches this.
func inconclusive(c *gin.Context, opt Options) bool {
	if !opt.BehindProxy || opt.PanelDomain != "" {
		return false
	}
	if strings.TrimSpace(c.GetHeader("X-Forwarded-Host")) != "" {
		return false
	}
	return hostIsOwnSocket(c.Request.Host, opt.Listen, opt.Port)
}

// hostIsOwnSocket reports whether host names the address this panel is listening
// on, which is what $proxy_host expands to.
//
// The port has to match: a Host carrying some other port is not this socket. A
// Host with no port at all is the shape nginx's `Host $host` produces for a
// browser on 443, so it is not one of ours either.
func hostIsOwnSocket(host, listen string, port int) bool {
	h, p, err := net.SplitHostPort(host)
	if err != nil || p != strconv.Itoa(port) {
		return false
	}
	if listen != "" && listen != "0.0.0.0" && listen != "::" {
		// Bound to one address, so that address is the only one a proxy on
		// another host can dial -- and loopback below may not even be bound.
		return strings.EqualFold(h, listen)
	}
	if strings.EqualFold(h, "localhost") {
		return true
	}
	ip := net.ParseIP(h)
	return ip != nil && ip.IsLoopback()
}

// proxyHostWarned keeps the stand-down to one line in the log. It is reached on
// every write of every session on a panel whose proxy forwards no host, and the
// buffer GetLogs reads holds ten thousand lines in total.
var proxyHostWarned sync.Once

func warnProxyHidesHost(host string) {
	proxyHostWarned.Do(func() {
		// Only the two header fixes. Setting the panel domain looks like a
		// third way out and is the opposite: it mounts DomainValidator, which
		// compares this same rewritten Host and aborts *every* request, GETs
		// included, so the panel stops serving the login page at all.
		logger.Warning("the reverse proxy in front of this panel forwards neither",
			" the browser's Host nor X-Forwarded-Host (this request arrived as \"", host,
			"\"), so the same-origin check cannot run. Add",
			" `proxy_set_header Host $host;` or `proxy_set_header X-Forwarded-Host $host;`",
			" to the vhost.")
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
