package listeners

import (
	"time"
)

func StartListeners() {
	statusHandler := NewStatusListener()
	go statusHandler.Start()

	uploadLogsListener := NewUploadLogsListener(time.Second, 1000)
	go uploadLogsListener.Start()

	procListener := NewProcListener()
	go procListener.Start()
}
