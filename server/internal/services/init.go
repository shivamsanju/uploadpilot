package services

import (
	"github.com/shivamsanju/uploader/internal/config"
)

func InitServices() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Initialize the logger
	initLogger()

	// Initialize MongoDB
	err = initMongoDB(cfg)
	if err != nil {
		panic(err)
	}

	// Initialize Supertokens.
	err = initSuperTokens(cfg)
	if err != nil {
		panic(err)
	}

	// Initialize the web server.
	err = initWebServer(cfg)
	if err != nil {
		panic(err)
	}
}
