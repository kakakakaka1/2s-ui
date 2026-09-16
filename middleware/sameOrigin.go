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
		candidates := candidateHosts(c, opt)
		from := hostname(u.Host)
		for _, candidate := range candidates {
			if strings.EqualFold(from, hostname(candidate)) {
				c.Next()
				return
			}
		}

		// A mismatch only means something if the names it was compared against
		// mean something.
		if inconclusive(opt, candidates) {
			warnProxyHidesHost(c.Request.Host)
			c.Next()
			return
		}
		rejectCrossSite(c, from, candidates)
	}
}

// candidateHosts returns every address a browser may legitimately have dialled
// to reach this panel, each as the raw host[:port] it arrived as -- the port is
// what inconclusive needs to recognise the panel's own socket.
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
func candidateHosts(c *gin.Context, opt Options) []string {
	hosts := []string{c.Request.Host}
	if opt.PanelDomain != "" {
		hosts = append(hosts, opt.PanelDomain)
	}
	if opt.BehindProxy {
		for _, forwarded := range strings.Split(c.GetHeader("X-Forwarded-Host"), ",") {
			if forwarded = strings.TrimSpace(forwarded); forwarded != "" {
				hosts = append(hosts, forwarded)
			}
		}
	}
	return hosts
}

// inconclusive reports whether a mismatch proves nothing, because every name the
// panel had to compare against is one it can see is its own socket rather than
// anything a browser dialled.
//
// nginx's default is `proxy_set_header Host $proxy_host` -- the address the
// vhost dials the panel on -- so on a proxied panel that was not told to forward
// the real one, Host names nothing a browser ever typed and refusing on it
// refuses every write on an install that works, the login POST with them: there
// is then nowhere left to go and fix it from.
//
// Proven, not assumed, and asked of the whole set rather than of Host alone. One
// candidate that is not our own socket -- a forwarded host, the configured
// domain, a Host the vhost passed through -- had to come from something that
// knew the public name, so a mismatch against it is real and the check stands.
// Keying this on "is X-Forwarded-Host absent" instead missed the vhost that sets
// it from $proxy_host by mistake, which then refused every write; keying it on
// Host alone stood down for every proxied panel with no configured domain,
// including the many whose vhost carries `proxy_set_header Host $host;` and were
// checked correctly all along.
//
// candidateHosts always yields at least c.Request.Host, so this is never
// vacuously true.
func inconclusive(opt Options, candidates []string) bool {
	if !opt.BehindProxy {
		return false
	}
	for _, candidate := range candidates {
		if !hostIsOwnSocket(candidate, opt.Listen, opt.Port) {
			return false
		}
	}
	return true
}

// hostIsOwnSocket reports whether host names the address this panel is listening
// on, which is what $proxy_host expands to.
//
// The port has to match: a Host carrying some other port is not this socket. A
// Host with no port at all is the shape nginx's `Host $host` produces for a
// browser on 443, so it is not one of ours either.
//
// Every spelling of loopback counts, whatever webListen says, because the two
// are not interchangeable as strings and the cost of the two answers is not
// symmetric. A panel bound to 127.0.0.1 behind `proxy_pass http://localhost:2095`
// -- both sides correct -- sees "localhost:2095", and reading that as somebody
// else's name refuses every write including the login, with no way back in
// through the UI. Reading it as ours only gives up a check that was already
// unable to run in this configuration. Nothing is loosened by the extra
// spellings either: a browser can only send a Host it dialled, and an Origin
// naming the same loopback address matches outright one step earlier.
func hostIsOwnSocket(host, listen string, port int) bool {
	h, p, err := net.SplitHostPort(host)
	if err != nil || p != strconv.Itoa(port) {
		return false
	}
	if listen != "" && strings.EqualFold(h, listen) {
		return true
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
func rejectCrossSite(c *gin.Context, from string, candidates []string) {
	msg := "cross-site request refused"
	if from != "" {
		msg += ": this request says it came from \"" + from + "\""
		if len(candidates) > 0 {
			// The names as compared, so the operator sees the same two strings
			// the check did rather than a port that was never part of it.
			expected := make([]string, 0, len(candidates))
			for _, candidate := range candidates {
				expected = append(expected, hostname(candidate))
			}
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
