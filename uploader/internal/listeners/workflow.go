package listeners

import (
	"context"

	"github.com/google/uuid"
	commonutils "github.com/uploadpilot/uploadpilot/go-core/common/utils"
	"github.com/uploadpilot/uploadpilot/go-core/db/pkg/models"
	"github.com/uploadpilot/uploadpilot/go-core/pubsub/pkg/events"
	"github.com/uploadpilot/uploadpilot/uploader/internal/infra"
	"github.com/uploadpilot/uploadpilot/uploader/internal/svc/processor"
)

type WorkflowListener struct {
	uploadEvent    *events.UploadStatusEvent
	uploadLogEvent *events.UploadLogEvent
	procSvc        *processor.Service
}

func NewWorkflowListener(procSvc *processor.Service) *WorkflowListener {
	return &WorkflowListener{
		procSvc:        procSvc,
		uploadEvent:    events.NewUploadStatusEvent(infra.RedisClient),
		uploadLogEvent: events.NewUploadLogEvent(infra.RedisClient),
	}
}

func (l *WorkflowListener) Start() {
	defer commonutils.Recover()
	infra.Log.Info("starting upload Workflow listener...")

	consumerGroup := "upload-workflow-listener"
	consumerKey := uuid.NewString()
	l.uploadEvent.Subscribe(consumerGroup, consumerKey, l.WorkflowTriggerHandler)
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
			l.uploadLogEvent.Publish(msg.WorkspaceID, msg.UploadID, &processor.ID, nil, "workflow triggered", string(models.UploadLogLevelInfo))
		}
	}

	return nil
}
