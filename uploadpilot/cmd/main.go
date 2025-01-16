package main

import (
	"os"
	"os/signal"
	"syscall"

	initializer "github.com/uploadpilot/uploadpilot/internal/init"
)

func main() {
	// run InitServices() in a separate goroutine
	go initializer.InitServices()

	// wait for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
