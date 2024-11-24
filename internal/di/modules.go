package di

import (
	"distributed-chat-system/internal/apis/handlers"
	"distributed-chat-system/internal/constants"
	"distributed-chat-system/internal/services"
	"distributed-chat-system/pkg/kafka"
	"log"
	"os"

	"go.uber.org/dig"
)

var Container *dig.Container

// InitializeDependencies sets up the dependency injection container
func InitializeDependencies() {
	if Container == nil {
		Container = dig.New()
	}

	// Provide KafkaClient
	err := Container.Provide(func() (*kafka.KafkaClient, error) {
		cfg := kafka.DefaultConfig()
		cfg.Brokers = os.Getenv("KAFKA_URL") // Update as per your environment
		cfg.Topics = []string{constants.ChatMessageTopic}
		cfg.GroupID = os.Getenv("CHAT_GROUP_ID")
		return kafka.NewKafkaClient(cfg)
	})
	if err != nil {
		log.Fatalf("Failed to provide KafkaClient: %v", err)
	}

	// Provide ChatMessageService
	err = Container.Provide(func(kafkaClient *kafka.KafkaClient) *services.ChatMessageService {
		service := services.NewChatMessageService(kafkaClient)
		service.StartMessageConsumption()
		return service
	})
	if err != nil {
		log.Fatalf("Failed to provide ChatMessageService: %v", err)
	}

	// Provide WebSocketHandler
	err = Container.Provide(func(chatService *services.ChatMessageService) *handlers.WebSocketHandler {
		return handlers.InitWebSocketHandler(chatService)
	})
	if err != nil {
		log.Fatalf("Failed to provide WebSocketHandler: %v", err)
	}
}

// Resolve resolves a dependency from the container
func Resolve(target interface{}) {
	if Container == nil {
		log.Fatal("Dependencies are not initialized. Call InitializeDependencies first.")
	}

	err := Container.Invoke(target)
	if err != nil {
		log.Fatalf("Failed to resolve dependency: %v", err)
	}
}
