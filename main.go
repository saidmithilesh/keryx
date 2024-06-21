package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"keryx/comms"
	"keryx/connections"
	"keryx/utils"

	"go.uber.org/zap"
)

func handleInterrupt(
	sigChan chan os.Signal,
	hub *connections.Hub,
	logger *zap.Logger,
) {
	sig := <-sigChan
	logger.Info("interrupt received", zap.String("signal", sig.String()))

	hub.Stop()

	time.Sleep(5 * time.Second)
	logger.Info("server shutting down")
	os.Exit(0)
}

func increaseResources() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
}

func main() {
	increaseResources()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	config := utils.NewConfig()
	logger := utils.NewLogger()

	logger.Info("logger initialised", zap.String("environment", config.Env))
	router := comms.NewRouter(config, logger)
	hub := connections.NewHub(config, logger, router)

	go handleInterrupt(sigChan, hub, logger)
	hub.Start()
}
