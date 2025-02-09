package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	BackendAPI           string   `env:"BACKEND_API"            env-default:"http://localhost:8080"`
	ContainerStatusTopic string   `env:"CONTAINER_STATUS_TOPIC"`
	Networks             []string `env:"NETWORKS"`
	Brokers              []string `env:"KAFKA_BROKERS"`
	ComposeProjects      []string `env:"COMPOSE_PROJECTS"`
	AllContainers        bool     `env:"ALL_CONTAINERS"         env-default:"true"`
	ParsLabels           bool     `env:"LABELS"                 env-default:"true"`
	PingInterval         int      `env:"PING_INTERVAL"          env-default:"10"`
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
