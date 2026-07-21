package service

import (
	"context"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shenaba/2s-ui/logger"
	"github.com/shenaba/2s-ui/util/common"
)

// 证书申请统一走系统 acme.sh,证书装到与 s-ui.sh 脚本一致的 /root/cert/{域名}/,
// 便于脚本 / nginx / sing-box 按固定文件名(fullchain.pem / privkey.pem)复用。
const (
	certBaseDir   = "/root/cert"        // 与脚本完全一致
	nginxConfDir  = "/etc/nginx/conf.d" // 自动生成的 ACME 验证 server 块所在目录
	acmeIssueTO   = 180 * time.Second   // 申请/安装超时
	acmeInstallTO = 120 * time.Second   // acme.sh / socat 安装超时
	cmdDetectTO   = 5 * time.Second     // 检测类命令超时
	// systemd 服务的 PATH 可能不全(且不继承登录 shell),补一个兜底,
	// 确保 exec 调用能定位 nginx/socat/systemctl/apt 等。
	fallbackPath = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
)

// 验证方式:standalone 独占 80 端口;nginx 借用运行中的 nginx,不中断 80 端口服务。
const (
	methodStandalone = "standalone"
	methodNginx      = "nginx"
)

// acmeIssuing 保证同一时刻只有一个证书申请在跑,避免前端重复点击导致并发申请、
// 撞 Let's Encrypt 限速。
var acmeIssuing sync.Mutex

// AcmeService 是无状态工具,不嵌入 SettingService(避免与 ApiService 已嵌入的
// SettingService 产生方法集二义性)。所有入参由调用方传入,不直接读写数据库。
type AcmeService struct{}

type NginxStatus struct {
	Installed  bool `json:"installed"`
	Active     bool `json:"active"`
	Port80Busy bool `json:"port80Busy"`
}

type IssueResult struct {
	CertFile  string `json:"certFile"`
	KeyFile   string `json:"keyFile"`
	Method    string `json:"method"`    // 实际使用的验证方式:standalone / nginx
	ReloadCmd string `json:"reloadCmd"` // 续期后的重载命令,空表示没配(前端据此提示自配钩子)
}

// withHome 返回把 HOME 固定为指定值、并补全 PATH 兜底的环境变量(去重)。
// systemd 服务即使以 root 运行也往往不设 HOME、PATH 也可能不全。
func withHome(home string) []string {
	env := os.Environ()
	out := make([]string, 0, len(env)+2)
	hasPath := false
	for _, e := range env {
		switch {
		case strings.HasPrefix(e, "HOME="):
			continue
		case strings.HasPrefix(e, "PATH="):
			hasPath = true
			e = e + ":" + fallbackPath // 合并兜底路径,重复项无害
		}
		out = append(out, e)
	}
	out = append(out, "HOME="+home)
	if !hasPath {
		out = append(out, "PATH="+fallbackPath)
	}
	return out
}

// resolveBin 定位外部命令:先按进程自身 PATH 找(exec 定位二进制用的是进程 PATH,
// withHome 注入的兜底 PATH 只对子进程内部生效,如 acme.sh 再调 socat),找不到再扫
// fallbackPath,保证极简 PATH 的服务环境下也能找到 nginx/systemctl/apt 等。
func resolveBin(name string) (string, bool) {
	if p, err := exec.LookPath(name); err == nil {
		return p, true
	}
	for _, dir := range strings.Split(fallbackPath, ":") {
		p := filepath.Join(dir, name)
		if st, err := os.Stat(p); err == nil && !st.IsDir() && st.Mode()&0111 != 0 {
			return p, true
		}
	}
	return name, false // 原样返回,让 exec 报标准的 not found 错误
}

// runCmd 执行外部命令(HOME 固定为 home),合并 stdout/stderr,超时或非零退出码
// 都包成 error,并把输出原文附在错误里回传前端,便于排查(80 端口被占、域名未解析等)。
func runCmd(timeout time.Duration, home, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	bin, _ := resolveBin(name)
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = withHome(home)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))

	if ctx.Err() == context.DeadlineExceeded {
		return output, common.NewErrorf("命令超时(%v): %s %s\n%s", timeout, name, strings.Join(args, " "), output)
	}
	if err != nil {
		return output, common.NewErrorf("命令执行失败: %s %s\n%s\n%v", name, strings.Join(args, " "), output, err)
	}
	return output, nil
}

