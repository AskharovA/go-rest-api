package config

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SecretKey string `envconfig:"SECRET_KEY" required:"true"`
	AppPort   int    `envconfig:"APP_PORT" default:"8080"`
}

var (
	once   sync.Once
	config *Config
)

func Get() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found.")
		}

		var cfg Config
		err = envconfig.Process("", &cfg)
		if err != nil {
			log.Fatalf("failed to process config: %v", err)
		}

		config = &cfg
	})
	return config
}
