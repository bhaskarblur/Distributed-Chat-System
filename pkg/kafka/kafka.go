package kafka

import (
	"context"
	"distributed-chat-system/internal/constants"
	"log"
	"strings"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	Producer     *kafka.Writer
	Consumer     *kafka.Reader
	Config       *Config
	DefaultTopic string
}

// NewKafkaClient initializes a new Kafka client
func NewKafkaClient(cfg *Config) *KafkaClient {
	// Convert comma-separated brokers string to a slice
	brokers := strings.Split(cfg.Brokers, ",")

	// Initialize Producer
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaClient{
		Producer:     producer,
		Config:       cfg,
		DefaultTopic: constants.ChatMessageTopic,
	}
}

// PublishMessage sends a message to a Kafka topic
func (k *KafkaClient) PublishMessage(receiverID, message string) error {
	err := k.Producer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: k.DefaultTopic, // Default topic is used
			Key:   []byte(receiverID),
			Value: []byte(message),
		},
	)
	if err != nil {
		log.Printf("Failed to deliver message to topic %s: %v\n", k.DefaultTopic, err)
		return err
	}

	log.Printf("Message delivered to topic %s\n", k.DefaultTopic)
	return nil
}

// ConsumeMessages starts consuming messages from a specified Kafka topic at runtime
func (k *KafkaClient) ConsumeMessages(ctx context.Context, topic string, handler func(message string)) error {
	// Convert comma-separated brokers string to a slice
	brokers := strings.Split(k.Config.Brokers, ",")

	log.Println("Kafka Consumer brokers: ", brokers)
	// Initialize Consumer dynamically for the topic
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     k.Config.GroupID,
		Topic:       topic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
	})

	go func() {
		defer consumer.Close()
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping Kafka consumer...")
				return
			default:
				m, err := consumer.ReadMessage(ctx)
				if err != nil {
					log.Printf("Consumer error for topic %s: %v\n", topic, err)
					continue
				}
				log.Printf("Received message from topic %s: %s\n", m.Topic, string(m.Value))
				handler(string(m.Value))
			}
		}
	}()
	return nil
}

// Close gracefully shuts down the Kafka producer
func (k *KafkaClient) Close() {
	if k.Producer != nil {
		if err := k.Producer.Close(); err != nil {
			log.Printf("Failed to close Kafka producer: %v\n", err)
		}
	}
}
