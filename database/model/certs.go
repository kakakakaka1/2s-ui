package model

// Cert 是一条手动登记的证书:用户自带的证书文件（Cloudflare 源证书、公司内部 CA
// 签发的，或升级前手填在设置里的那两个路径），面板只记下它在哪，续期归用户自己管。
//
// acme.sh 申请的证书【不】进这张表。那半边的事实来源是 acme.sh 自己的家目录，
// 面板扫盘就能读到全部信息（到期时间、下次续期、CA），再存一份副本只会两边不同步。
type Cert struct {
	Id uint `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	// 一个域名只会有一份在用的证书，唯一索引把「同一域名登记两次」挡在数据库层，
	// 省得上层每次写入都要先查一遍。
	Domain   string `json:"domain" form:"domain" gorm:"uniqueIndex"`
	CertFile string `json:"certFile" form:"certFile"`
	KeyFile  string `json:"keyFile" form:"keyFile"`
}
