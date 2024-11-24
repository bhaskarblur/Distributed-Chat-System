package di

import (
	"distributed-chat-system/internal/apis/handlers"
	"distributed-chat-system/internal/constants"
	"distributed-chat-system/internal/services"
	"distributed-chat-system/pkg/kafka"
	"distributed-chat-system/pkg/redis"
	"log"
	"os"

	redisClient "github.com/redis/go-redis/v9"
	"go.uber.org/dig"
)

var Container *dig.Container

// InitializeDependencies sets up the dependency injection container
func InitializeDependencies() {
	if Container == nil {
		Container = dig.New()
	}

	redisRepo := initRedis()
	log.Println("Kafka url:", os.Getenv("KAFKA_URL"))
	// Provide KafkaClient
	err := Container.Provide(func() (*kafka.KafkaClient, error) {
		cfg := kafka.DefaultConfig()
		cfg.Brokers = os.Getenv("KAFKA_URL") // Update as per your environment
		cfg.Topics = []string{constants.ChatMessageTopic}
		cfg.GroupID = os.Getenv("CHAT_GROUP_ID")
		return kafka.NewKafkaClient(cfg), nil
	})
	if err != nil {
		log.Fatalf("Failed to provide KafkaClient: %v", err)
	}

	// Provide ChatMessageService
	err = Container.Provide(func(kafkaClient *kafka.KafkaClient) *services.ChatMessageService {
		service := services.NewChatMessageService(kafkaClient, redisRepo)
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

func initRedis() redis.IRedisRepositories {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUsername := os.Getenv("REDIS_USERNAME")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Provide Redis client
	client, err := redis.RedisClient(redisHost, redisPort, redisUsername, redisPassword)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	// Provide Redis client
	Container.Provide(func() (*redisClient.Client, error) {
		return client, nil
	})

	redis_repo := redis.NewRedisRepositories(client)
	// Provides Redis CRUD Repository
	Container.Provide(func(redisClient *redisClient.Client) redis.IRedisRepositories {
		return redis_repo
	})

	return redis_repo
}
