package kafka

import (
	"backend/pkg/logger/sl"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func CreateTopic(brokers []string, topic string, numPartitions, replicationFactor int) error {
	// Create a Kafka Admin client
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer func(conn *kafka.Conn) {
		if err = conn.Close(); err != nil {
			slog.Error("Failed to close Kafka connection", sl.Err(err))
		}
	}(conn)

	// Check if the topic exists
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	for i := range partitions {
		p := partitions[i]
		if p.Topic == topic {
			// Topic already exists
			return nil
		}
	}

	// Create the topic if it does not exist
	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get Kafka controller: %w", err)
	}

	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))

	adminConn, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka controller: %w", err)
	}
	defer func(adminConn *kafka.Conn) {
		if err = adminConn.Close(); err != nil {
			slog.Error("failed to close Kafka admin connection", sl.Err(err))
		}
	}(adminConn)

	err = adminConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}
