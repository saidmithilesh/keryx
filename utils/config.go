package utils

import (
	"os"
	"strconv"
	"time"
)

// Constants representing different environment types.
const (
	EnvProd  = "production"  // Production environment
	EnvDev   = "development" // Development environment
	EnvDebug = "debug"       // Debug environment

	DefaultEnv              = EnvProd // Default environment set to production
	DefaultPort             = "8080"  // Default port set to 8080
	DefaultRegistryInterval = 30 * time.Second
)

// Config holds the configuration settings for the application.
type Config struct {
	Env              string        // Environment type
	Port             string        // Port number
	RegistryInterval time.Duration // interval of updating registry
}

// hasValidEnv checks if the environment setting is valid.
func (cfg *Config) hasValidEnv() bool {
	switch cfg.Env {
	case EnvProd, EnvDebug, EnvDev:
		return true // Valid environment
	default:
		return false // Invalid environment
	}
}

// hasValidPort checks if the port setting is a valid integer.
func (cfg *Config) hasValidPort() bool {
	_, err := strconv.Atoi(cfg.Port)
	return err != nil
}

// setDefaults sets the default environment and port
// if the current settings are invalid.
func (cfg *Config) setDefaults() error {
	if !cfg.hasValidEnv() {
		cfg.Env = DefaultEnv // Set default environment
	}

	if !cfg.hasValidPort() {
		cfg.Port = DefaultPort // Set default port
	}

	return nil
}

var config *Config

// New initializes the configuration settings by reading environment variables.
// Uses the singleton pattern to ensure only one instance of Config is created.
func NewConfig() *Config {
	cfgOnce.Do(func() {
		var err error
		config = &Config{}

		config.Env = os.Getenv("ENV")   // Get environment setting from ENV variable
		config.Port = os.Getenv("PORT") // Get port setting from PORT variable
		config.RegistryInterval, err = time.ParseDuration(os.Getenv("HUB_REGISTRY_INTERVAL"))
		if err != nil {
			config.RegistryInterval = DefaultRegistryInterval
		}
		config.setDefaults() // Set default values if necessary
	})

	return config
}
