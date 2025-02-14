package listeners

import (
	"time"

	"github.com/uploadpilot/uploadpilot/uploader/internal/svc"
)

func StartListeners(svc *svc.Services) {
	statusHandler := NewStatusListener(svc.UploadService)
	go statusHandler.Start()

	uploadLogsListener := NewUploadLogsListener(time.Second, 1000, svc.UploadService)
	go uploadLogsListener.Start()

	workflowListener := NewWorkflowListener(svc.ProcessorService)
	go workflowListener.Start()
}
