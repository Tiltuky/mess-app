package kafkaConn

import (
	"context"
	"fmt"
	"log"
	"progekt/dating-app/geolocation-service/internal/modules/geoservice/service"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) (service.MessageQueuer, error) {
	// Проверка существования топика
	err := createTopicIfNotExists(brokers, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &KafkaProducer{writer: writer}, nil
}

func createTopicIfNotExists(brokers []string, topic string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		// Если топик не существует, создаем его
		controller, err := conn.Controller()
		if err != nil {
			return fmt.Errorf("failed to get controller: %w", err)
		}

		controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
		if err != nil {
			return fmt.Errorf("failed to dial controller: %w", err)
		}
		defer controllerConn.Close()

		topicConfig := kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}

		err = controllerConn.CreateTopics(topicConfig)
		if err != nil {
			return fmt.Errorf("failed to create topic: %w", err)
		}
	} else {
		log.Printf("Topic %s already exists with %d partitions", topic, len(partitions))
	}

	return nil
}

func (kp *KafkaProducer) Close() error {
	err := kp.writer.Close()
	return err
}

func (k *KafkaProducer) Publish(topic string, message []byte) error {
	msg := kafka.Message{

		Value: message,
	}
	err := k.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to publish message to Kafka: %v", err)
	} else {
		log.Printf("Message published to Kafka: %s", message)
	}
	return err
}
