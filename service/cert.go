package service

import (
	"crypto/tls"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shenaba/2s-ui/database"
	"github.com/shenaba/2s-ui/database/model"
	"github.com/shenaba/2s-ui/logger"
	"github.com/shenaba/2s-ui/util/common"

	"gorm.io/gorm/clause"
)

// CertService 管手动登记的证书，并把它们与 acme.sh 托管的那些合并成「域名与证书」
// 页面上的那一份清单。
//
// 分工：acme.sh 那半边的事实来源是它自己的家目录（AcmeService 扫盘直接读），只有
// 手动登记的这半边才落库。同一个域名两边都有时以 acme.sh 为准——它在自动续期，
// 手动记录多半是升级时归档进来的旧路径。
//
// 域名一律以小写存储和比对：DNS 不区分大小写，而 uniqueIndex、合并去重、前端匹配
// 全是精确比对，混着存会让同一个域名分裂成两条互相矛盾的记录。
//
// 字段用具名的 settings 而不是嵌入 SettingService：ApiService 已经嵌入了后者，
// 这里再嵌入会让 CertService 的方法提升到 ApiService 上产生二义性（AcmeService
// 那个注释说的也是这件事）。
type CertService struct {
	acme     AcmeService
	settings SettingService
}

// List 返回合并后的证书清单，按域名排序。
func (c *CertService) List() ([]CertInfo, error) {
	merged, err := c.acme.ListCerts()
	if err != nil {
		return nil, err
	}
	// ListCerts 已统一小写并按域名去重
	seen := make(map[string]bool, len(merged))
	for _, m := range merged {
		seen[m.Domain] = true
	}

	manual, err := c.listManual()
	if err != nil {
		return nil, err
	}
	for _, m := range manual {
		domain := strings.ToLower(m.Domain) // 存量库里可能还有大小写混杂的旧记录
		if seen[domain] {
			// acme.sh 已经在管这个域名了。手动记录留在库里不动（用户可能只是重复
			// 登记了一次），但清单上只显示会自动续期的那份，免得两行同名互相矛盾。
			continue
		}
		info := CertInfo{
			Domain:   domain,
			CertFile: m.CertFile,
			KeyFile:  m.KeyFile,
			Managed:  false,
		}
		// 叶子证书只解析一次，CA 和到期时间都从这一份上取
		if leaf := readLeafCert(m.CertFile); leaf != nil {
			info.CA = certIssuer(leaf)
			info.NotAfter = leaf.NotAfter.Unix()
		}
		merged = append(merged, info)
	}

	sort.Slice(merged, func(i, j int) bool { return merged[i].Domain < merged[j].Domain })
	return merged, nil
}

// FindByDomain 查某个域名当前可用的证书。
//
// 错误必须上抛而不是吞成「没找到」：调用方拿 found/Managed 决定走哪条分支，把一次
// 数据库故障当成「不存在」会让删除走错路（acme.sh 托管的被当成手动记录，DeleteManual
// 删了个寂寞还报成功）。
func (c *CertService) FindByDomain(domain string) (CertInfo, bool, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	if domain == "" {
		return CertInfo{}, false, nil
	}
	certs, err := c.List()
	if err != nil {
		return CertInfo{}, false, err
	}
	for _, ci := range certs {
		if ci.Domain == domain {
			return ci, true, nil
		}
	}
	return CertInfo{}, false, nil
}

func (c *CertService) listManual() ([]model.Cert, error) {
	db := database.GetDB()
	certs := []model.Cert{}
	if err := db.Model(model.Cert{}).Find(&certs).Error; err != nil {
		return nil, err
	}
	return certs, nil
}

// SaveManual 登记（或改写）一份自带的证书。
func (c *CertService) SaveManual(domain, certFile, keyFile string) error {
	return c.saveManual(domain, certFile, keyFile, true)
}

// saveManual 是登记的实体。verifyHost 控制「证书必须覆盖该域名」这一条校验：
// API 入口开着；启动归档关着——那是用户今天就在跑的组合，归档只求如实记录。
func (c *CertService) saveManual(domain, certFile, keyFile string, verifyHost bool) error {
	domain = strings.ToLower(strings.TrimSpace(domain))
	certFile = strings.TrimSpace(certFile)
	keyFile = strings.TrimSpace(keyFile)

	// 通配符放行：Cloudflare 源证书之类的自带证书常常就是 *.example.com
	if !validCertDomain(domain) {
		return common.NewErrorf("域名无效: %q", domain)
	}
	if certFile == "" || keyFile == "" {
		return common.NewError("证书路径和私钥路径都要填")
	}
	// 相对路径相对的是面板进程此刻的工作目录。与其拒绝（Windows 上 filepath.IsAbs
	// 连 /etc/ssl/cert.pem 都判为相对），不如在登记这一刻把它定格成绝对路径存下来：
	// 此后无论面板从哪里被拉起，读的都是登记时验证过的那个文件。
	if p, err := filepath.Abs(certFile); err == nil {
		certFile = p
	}
	if p, err := filepath.Abs(keyFile); err == nil {
		keyFile = p
	}
	if err := readableFile(certFile); err != nil {
		return common.NewErrorf("证书文件读不到: %v", err)
	}
	if err := readableFile(keyFile); err != nil {
		return common.NewErrorf("私钥文件读不到: %v", err)
	}
	// 证书与私钥必须配对。不配对的组合一旦被选中，面板下次重启在 LoadX509KeyPair
	// 上直接起不来——而 SIGHUP 重启路径会丢弃错误，表现为进程活着、服务全哑。
	if _, err := tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return common.NewErrorf("证书与私钥不匹配或无法解析: %v", err)
	}
	leaf := readLeafCert(certFile)
	if leaf == nil {
		return common.NewErrorf("%s 不是有效的 PEM 证书", certFile)
	}
	// 证书还得覆盖登记的这个域名。名字对不上的证书面板照常启动、握手也成功，但所有
	// 浏览器报名字不匹配，而面板侧一行错误日志都没有——必须在登记这一刻拦下。
	if verifyHost {
		if strings.HasPrefix(domain, "*.") {
			// VerifyHostname 不接受通配符当查询串，改查 SAN 里有没有这一条
			if !containsFold(leaf.DNSNames, domain) {
				return common.NewErrorf("证书不包含 %s（SAN: %s）", domain, strings.Join(leaf.DNSNames, ", "))
			}
		} else if err := leaf.VerifyHostname(domain); err != nil {
			return common.NewErrorf("证书不覆盖域名 %s: %v", domain, err)
		}
	}

	// 唯一索引已经在 domain 上，用 upsert 一次写入：手写 First→Create/Save 在两个
	// 请求同时登记同一域名时会撞索引报一条用户看不懂的错。
	db := database.GetDB()
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "domain"}},
		DoUpdates: clause.AssignmentColumns([]string{"cert_file", "key_file"}),
	}).Create(&model.Cert{Domain: domain, CertFile: certFile, KeyFile: keyFile}).Error
}

