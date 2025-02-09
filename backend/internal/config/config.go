package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string   `env:"PORT"                      env-default:"80"`
	DBUrl                  string   `env:"DB_URL"                                       env-required:"true"`
	ContainerStatusTopic   string   `env:"CONTAINER_STATUS_TOPIC"                       env-required:"true"`
	ContainerStatusGroupID string   `env:"CONTAINER_STATUS_GROUP_ID"                    env-required:"true"`
	Brokers                []string `env:"KAFKA_BROKERS"                                env-required:"true"`
	Production             bool     `env:"PRODUCTION"                env-default:"true"`
}

var (
	config Config    //nolint:gochecknoglobals,lll // Global config is initialized once and accessed throughout the application.
	once   sync.Once //nolint:gochecknoglobals,lll // Ensures the config is initialized only once, which requires a global sync.Once.
)

func Get() *Config {
	once.Do(func() {
		err := godotenv.Load()

		if err != nil {
			log.Println("error loading .env file")
		}
		err = cleanenv.ReadEnv(&config)
		if err != nil {
			panic(fmt.Sprintf("Failed to get config: %s", err))
		}
	})

	return &config
}
