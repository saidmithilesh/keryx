package conf

import (
	"log"
	"os"
	"sync"
)

type Config struct {
	Port  string
	DBURL string
}

var config Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		config.Port = os.Getenv("PORT")
		config.DBURL = os.Getenv("DB_URL")

		log.Println("Config loaded from environment variables")
	})

	return &config
}
