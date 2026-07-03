package util

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"net"
	"os"
	"strings"
	"time"

	"github.com/shenaba/2s-ui/util/common"

	utls "github.com/refraction-networking/utls"
)

func CertPEMFromTLS(tlsConfig map[string]interface{}) string {
	if tlsConfig == nil {
		return ""
	}
	switch c := tlsConfig["certificate"].(type) {
	case string:
		if c != "" {
			return c
		}
	case []interface{}:
		lines := make([]string, 0, len(c))
		for _, l := range c {
			if s, ok := l.(string); ok {
				lines = append(lines, s)
			}
		}
		if len(lines) > 0 {
			return strings.Join(lines, "\n")
		}
	}
	if path, ok := tlsConfig["certificate_path"].(string); ok && path != "" {
		if data, err := os.ReadFile(path); err == nil {
			return string(data)
		}
	}
	return ""
}

func parseLeafCert(pemData string) *x509.Certificate {
	rest := []byte(pemData)
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			return nil
		}
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil
			}
			return cert
		}
	}
}

// CertIsSelfSigned reports whether the leaf certificate in pemData is
// self-signed, i.e. its signature verifies against its own public key. Only
// self-signed certificates should be pinned via certificate_public_key_sha256;
// CA-signed certificates are validated normally.
func CertIsSelfSigned(pemData string) bool {
	cert := parseLeafCert(pemData)
	if cert == nil {
		return false
	}
	return cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature) == nil
}

// CertPublicKeySha256 returns the base64-encoded SHA256 of the certificate's
// SubjectPublicKeyInfo (sing-box `certificate_public_key_sha256` / link pinSHA256).
func CertPublicKeySha256(pemData string) string {
	cert := parseLeafCert(pemData)
	if cert == nil {
		return ""
	}
	sum := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
	return base64.StdEncoding.EncodeToString(sum[:])
}

// CertSha256Hex returns the lowercase hex SHA256 of the whole certificate (DER),
// matching `openssl x509 -fingerprint -sha256` and Clash/mihomo's `fingerprint`.
func CertSha256Hex(pemData string) string {
	cert := parseLeafCert(pemData)
	if cert == nil {
		return ""
	}
	sum := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(sum[:])
}

// GetTlsPing performs a TLS handshake against domain:port (uTLS Chrome hello,
// no verification -- the point is to fetch whatever certificate is served) and
// returns the leaf certificate's SPKI SHA256 so the UI can show/pin it.
func GetTlsPing(domain string, port string) (any, error) {
	if domain == "" {
		return "", common.NewError("domain is empty")
	}
	if port == "" {
		port = "443"
	}

	d := net.Dialer{Timeout: 10 * time.Second}
	tcpConn, err := d.Dial("tcp", net.JoinHostPort(domain, port))
	if err != nil {
		return "", common.NewErrorf("Failed to dial tcp: %s", err)
	}
	defer tcpConn.Close()
	tlsConn := utls.UClient(tcpConn, &utls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2", "http/1.1"},
	}, utls.HelloChrome_Auto)
	if err = tlsConn.Handshake(); err != nil {
		return "", common.NewErrorf("Failed to handshake: %s", err)
	}
	// Prefer the certificate that carries SANs; fall back to the first one so
	// CN-only (typically self-signed) certificates don't crash the probe.
	var leaf *x509.Certificate
	certs := tlsConn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		if len(cert.DNSNames) != 0 {
			leaf = cert
			break
		}
	}
	if leaf == nil {
		if len(certs) == 0 {
			return "", common.NewError("no peer certificate received")
		}
		leaf = certs[0]
	}
	sum := sha256.Sum256(leaf.RawSubjectPublicKeyInfo)
	return map[string]string{
		"leafHash": base64.StdEncoding.EncodeToString(sum[:]),
	}, nil
}
