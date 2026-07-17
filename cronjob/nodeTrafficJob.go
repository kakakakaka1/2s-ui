package cronjob

import (
	"sync"

	"github.com/shenaba/2s-ui/service"
)

type NodeTrafficJob struct {
	service.NodeSyncService
	running sync.Mutex
}

func NewNodeTrafficJob() *NodeTrafficJob {
	return &NodeTrafficJob{}
}

func (s *NodeTrafficJob) Run() {
	if !s.running.TryLock() {
		return
	}
	defer s.running.Unlock()
	s.NodeSyncService.CollectTraffic()
}
