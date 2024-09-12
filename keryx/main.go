package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/saidmithilesh/keryx/config"
	"github.com/saidmithilesh/keryx/logging"
	"github.com/saidmithilesh/keryx/server"

	"go.uber.org/zap"
)

func interruptHandler(signalChan chan os.Signal, logger *zap.Logger) {
	sig := <-signalChan
	logger.Info("Received signal", zap.String("signal", sig.String()))
	logger.Info("Shutting down the application")
	logger.Sync()
	os.Exit(0)
}

func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger(cfg)

	// interrupt handling mechanism for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go interruptHandler(signalChan, logger)

	httpServer := server.New(cfg, logger)
	httpServer.Run()
}
