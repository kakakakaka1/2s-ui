package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DomainValidator(domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		if colonIndex := strings.LastIndex(host, ":"); colonIndex != -1 {
			host, _, _ = net.SplitHostPort(c.Request.Host)
		}

		// 域名大小写不敏感(network/tls.go 的 SNI 校验也用 EqualFold)。精确比对会在
		// 设置里存了大小写混杂的域名时把整个面板 403 掉——TLS 握手反而成功,极难定位。
		if !strings.EqualFold(host, domain) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
