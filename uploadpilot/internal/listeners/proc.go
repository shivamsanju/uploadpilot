package listeners

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/uploadpilot/uploadpilot/internal/db"
	"github.com/uploadpilot/uploadpilot/internal/db/models"
	"github.com/uploadpilot/uploadpilot/internal/events"
	"github.com/uploadpilot/uploadpilot/internal/infra"
	"github.com/uploadpilot/uploadpilot/internal/proc"
	"github.com/uploadpilot/uploadpilot/internal/utils"
)

type ProcListener struct {
	eventChan chan events.UploadEvent
	pRepo     *db.ProcessorRepo
	upRepo    *db.UploadRepo
	leb       *events.LogEventBus
	ueb       *events.UploadEventBus
}

func NewProcListener() *ProcListener {
	eventBus := events.GetUploadEventBus()

	eventChan := make(chan events.UploadEvent)
	eventBus.Subscribe(events.EventUploadComplete, eventChan)

	return &ProcListener{
		eventChan: eventChan,
		pRepo:     db.NewProcessorRepo(),
		upRepo:    db.NewUploadRepo(),
		leb:       events.GetLogEventBus(),
		ueb:       eventBus,
	}
}

func (l *ProcListener) Start() {
	defer utils.Recover()
	infra.Log.Info("starting processing listener...")
	for event := range l.eventChan {
		infra.Log.Infof("processing upload complete event %s", event.Key)
		go l.startProcessing(event)
	}
}

func (l *ProcListener) startProcessing(event events.UploadEvent) {
	defer utils.Recover()
	processors, err := l.pRepo.GetAll(event.Context, event.Upload.WorkspaceID.Hex())
	if err != nil {
		infra.Log.Errorf("failed to get processors: %s", err.Error())
		return
	}

	if len(processors) == 0 {
		infra.Log.Info("no processors found, skipping processing...")
		return
	}

	wg := &sync.WaitGroup{}
	for _, processor := range processors {
		wg.Add(1)
		go l.startSingleProcessor(wg, event, &processor)
	}
	wg.Wait()
}

func (l *ProcListener) startSingleProcessor(wg *sync.WaitGroup, event events.UploadEvent, processor *models.Processor) {
	defer utils.Recover()
	defer wg.Done()

	wID := event.Upload.WorkspaceID.Hex()
	uID := event.Upload.ID.Hex()
	pID := processor.ID.Hex()
	ctx := context.Background()
	upload, err := l.upRepo.Get(ctx, uID)
	if err != nil {
		l.leb.Publish(events.NewLogEvent(ctx, wID, uID, err.Error(), nil, nil, models.UploadLogLevelError))
		return
	}

	p, err := l.pRepo.Get(ctx, pID)
	if err != nil {
		l.leb.Publish(events.NewLogEvent(ctx, wID, uID, err.Error(), nil, nil, models.UploadLogLevelError))
		return
	}

	fileType := upload.Metadata["filetype"].(string)
	if len(p.Triggers) > 0 && !slices.Contains(p.Triggers, fileType) {
		return
	}

	m := fmt.Sprintf("started processing upload for processor '%s'", processor.Name)
	l.ueb.Publish(events.NewUploadEvent(event.Context, events.EventUploadProcessing, &event.Upload, "", nil))
	l.leb.Publish(events.NewLogEvent(event.Context, wID, uID, m, nil, nil, models.UploadLogLevelInfo))

	r := proc.NewProcWorkflowRunner()
	if err := r.Build(event.Context, wID, pID, uID); err != nil {
		m := fmt.Sprintf("failed to build workflow for processor %s and upload %s: %s", processor.Name, uID, err.Error())
		infra.Log.Error(m)
		l.leb.Publish(events.NewLogEvent(event.Context, wID, uID, m, &pID, nil, models.UploadLogLevelError))
		return
	}

	if err := r.Run(event.Context); err != nil {
		m := fmt.Sprintf("workflow run failed for processor %s and upload %s", processor.Name, uID)
		infra.Log.Error(m, err)
		return
	}
}
