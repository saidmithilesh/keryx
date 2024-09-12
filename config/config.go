package config

import (
	"os"
	"sync"
)

var once sync.Once

const (
	EnvProduction = "production"
	EnvDebug      = "debug"
	EnvTest       = "test"
)

type Config struct {
	AppPort string
	AppEnv  string
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == EnvProduction
}

func (c *Config) validate() {
	if c.AppPort == "" {
		c.AppPort = DefaultAppPort
	}
	if c.AppEnv == "" {
		c.AppEnv = DefaultAppEnv
	}
}

func buildConfig() {
	config = &Config{
		AppPort: os.Getenv("APP_PORT"),
		AppEnv:  os.Getenv("APP_ENV"),
	}
	config.validate()
}

var config *Config

func GetConfig() *Config {
	once.Do(buildConfig)
	return config
}
