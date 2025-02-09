package app

import (
	"log/slog"
	"os"
	"os/signal"
	"pinger/internal/api"
	"pinger/internal/config"
	"pinger/internal/kafka/producer"
	"pinger/internal/pinger"
	"syscall"
	"time"
)

func Run(config *config.Config) {
	var saver pinger.StatusSaver

	switch {
	case len(config.Brokers) > 0 && config.ContainerStatusTopic != "":
		contStatProducer := producer.NewContainerStatusProducer(config.ContainerStatusTopic, config.Brokers)
		defer contStatProducer.Close()
		saver = contStatProducer

	case config.BackendAPI != "":
		saver = api.NewBackendAPI(config.BackendAPI)

	default:
		slog.Error("kafka topic or backend api is not set")
		return
	}

	pingInterval := time.Duration(config.PingInterval) * time.Second
	containerPinger, err := pinger.NewPinger(saver, pingInterval, pingerOptionsFromConfig(config)...)
	if err != nil {
		slog.Error("error creating pinger", slog.Any("error", err))
		return
	}

	containerPinger.Run()

	// graceful shutdown
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	containerPinger.Stop()
}

func pingerOptionsFromConfig(config *config.Config) []pinger.Option {
	var opts []pinger.Option

	if len(config.Networks) > 0 {
		opts = append(opts, pinger.Networks(config.Networks))
	}

	if len(config.ComposeProjects) > 0 {
		opts = append(opts, pinger.ComposeProjects(config.ComposeProjects))
	}

	opts = append(opts, pinger.AllContainers(config.AllContainers))

	if config.ParsLabels {
		opts = append(opts, pinger.PingerLabel())
	}

	return opts
}