// DeleteManual 删掉一条登记记录。只忘掉路径，不碰证书文件本身——那是用户自己的
// 文件，多半还被别的服务用着。对不存在的域名是 no-op。
func (c *CertService) DeleteManual(domain string) error {
	domain = strings.ToLower(strings.TrimSpace(domain))
	if domain == "" {
		return common.NewError("域名不能为空")
	}
	return database.GetDB().Where("domain = ?", domain).Delete(model.Cert{}).Error
}

// HasManual 判断某个域名有没有手动登记记录。唯一索引上的一次 COUNT，零磁盘 I/O，
// 给 ArchiveLegacy 当每次启动的快速出口。
func (c *CertService) HasManual(domain string) bool {
	domain = strings.ToLower(strings.TrimSpace(domain))
	var count int64
	err := database.GetDB().Model(model.Cert{}).Where("domain = ?", domain).Count(&count).Error
	return err == nil && count > 0
}

// ArchiveLegacy 把升级前手填在设置里的证书路径归档成手动登记记录。
//
// 面板启动时调用，幂等：只在「域名非空 + 两个路径都非空且可读 + 该域名还没有任何
// 证书」时补一条。设置里的 webCertFile/webKeyFile 不动，面板照旧读它启动——这一步
// 只是让那份证书在新页面上看得见、可管理，不改变任何现有行为。
//
// 不放在 cmd/migration 里：那套迁移只在 `sui migrate` 和导入备份时跑，正常启动
// 根本不经过，挂在那儿等于对绝大多数升级用户没生效。
//
// 成本控制：正常启动在 HasManual（索引 COUNT）就短路了；acme.sh 家目录的全量扫描
// 只在两侧至少有一侧真要归档时才做，而且只做一次。
func (c *CertService) ArchiveLegacy() {
	all, err := c.settings.GetAllSetting()
	if err != nil {
		logger.Warning("归档存量证书路径失败，读设置出错: ", err)
		return
	}
	s := *all

	var managed map[string]bool
	managedHas := func(domain string) bool {
		if managed == nil {
			managed = map[string]bool{}
			if list, err := c.acme.ListCerts(); err == nil {
				for _, ci := range list {
					managed[ci.Domain] = true
				}
			}
		}
		return managed[domain]
	}

	for _, side := range []struct{ domain, cert, key string }{
		{s["webDomain"], s["webCertFile"], s["webKeyFile"]},
		{s["subDomain"], s["subCertFile"], s["subKeyFile"]},
	} {
		domain := strings.ToLower(strings.TrimSpace(side.domain))
		cert := strings.TrimSpace(side.cert)
		key := strings.TrimSpace(side.key)
		if domain == "" || cert == "" || key == "" || !validDomain(domain) {
			continue
		}
		if c.HasManual(domain) {
			continue
		}
		if readableFile(cert) != nil || readableFile(key) != nil {
			// 路径早就失效了（换过机器、证书搬过家）。归档一条读不到的记录只会在
			// 页面上多一行红字，不如放着——用户真要用会重新申请或重新登记。
			continue
		}
		if managedHas(domain) {
			continue
		}
		if err := c.saveManual(domain, cert, key, false); err != nil {
			logger.Warning("归档 ", domain, " 的存量证书路径失败: ", err)
			continue
		}
		logger.Info("已把 ", domain, " 手填的证书路径归档为手动登记的证书")
	}
}

// containsFold 大小写不敏感地找一个串（证书 SAN 里的通配符条目理论上可能带大写）。
func containsFold(list []string, want string) bool {
	for _, s := range list {
		if strings.EqualFold(s, want) {
			return true
		}
	}
	return false
}

// readableFile 要求路径存在、是普通文件、且当前进程读得开。只 os.Stat 不够：
// 目录、以及 root 之外的进程读不了的私钥，都能通过 Stat。
func readableFile(path string) error {
	st, err := os.Stat(path)
	if err != nil {
		return err
	}
	if st.IsDir() {
		return common.NewErrorf("%s 是目录", path)
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return f.Close()
}
