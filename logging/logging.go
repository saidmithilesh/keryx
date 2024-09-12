package logging

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/saidmithilesh/keryx/config"
)

var once sync.Once
var logger *zap.Logger

// buildProdctionLogger creates a new production logger
// with JSON encoding and log to stderr.
func buildProdctionLogger() *zap.Logger {
	// NewProductionEncoderConfig returns an opinionated EncoderConfig for production environments.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"                   // Change the default time key to timestamp
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder // Use ISO8601 time format for timestamp

	// Get the hostname of the machine
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel), // Set the lowest log level to Info
		Development:       false,                               // Disable development mode
		DisableCaller:     false,                               // Enable caller information
		DisableStacktrace: false,                               // Enable stacktrace information for errors and panics
		Sampling:          nil,                                 // Disable sampling
		Encoding:          "json",                              // Use JSON encoding for logs in production
		EncoderConfig:     encoderCfg,                          // Use the custom encoder configuration
		OutputPaths: []string{
			"stdout", // Log to stdout in production
		},
		ErrorOutputPaths: []string{
			"stderr", // Log errors to stderr in production
		},
		InitialFields: map[string]interface{}{
			"pid":      os.Getpid(), // Log the process ID of the application as pid
			"hostname": hostname,    // Log the hostname of the machine as hostname
		},
	}

	return zap.Must(config.Build()) // Build the logger and panic if there is an error
}

// GetLogger returns a singleton instance of the logger.
func GetLogger(cfg *config.Config) *zap.Logger {
	once.Do(func() {
		if cfg.IsProduction() {
			// Create a new production logger
			// if the application environment is production
			// and store it in the logger variable
			logger = buildProdctionLogger()
		} else {
			// Create a new development logger
			// if the application environment is not production
			// and store it in the logger variable
			logger = zap.Must(zap.NewDevelopment())
		}
	})

	return logger
}
