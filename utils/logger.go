package utils

import (
	"log" // Standard library package for logging

	"go.uber.org/zap" // Zap package for structured logging
)

// Global variables for logger instance and initialization error
var logger *zap.Logger
var logInitErr error

// NewLogger initializes and returns a zap.Logger instance.
func NewLogger() *zap.Logger {
	logOnce.Do(func() { // Ensure this block runs only once
		// Initialize logger based on the environment
		if config.Env == EnvProd {
			logger, logInitErr = zap.NewProduction() // Production logger
		} else {
			logger, logInitErr = zap.NewDevelopment() // Development logger
		}

		// Check for errors during logger initialization
		if logInitErr != nil {
			// Log fatal error and exit if logger initialization fails
			log.Fatalf(
				"Failed to initialise structured logging %#v",
				logInitErr,
			)
		}
	})

	return logger // Return the initialized logger
}
