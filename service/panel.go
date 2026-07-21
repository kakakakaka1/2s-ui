package service

import (
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/shenaba/2s-ui/logger"
)

type PanelService struct {
}

// RestartPanel 在 delay 之后给自己发 SIGHUP(Windows 上 Kill)触发重启。
// delay 是 time.Duration,别传裸数字——那是纳秒,等于立即重启,调用方的 HTTP 响应
// 会来不及刷出就随 gin server 一起被拆掉。它存在的意义就是留出这段刷响应的时间。
func (s *PanelService) RestartPanel(delay time.Duration) error {
	p, err := os.FindProcess(syscall.Getpid())
	if err != nil {
		return err
	}
	go func() {
		time.Sleep(delay)
		if runtime.GOOS == "windows" {
			err = p.Kill()
		} else {
			err = p.Signal(syscall.SIGHUP)
		}
		if err != nil {
			logger.Error("send signal SIGHUP failed:", err)
		}
	}()
	return nil
}
