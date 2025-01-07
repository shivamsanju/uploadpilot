package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/shivamsanju/uploader/internal/services"
)

func main() {
	// run InitServices() in a separate goroutine
	go services.InitServices()

	// wait for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