// resolveAcmeSh 在常见位置查找已安装的 acme.sh,返回可执行路径及其对应的 HOME。
// systemd 服务下 HOME 常为空,acme.sh 会装到 /.acme.sh;s-ui.sh 脚本则装到
// /root/.acme.sh。都纳入探测,避免硬编码单一路径导致"找不到"。
func resolveAcmeSh() (bin, home string) {
	candidates := make([]string, 0, 3)
	if h := os.Getenv("HOME"); h != "" {
		candidates = append(candidates, filepath.Join(h, ".acme.sh", "acme.sh"))
	}
	candidates = append(candidates, "/root/.acme.sh/acme.sh", "/.acme.sh/acme.sh")
	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			// home = 可执行文件上两级目录:/x/.acme.sh/acme.sh -> /x
			return p, filepath.Dir(filepath.Dir(p))
		}
	}
	return "", ""
}

// port80Free 探测本机 80 端口是否空闲(standalone 验证需要独占 80)。
func port80Free() bool {
	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// port80Hint 在面板非 root 运行时补一句说明:非特权进程绑 :80 拿到的是 EACCES,
// 在 port80Free 里与「端口被占用」同样是 false,错误文案不区分会把排查方向带偏。
func port80Hint() string {
	if os.Geteuid() != 0 {
		return "(注意:面板当前不是以 root 运行,绑定 80 端口会被内核直接拒绝,这也可能是本次判定不可用的真正原因)"
	}
	return ""
}

// ipv6Available 探测内核能否创建 IPv6 监听(ipv6.disable=1 的主机不行):
// 生成 nginx 验证块时据此决定是否加 listen [::]:80,避免 reload 因开不了 v6 socket 失败。
func ipv6Available() bool {
	ln, err := net.Listen("tcp6", "[::1]:0")
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// hasGlobalIPv4 判断主机是否拥有全局 IPv4 地址(含 NAT 后的私网地址,与 ip -4 addr
// show scope global 语义一致,只排除回环与链路本地)。
// 用途:acme.sh 的 standalone 服务器默认只绑 IPv4,而 --listen-v6 是【排他】的
// (v6-only),双栈主机加上它反而会让 A 记录指向的 v4 验证失败——故仅在纯 IPv6 主机
// 上才加该标志。探测失败按「有 v4」处理:宁可维持默认,也不要把好用的双栈主机弄挂。
func hasGlobalIPv4() bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return true
	}
	for _, a := range addrs {
		ipn, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		if ipn.IP.To4() != nil && !ipn.IP.IsLoopback() && !ipn.IP.IsLinkLocalUnicast() {
			return true
		}
	}
	return false
}

// validDomain 只放行主机名合法字符:域名会拼进文件路径与外部命令参数,严格校验兜底。
func validDomain(d string) bool {
	if len(d) == 0 || len(d) > 253 || d[0] == '-' || d[0] == '.' {
		return false
	}
	for _, r := range d {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '.', r == '-':
		default:
			return false
		}
	}
	return true
}

// DetectNginx 检测 nginx 是否安装并运行,以及 80 端口是否被占用。Windows 直接返回零值。
func (a *AcmeService) DetectNginx() NginxStatus {
	status := NginxStatus{}
	if runtime.GOOS == "windows" {
		return status
	}
	status.Port80Busy = !port80Free()
	if _, ok := resolveBin("nginx"); ok {
		status.Installed = true
	}
	// systemctl is-active 在运行时退出码为 0、输出 "active"。
	// 局限:无 systemd 的环境(如容器里直跑 nginx)检不出 active,auto 只会走 standalone/报错。
	if out, err := runCmd(cmdDetectTO, "/root", "systemctl", "is-active", "nginx"); err == nil && strings.TrimSpace(out) == "active" {
		status.Active = true
		status.Installed = true
	}
	return status
}

