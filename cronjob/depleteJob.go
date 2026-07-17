package cronjob

import (
	"github.com/shenaba/2s-ui/database"
	"github.com/shenaba/2s-ui/logger"
	"github.com/shenaba/2s-ui/service"
)

type DepleteJob struct {
	service.ClientService
	service.InboundService
	service.NodeSyncService
}

func NewDepleteJob() *DepleteJob {
	return new(DepleteJob)
}

func (s *DepleteJob) Run() {
	inboundIds, disabled, err := s.ClientService.DepleteClients()
	if err != nil {
		logger.Warning("Disable depleted users failed: ", err)
		return
	}
	if len(inboundIds) > 0 {
		// RestartInbounds already filters to node_id IS NULL — only local
		// inbounds are hot-restarted here.
		err := s.InboundService.RestartInbounds(database.GetDB(), inboundIds)
		if err != nil {
			logger.Error("unable to restart inbounds: ", err)
		}
	}
	// Fan the disable out to nodes: reconcile sees enable=false in the expected
	// set and pushes an edit; the node then hot-restarts and drops connections.
	if len(disabled) > 0 {
		s.NodeSyncService.MarkAllDirty()
		go s.NodeSyncService.ReconcileDirtyOnline()
	}
}
