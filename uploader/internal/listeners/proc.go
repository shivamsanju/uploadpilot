package listeners

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/uploadpilot/uploadpilot/common/pkg/db"
	"github.com/uploadpilot/uploadpilot/common/pkg/events"
	"github.com/uploadpilot/uploadpilot/common/pkg/infra"
	"github.com/uploadpilot/uploadpilot/common/pkg/models"
	"github.com/uploadpilot/uploadpilot/common/pkg/pubsub"
	commonutils "github.com/uploadpilot/uploadpilot/common/pkg/utils"
	"github.com/uploadpilot/uploadpilot/uploader/internal/config"
)

type ProcListener struct {
	uploadEb *pubsub.EventBus[events.UploadEventMsg]
	logEb    *pubsub.EventBus[events.UploadLogEventMsg]

	pRepo  *db.ProcessorRepo
	upRepo *db.UploadRepo
}

func NewProcListener() *ProcListener {
	return &ProcListener{
		pRepo:    db.NewProcessorRepo(),
		upRepo:   db.NewUploadRepo(),
		uploadEb: events.NewUploadStatusEvent(config.EventBusRedisConfig, uuid.New().String()),
		logEb:    events.NewUploadLogEventBus(config.EventBusRedisConfig, uuid.New().String()),
	}
}

func (l *ProcListener) Start() {
	infra.Log.Info("starting listening for upload complete events...")
	group := "proc-listener"
	l.uploadEb.Subscribe(group, l.procHandler)
}

func (l *ProcListener) Stop() {
	l.uploadEb.Unsubscribe()
}

func (l *ProcListener) procHandler(msg *events.UploadEventMsg) error {
	defer commonutils.Recover()
	ctx := context.Background()
	if msg.Status != string(models.UploadStatusComplete) {
		return nil
	}

	processors, err := l.pRepo.GetAll(ctx, msg.WorkspaceID)
	if err != nil {
		infra.Log.Errorf("failed to get processors: %s", err.Error())
		return err
	}

	if len(processors) == 0 {
		infra.Log.Info("no processors found, skipping processing...")
		msg.Status = string(models.UploadStatusSkipped)
		l.uploadEb.Publish(msg)
		return err
	}

	wg := &sync.WaitGroup{}
	for _, processor := range processors {
		wg.Add(1)
		go l.startTaskProcessor(wg, msg, &processor)
	}
	wg.Wait()

	return nil
}

func (l *ProcListener) startTaskProcessor(wg *sync.WaitGroup, msg *events.UploadEventMsg, processor *models.Processor) {
	defer commonutils.Recover()
	defer wg.Done()

	wID := msg.WorkspaceID
	uID := msg.UploadID
	pID := processor.ID

	// context with timeout of 15 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	upload, err := l.upRepo.Get(ctx, uID)
	if err != nil {
		l.logEb.Publish(events.NewUploadLogEventMessage(wID, uID, &pID, nil, err.Error(), models.UploadLogLevelError))
		return
	}

	p, err := l.pRepo.Get(ctx, pID)
	if err != nil {
		l.logEb.Publish(events.NewUploadLogEventMessage(wID, uID, &pID, nil, err.Error(), models.UploadLogLevelError))
		return
	}

	fileType := upload.Metadata["filetype"].(string)
	if len(p.Triggers) > 0 && !slices.Contains(p.Triggers, fileType) {
		return
	}

	msg.Status = string(models.UploadStatusProcessing)
	l.uploadEb.Publish(msg)

	// r := proc.NewProcWorkflowRunner()
	// if err := r.Build(ctx, wID, pID, uID); err != nil {
	// 	m := fmt.Sprintf("failed to build workflow for processor %s and upload %s: %s", processor.Name, uID, err.Error())
	// 	infra.Log.Error(m)
	// 	l.leb.Publish(events.NewLogEvent(ctx, wID, uID, m, &pID, nil, models.UploadLogLevelError))
	// 	return
	// }

	// if err := r.Run(ctx); err != nil {
	// 	m := fmt.Sprintf("workflow run failed for processor %s and upload %s", processor.Name, uID)
	// 	infra.Log.Error(m, err)
	// 	return
	// }

	msg.Status = string(models.UploadStatusProcessingComplete)
	l.uploadEb.Publish(msg)
}