// resolveMethod 把前端传入的验证方式解析为实际可执行的方式并校验可行性。
//   - standalone:需独占 80 端口。
//   - nginx     :需 nginx 正在运行(不中断 80 端口服务)。
//   - auto / 空 :80 空闲优先 standalone;80 被占且 nginx 在跑则借用 nginx;否则报错。
func (a *AcmeService) resolveMethod(method string) (string, error) {
	switch method {
	case methodStandalone:
		if !port80Free() {
			return "", common.NewErrorf("80 端口被占用,无法用 standalone 申请;若 nginx 正在运行请改选 nginx 验证或「自动」,否则请先停止占用 80 端口的服务%s", port80Hint())
		}
		return methodStandalone, nil
	case methodNginx:
		if !a.DetectNginx().Active {
			return "", common.NewError("未检测到正在运行的 nginx,无法用 nginx 验证;请先启动 nginx,或改选「自动」")
		}
		return methodNginx, nil
	case "", "auto":
		if port80Free() {
			return methodStandalone, nil
		}
		if a.DetectNginx().Active {
			return methodNginx, nil
		}
		return "", common.NewErrorf("80 端口不可用且未检测到运行中的 nginx:请停止占用 80 端口的程序,或启动 nginx 后重试%s", port80Hint())
	default:
		return "", common.NewErrorf("未知的验证方式: %q", method)
	}
}

// nginxHasServerName 在 nginx -T 的完整生效配置里查找是否已有 server_name 含该域名。
// 按行解析,值换行的合法写法(server_name 与域名不同行)检不出——后果只是多生成一个
// 冗余块:nginx 对重名仅告警 conflicting server name,-t 仍通过,无害。
func nginxHasServerName(conf, domain string) bool {
	const directive = "server_name"
	for _, line := range strings.Split(conf, "\n") {
		line = strings.TrimSpace(line)
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = line[:i]
		}
		if !strings.HasPrefix(line, directive) {
			continue
		}
		rest := line[len(directive):]
		if rest == "" || (rest[0] != ' ' && rest[0] != '\t') {
			continue // 跳过 server_names_hash_bucket_size 等同前缀指令
		}
		for _, name := range strings.Fields(strings.TrimSuffix(strings.TrimSpace(rest), ";")) {
			if strings.EqualFold(name, domain) {
				return true
			}
		}
	}
	return false
}

// ensureNginxServerBlock 确保 nginx 配置里有 server_name 匹配该域名的 server 块——
// acme.sh --nginx 靠它定位要临时改写的配置,缺失会报 "Can not find conf file"。
// 缺失时生成自包含的最小验证块(只加自己的文件,不碰用户已有配置):nginx -t 通过才
// reload,失败立即回滚删除;文件常驻,以后自动续期仍靠它验证。
func (a *AcmeService) ensureNginxServerBlock(domain string) error {
	if out, err := runCmd(cmdDetectTO, "/root", "nginx", "-T"); err == nil && nginxHasServerName(out, domain) {
		return nil
	}
	confPath := filepath.Join(nginxConfDir, "s-ui-acme-"+domain+".conf")
	// 域名有 AAAA 记录时 Let's Encrypt 优先走 IPv6:若本块只听 v4,校验请求会被别的
	// [::]:80 块(如发行版默认站点的 default_server)接走导致 404,故 v6 可用时一并监听。
	listen := "    listen 80;\n"
	if ipv6Available() {
		listen += "    listen [::]:80;\n"
	}
	content := "# Generated by s-ui: ACME HTTP-01 validation block for " + domain + ".\n" +
		"# Kept in place so automatic certificate renewal keeps working.\n" +
		"server {\n" +
		listen +
		"    server_name " + domain + ";\n" +
		"    root /var/www/html;\n" +
		"}\n"
	// 先确认目录本身在:Alpine 用的是 /etc/nginx/http.d,源码编译的 nginx 可能压根没有
	// conf.d。不预判就会撞上一条没有指向性的 ENOENT,与下面「文件写了但没被 include」
	// 是同一类问题(本机 nginx 布局不符合假设),给同样指向根因的错误。
	if st, err := os.Stat(nginxConfDir); err != nil || !st.IsDir() {
		return common.NewErrorf("nginx 配置目录 %s 不存在,无法自动生成验证配置"+
			"(Alpine 通常是 /etc/nginx/http.d,源码编译的 nginx 可能没有此目录);"+
			"请改用 standalone 验证,或手动为 %s 添加一个 server_name 匹配的 server 块",
			nginxConfDir, domain)
	}
	if err := os.WriteFile(confPath, []byte(content), 0644); err != nil {
		return common.NewErrorf("写入 nginx 验证配置失败 %s: %v", confPath, err)
	}
	if out, err := runCmd(cmdDetectTO, "/root", "nginx", "-t"); err != nil {
		_ = os.Remove(confPath) // 绝不留下会卡住 reload 的配置
		return common.NewErrorf("自动生成的 nginx 验证配置未通过 nginx -t,已回滚删除 %s:\n%s", confPath, out)
	}
	if out, err := runCmd(cmdDetectTO, "/root", "systemctl", "reload", "nginx"); err != nil {
		_ = os.Remove(confPath)
		return common.NewErrorf("reload nginx 失败,已回滚删除 %s:\n%s", confPath, out)
	}
	// 复验块是否真的生效:若本机 nginx.conf 没有 include conf.d/*.conf(源码编译、
	// openresty、手写配置都常见),上面三步会全部「成功」——文件根本没被解析,nginx -t
	// 自然通过——随后 acme.sh 才报它自己的 "Can not find conf file",离真因十万八千里。
	// 复验把它变成一条指向根因的错误。
	if out, err := runCmd(cmdDetectTO, "/root", "nginx", "-T"); err != nil || !nginxHasServerName(out, domain) {
		_ = os.Remove(confPath)
		// 若文件其实已被 include(nginx -T 因别的原因失败),上面那次 reload 已把块读进
		// 内存,光删文件不足以复原;再 reload 一次抹掉。没被 include 时这次是空操作。
		_, _ = runCmd(cmdDetectTO, "/root", "systemctl", "reload", "nginx")
		return common.NewErrorf("已写入 %s 但它未出现在 nginx 生效配置中,"+
			"本机 nginx.conf 可能没有 include %s/*.conf;已回滚删除,请改用 standalone 验证",
			confPath, nginxConfDir)
	}
	logger.Info("已生成 nginx 验证配置:", confPath)
	return nil
}

