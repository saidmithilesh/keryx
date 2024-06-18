package main

import (
	"keryx/connections"
	"keryx/utils"

	"go.uber.org/zap"
)

func main() {
	config := utils.NewConfig()
	logger := utils.NewLogger()

	logger.Info("logger initialised", zap.String("environment", config.Env))
	hub := connections.NewHub(config, logger)
	hub.Start()
}
