package main

import (
	"github.com/saidmithilesh/keryx/config"
	"github.com/saidmithilesh/keryx/logging"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger(cfg)

	logger.Info("Starting the application")
	logger.Info("Shutting down the application")
}
