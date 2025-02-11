package listeners

import (
	"context"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/common/pkg/events"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
	commonutils "github.com/uploadpilot/uploadpilot/common/pkg/utils"
	"github.com/uploadpilot/uploadpilot/uploader/internal/config"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/proc"
)

type WorkflowListener struct {
	logEb    *pubsub.EventBus[events.UploadLogEventMsg]
	uploadEb *pubsub.EventBus[events.UploadEventMsg]
	procSvc  *proc.ProcessorService
}

func NewWorkflowListener() *WorkflowListener {
	return &WorkflowListener{
		procSvc:  proc.NewProcessorService(),
		uploadEb: events.NewUploadStatusEvent(config.EventBusRedisConfig, uuid.New().String()),
		logEb:    events.NewUploadLogEventBus(config.EventBusRedisConfig, uuid.New().String()),
	}
}

func (l *WorkflowListener) Start() {
	defer commonutils.Recover()
	infra.Log.Info("starting upload Workflow listener...")
	group := "upload-workflow-listener"
	l.uploadEb.Subscribe(group, l.WorkflowTriggerHandler)
}

func (l *WorkflowListener) WorkflowTriggerHandler(msg *events.UploadEventMsg) error {
	ctx := context.Background()
	if msg.Status == string(models.UploadStatusComplete) {
		processors, err := l.procSvc.GetProcessors(ctx, msg.WorkspaceID)
		if err != nil {
			return err
		}
		for _, processor := range processors {
			l.procSvc.TriggerWorkflow(ctx, processor.Workflow)
			l.logEb.Publish(events.NewUploadLogEventMessage(msg.WorkspaceID, msg.UploadID, &processor.ID, nil, "workflow triggered", models.UploadLogLevelInfo))
		}
	}

	return nil
}