// ensureAcmeSh 确保 acme.sh 可用,返回其可执行路径与对应 HOME;缺失时自动安装。
// 安装时在 shell 内显式 export HOME=/root,确保装到 /root/.acme.sh(systemd 下
// HOME 常为空,否则会装到 /.acme.sh);并 curl/wget 自适应(最小化系统可能只有其一)。
func ensureAcmeSh() (bin, home string, err error) {
	if bin, home = resolveAcmeSh(); bin != "" {
		return bin, home, nil
	}
	logger.Info("acme.sh 未安装,开始自动安装...")
	installScript := "export HOME=/root; " +
		"if command -v curl >/dev/null 2>&1; then curl https://get.acme.sh | sh; " +
		"else wget -O - https://get.acme.sh | sh; fi"
	out, e := runCmd(acmeInstallTO, "/root", "sh", "-c", installScript)
	if e != nil {
		return "", "", common.NewErrorf("自动安装 acme.sh 失败,请在服务器手动安装(s-ui 脚本 SSL 菜单)。详情:\n%s", out)
	}
	if bin, home = resolveAcmeSh(); bin != "" {
		logger.Info("acme.sh 安装成功:", bin)
		return bin, home, nil
	}
	return "", "", common.NewErrorf("acme.sh 安装后仍未找到(已查 $HOME/.acme.sh、/root/.acme.sh、/.acme.sh)。安装输出:\n%s", out)
}

// ensureSocat 仅 standalone 申请需要,best-effort 安装,失败不致命(由后续申请报真实错)。
func ensureSocat() {
	if _, ok := resolveBin("socat"); ok {
		return
	}
	logger.Info("socat 未安装,尝试自动安装(standalone 申请需要)...")
	managers := [][]string{
		{"apt", "-y", "install", "socat"},
		{"yum", "-y", "install", "socat"},
		{"dnf", "-y", "install", "socat"},
		{"pacman", "-Sy", "--noconfirm", "socat"},
	}
	for _, m := range managers {
		if _, ok := resolveBin(m[0]); ok {
			if _, err := runCmd(acmeInstallTO, "/root", m[0], m[1:]...); err == nil {
				return
			}
		}
	}
	logger.Warning("socat 自动安装失败,若 standalone 申请报错请手动安装 socat")
}

