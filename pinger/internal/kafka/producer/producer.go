package producer

import (
	"context"
	"encoding/json"
	"pinger/internal/model"
	"time"

	"github.com/segmentio/kafka-go"
)

type ContainerStatusProducer struct {
	writer  *kafka.Writer
	Topic   string
	brokers []string
}

const writeTimeout = 5 * time.Second
const batchSize = 1

func NewContainerStatusProducer(topic string, brokers []string) *ContainerStatusProducer {
	writer := kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		WriteTimeout: writeTimeout,
		RequiredAcks: kafka.RequireOne,
		BatchSize:    batchSize,
	}
	return &ContainerStatusProducer{
		Topic:   topic,
		brokers: brokers,
		writer:  &writer,
	}
}

func (c *ContainerStatusProducer) Close() {
	err := c.writer.Close()
	if err != nil {
		return
	}
}

func (c *ContainerStatusProducer) SaveContainerStatus(containerStatus model.ContainerStatus) error {
	message, err := json.Marshal(containerStatus)
	if err != nil {
		return err
	}
	err = c.writer.WriteMessages(context.Background(), kafka.Message{
		Value: message,
	})
	if err != nil {
		return err
	}

	return nil
}
