package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaClient struct {
	Producer *kafka.Producer
	Consumer *kafka.Consumer
	Config   *Config
}

// NewKafkaClient initializes a new Kafka client
func NewKafkaClient(cfg *Config) (*KafkaClient, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Brokers})
	if err != nil {
		return nil, err
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Brokers,
		"group.id":          cfg.GroupID,
		"auto.offset.reset": cfg.AutoOffset,
	})
	if err != nil {
		producer.Close()
		return nil, err
	}

	return &KafkaClient{
		Producer: producer,
		Consumer: consumer,
		Config:   cfg,
	}, nil
}

// PublishMessage sends a message to a Kafka topic
func (k *KafkaClient) PublishMessage(topic, message string) error {
	deliveryChan := make(chan kafka.Event, 1)

	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)
	if err != nil {
		return err
	}

	// Wait for delivery report
	e := <-deliveryChan
	msg := e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		log.Printf("Failed to deliver message: %v\n", msg.TopicPartition.Error)
	} else {
		log.Printf("Message delivered to topic %s [%d] at offset %v\n",
			*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	}
	close(deliveryChan)
	return nil
}

// ConsumeMessages starts consuming messages from Kafka
func (k *KafkaClient) ConsumeMessages(ctx context.Context, handler func(message string)) error {
	err := k.Consumer.SubscribeTopics(k.Config.Topics, nil)
	if err != nil {
		return err
	}

	go func() {
		defer k.Consumer.Close()
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping Kafka consumer...")
				return
			default:
				msg, err := k.Consumer.ReadMessage(-1)
				if err == nil {
					log.Printf("Received message: %s\n", string(msg.Value))
					handler(string(msg.Value))
				} else {
					log.Printf("Consumer error: %v\n", err)
				}
			}
		}
	}()
	return nil
}

// Close gracefully shuts down the Kafka producer and consumer
func (k *KafkaClient) Close() {
	if k.Producer != nil {
		k.Producer.Close()
	}
	if k.Consumer != nil {
		k.Consumer.Close()
	}
}