// IssueWeb 为面板/订阅申请证书并安装到 /root/cert/{域名}/。
//   - method     :standalone / nginx / auto(空视同 auto),解析与可行性校验见 resolveMethod。
//   - force      :域名已有未到期证书时 acme.sh 默认跳过签发,force 时加 --force 强制续期。
//   - behindProxy:webNginx=true,即反向代理终结 TLS、nginx 是证书消费方(决定 reloadcmd)。
func (a *AcmeService) IssueWeb(domain, email, method string, force, behindProxy bool) (*IssueResult, error) {
	if runtime.GOOS == "windows" {
		return nil, common.NewError("Windows 不支持 acme.sh 申请证书")
	}
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return nil, common.NewError("域名不能为空")
	}
	if !validDomain(domain) {
		return nil, common.NewErrorf("域名含非法字符: %q", domain)
	}
	// 同一时刻只允许一个申请,避免重复点击并发申请撞 Let's Encrypt 限速
	if !acmeIssuing.TryLock() {
		return nil, common.NewError("已有证书申请正在进行,请等当前申请完成后再试")
	}
	defer acmeIssuing.Unlock()

	resolved, err := a.resolveMethod(method)
	if err != nil {
		return nil, err
	}

	bin, home, err := ensureAcmeSh()
	if err != nil {
		return nil, err
	}

	// 默认 CA 设为 Let's Encrypt(与脚本一致)
	if out, err := runCmd(cmdDetectTO, home, bin, "--set-default-ca", "--server", "letsencrypt"); err != nil {
		return nil, common.NewErrorf("设置默认 CA 失败:\n%s", out)
	}

	// 申请证书
	issueArgs := []string{"--issue", "-d", domain}
	if email != "" {
		issueArgs = append(issueArgs, "--accountemail", email)
	}
	if resolved == methodNginx {
		if err := a.ensureNginxServerBlock(domain); err != nil {
			return nil, err
		}
		issueArgs = append(issueArgs, "--nginx")
	} else {
		ensureSocat()
		// resolveMethod 的 80 端口预检发生在 ensureAcmeSh / ensureSocat 之前,而这两步
		// 首次运行各可能耗掉 120s 装包——恰恰是包管理器可能拉起 :80 上某个服务的时刻。
		// 预检负责给出清晰的早期错误,这里紧贴使用点再确认一次,收窄那段窗口。
		if !port80Free() {
			return nil, common.NewError("80 端口在准备阶段被占用,无法用 standalone 申请;请停止占用 80 端口的程序后重试")
		}
		issueArgs = append(issueArgs, "--standalone", "--httpport", "80")
		// 纯 IPv6 主机必须显式切到 v6 监听,否则 acme.sh 只绑 v4、LE 的请求根本进不来。
		// 判据是「有没有全局 v4」而非「内核支不支持 v6」——该标志是排他的,见 hasGlobalIPv4。
		if !hasGlobalIPv4() {
			issueArgs = append(issueArgs, "--listen-v6")
		}
	}
	// 域名已有未到期证书时 acme.sh 会跳过("Skipping. Next renewal time is ...")，
	// --force 强制重新签发以续期。会消耗 Let's Encrypt 限速额度，故由前端「强制续期」显式触发。
	if force {
		issueArgs = append(issueArgs, "--force")
	}
	if out, err := runCmd(acmeIssueTO, home, bin, issueArgs...); err != nil {
		hint := "证书申请失败"
		if !force {
			hint = "证书申请失败(若域名已有未到期证书,请改用「强制续期」)"
		}
		return nil, common.NewErrorf("%s:\n%s", hint, out)
	}

	// 安装证书到 /root/cert/{域名}/
	certDir := filepath.Join(certBaseDir, domain)
	if err := os.MkdirAll(certDir, 0755); err != nil {
		return nil, common.NewErrorf("创建证书目录失败 %s: %v", certDir, err)
	}
	keyFile := filepath.Join(certDir, "privkey.pem")
	certFile := filepath.Join(certDir, "fullchain.pem")
	installArgs := []string{
		"--installcert", "-d", domain,
		"--key-file", keyFile,
		"--fullchain-file", certFile,
	}
	// reloadcmd 按「谁消费证书」决定(见 buildReloadCmd)。standalone 且面板自用时不配:
	//  - acme.sh 首次 --installcert 会内联执行一次 reloadcmd,若是重启面板会杀掉正在
	//    处理本次申请请求的进程,前端拿不到结果(即便证书已申请成功);
	//  - 面板/订阅侧也无需 reloadcmd:network/tls.go 的 certReloader 每次 TLS 握手按
	//    文件 mtime 热加载,续期覆盖 /root/cert 下的文件后自动生效。
	reloadCmd := a.buildReloadCmd(resolved, behindProxy)
	if reloadCmd != "" {
		installArgs = append(installArgs, "--reloadcmd", reloadCmd)
	}
	if out, err := runCmd(acmeIssueTO, home, bin, installArgs...); err != nil {
		return nil, common.NewErrorf("安装证书失败:\n%s", out)
	}

	// 启用 acme.sh 自带 cron 自动续期(失败不影响本次证书)
	if _, err := runCmd(acmeInstallTO, home, bin, "--upgrade", "--auto-upgrade"); err != nil {
		logger.Warning("启用 acme.sh 自动续期失败(不影响本次证书):", err)
	}

	return &IssueResult{CertFile: certFile, KeyFile: keyFile, Method: resolved, ReloadCmd: reloadCmd}, nil
}

