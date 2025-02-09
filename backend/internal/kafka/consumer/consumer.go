package consumer

import (
	"backend/internal/model"
	"backend/internal/service"
	"backend/pkg/logger/sl"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"

	kafkainternal "backend/internal/kafka"
)

const maxBytesRead = 10e6 // 10MB
const maxWaitRead = 10 * time.Millisecond
const minBytesRead = 1

type ContainerStatusConsumer struct {
	Topic   string
	GroupID string
	reader  *kafka.Reader
	service *service.ContainerStatusService
	brokers []string
}

func NewContainerStatusConsumer(topic, groupID string,
	containerStatusService *service.ContainerStatusService,
	brokers []string) (*ContainerStatusConsumer, error) {
	err := kafkainternal.CreateTopic(brokers, topic, 1, 1)
	if err != nil {
		return nil, err
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MaxBytes: maxBytesRead,
		MaxWait:  maxWaitRead,
		MinBytes: minBytesRead,
	})

	return &ContainerStatusConsumer{
		Topic:   topic,
		GroupID: groupID,
		reader:  reader,
		service: containerStatusService,
		brokers: brokers,
	}, nil
}

func (c *ContainerStatusConsumer) Start() {
	slog.Info("start consumer",
		"topic", c.Topic,
		"group_id", c.GroupID,
		"brokers", c.brokers)
	go func() {
		for {
			msg, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				slog.Error("ReadMessage",
					sl.Err(err),
					slog.String("topic", c.Topic),
					slog.String("group_id", c.GroupID),
					slog.Any("brokers", c.brokers))
				continue
			}

			var data model.ContainerStatus
			err = json.Unmarshal(msg.Value, &data)
			if err != nil {
				slog.Error("Unmarshal",
					sl.Err(err),
					slog.String("topic", c.Topic),
					slog.String("group_id", c.GroupID),
					slog.Any("brokers", c.brokers))
				continue
			}

			err = c.service.UpsertContainerStatus(&data)
			if err != nil {
				slog.Error("UpsertContainerStatus",
					sl.Err(err),
					slog.String("topic", c.Topic),
					slog.String("group_id", c.GroupID),
					slog.Any("brokers", c.brokers))
			}
		}
	}()
}

func (c *ContainerStatusConsumer) Stop() error {
	err := c.reader.Close()
	if err != nil {
		return err
	}
	return nil
}
