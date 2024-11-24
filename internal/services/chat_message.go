package services

import (
	"context"
	"distributed-chat-system/internal/apis/dtos"
	"distributed-chat-system/internal/constants"
	"distributed-chat-system/internal/models"
	"distributed-chat-system/internal/utils"
	"distributed-chat-system/pkg/kafka"
	"distributed-chat-system/pkg/redis"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ChatConsumerInterface interface {
	Notify(senderUserID string, message models.ChatMessage) error
}

type ChatMessageService struct {
	// Mutex to ensure thread-safe operations
	kafkaClient   *kafka.KafkaClient
	mutex         sync.RWMutex
	chatConsumers *ChatConsumerInterface
	redisRepo     redis.IRedisRepositories
}

func NewChatMessageService(kafkaClient *kafka.KafkaClient, redisRepo redis.IRedisRepositories) *ChatMessageService {
	return &ChatMessageService{
		kafkaClient:   kafkaClient,
		chatConsumers: nil,
		redisRepo:     redisRepo,
	}
}

// Starts listening to messages from Kafka consumer
func (s *ChatMessageService) StartMessageConsumption() {
	s.kafkaClient.ConsumeMessages(context.Background(), constants.ChatMessageTopic, s.consumeChatMessage)
}

// Handles a single chat message from from Kafka consumer & routes to chat consumers
func (s *ChatMessageService) consumeChatMessage(message string) {
	log.Println("Received chat message: ", message)
	// Unmarshal the message to models.ChatMessage
	var chatMessage *models.ChatMessage
	err := json.Unmarshal([]byte(message), &chatMessage)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Unmarshaled chat message: ", chatMessage)
	// Notify all registered chat consumers
	if s.chatConsumers != nil {
		(*s.chatConsumers).Notify(chatMessage.SenderUserID, *chatMessage)
	}
}

// SubscribeUserToChatServer adds a consumer user to service registry lookup store
func (s *ChatMessageService) SubscribeUserToChatServer(userId string) {
	jsonData := map[string]interface{}{
		"server_id": os.Getenv("SERVER_ID"),
	}

	jsonString, err := json.Marshal(jsonData)
	if err != nil {
		log.Printf("Error marshalling user to chat server")
		return
	}
	s.redisRepo.Set(userId, jsonString, time.Minute*5, context.Background())
	log.Println("User added to service registry lookup store: ", userId)
}

// SubscribeToChatMessage adds a consumer user to service registry lookup store
func (s *ChatMessageService) UnsubscribeUserToChatServer(userId string) {
	err := s.redisRepo.Del(userId, context.Background())
	if err != nil {
		return
	}
	log.Println("User removed from service registry lookup store: ", userId)
}

// LookupUserChatServer finds which server is the user currently connected to
func (s *ChatMessageService) LookupUserChatServer(userId string) *string {
	userLookupInfo, err := s.redisRepo.Get(userId, context.Background())
	if err != nil {
		return nil
	}

	var lookupData map[string]interface{}

	err = json.Unmarshal([]byte(userLookupInfo), &lookupData)
	if err != nil {
		return nil
	}
	return utils.StringPointer(lookupData["server_id"].(string))
}
func (s *ChatMessageService) SetChatConsumer(consumer ChatConsumerInterface) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.chatConsumers = &consumer
	log.Println("Consumer subscribed")
}

func (s *ChatMessageService) UnsetChatConsumer() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.chatConsumers = nil
}

// Publishes message to Kafka
func (s *ChatMessageService) SendMessageToUser(senderUserID string, message dtos.ChatMessageDto) error {
	// Here convert the message to string and publish to topic: chat-message
	chatMessage := &models.ChatMessage{
		EventID:        uuid.New().String(), // (Optional) For tracing purpose.
		SenderUserID:   senderUserID,
		ChatID:         message.ChatID,
		ReceiverUserID: message.ReceiverUserID,
		MessageType:    message.MessageType,
		Message:        message.Message,
	}

	messageJson, err := json.Marshal(chatMessage)
	if err != nil {
		return err
	}

	serverLookupId := s.LookupUserChatServer(chatMessage.ReceiverUserID)
	if serverLookupId == nil {
		return fmt.Errorf("user not connected to any server")
	}
	// Publish message to topic: chat-message
	err = s.kafkaClient.PublishMessage(*serverLookupId, string(messageJson))
	if err != nil {
		return err
	}
	log.Println("Message published successfully with event id", chatMessage.EventID)
	return nil
}