// buildReloadCmd 决定续期成功后的重载命令,按「谁消费证书」而非验证方式:
//   - nginx 验证或反代终结 TLS(behindProxy):nginx 持有/使用证书,续期后必须让它
//     重读,否则续期只覆盖了磁盘文件,nginx 仍用内存里的旧证书,90 天后线上过期。
//     用 try-reload-or-restart:nginx 在跑→reload;没在跑→无事退 0,避免 reloadcmd
//     失败把 --installcert 整体判失败(彼时证书文件其实已安装成功)。
//   - 其余(standalone 且面板/订阅自用):返回空——证书由 certReloader 热加载,
//     无需外部命令(见 IssueWeb 注释)。
//
// behindProxy 只说明「TLS 由反向代理终结」,没说那代理是 nginx——Caddy / Traefik /
// HAProxy 的用户同样会开这个开关。此时不能硬发 nginx 命令:try-reload-or-restart
// 的「不在跑就空转退 0」只对【已知但未激活】的 unit 成立,unit 根本不存在时 systemd
// 以 5 退出,会让 --installcert 整体判失败(彼时证书其实已落盘),且这条注定失败的
// 命令还会被写进 acme.sh 的域名 conf,让此后每次续期都报错。systemctl 本身缺席
// (Alpine/OpenRC、Devuan、容器里直跑 nginx)以 127 退出,后果完全一样。
// 这条路径经 behindProxy 可达(代理只听 443 → 80 空闲 → standalone 验证 → 走到
// 这里),不是理论情况;确认不了就留空,由调用方提示用户自行配置重载钩子。
func (a *AcmeService) buildReloadCmd(method string, behindProxy bool) string {
	if method != methodNginx && !behindProxy {
		return ""
	}
	// 判据必须是「systemd 认识 nginx 这个 unit」,不能是「nginx 二进制在」:源码编译的
	// nginx 装在 /usr/local/sbin(正好在 fallbackPath 上)却往往没有 unit 文件,查
	// 二进制会放行,随后 try-reload-or-restart 退 5,正是上面要避免的那种失败。
	// systemctl cat 恰好是需要的语义:unit 已知→退 0(哪怕没在跑),不存在→非 0;
	// systemctl 自身缺席时 runCmd 直接报 not found,同样落进 err 分支。
	if _, err := runCmd(cmdDetectTO, "/root", "systemctl", "cat", "nginx"); err != nil {
		return ""
	}
	return "systemctl try-reload-or-restart nginx"
}

// ListCerts 返回 acme.sh 已管理的证书列表,供前端展示、避免重复申请触发 LE 限速。
func (a *AcmeService) ListCerts() (string, error) {
	if runtime.GOOS == "windows" {
		return "", common.NewError("Windows 不支持 acme.sh")
	}
	bin, home := resolveAcmeSh()
	if bin == "" {
		return "", nil
	}
	out, _ := runCmd(cmdDetectTO, home, bin, "--list")
	return out, nil
}
