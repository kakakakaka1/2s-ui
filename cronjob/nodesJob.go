package cronjob

import (
	"sync"

	"github.com/shenaba/2s-ui/service"
)

type NodesJob struct {
	service.NodeService
	running sync.Mutex
}

func NewNodesJob() *NodesJob {
	return &NodesJob{}
}

func (s *NodesJob) Run() {
	// robfig/cron does not wait for the previous run; with a 5s cadence and a
	// 4s probe timeout a slow batch could overlap itself — skip instead.
	if !s.running.TryLock() {
		return
	}
	defer s.running.Unlock()
	s.NodeService.RefreshAll()
}
