package services

import (
	"context"
	"distributed-chat-system/internal/apis/dtos"
	"distributed-chat-system/internal/constants"
	"distributed-chat-system/internal/models"
	"distributed-chat-system/pkg/kafka"
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
)

type ChatConsumerInterface interface {
	Notify(senderUserID string, message models.ChatMessage) error
}

type ChatMessageService struct {
	// Mutex to ensure thread-safe operations
	kafkaClient   *kafka.KafkaClient
	mutex         sync.RWMutex
	chatConsumers []ChatConsumerInterface
}

func NewChatMessageService(kafkaClient *kafka.KafkaClient) *ChatMessageService {
	return &ChatMessageService{
		kafkaClient:   kafkaClient,
		chatConsumers: make([]ChatConsumerInterface, 0),
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
	for _, consumer := range s.chatConsumers {
		consumer.Notify(chatMessage.SenderUserID, *chatMessage)
	}
}

// SubscribeToChatMessage adds a consumer to the list
func (s *ChatMessageService) SubscribeToChatMessage(consumer ChatConsumerInterface) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.chatConsumers = append(s.chatConsumers, consumer)
	log.Println("New consumer subscribed")
}

// UnsubscribeToChatMessage removes a consumer from the list
func (s *ChatMessageService) UnsubscribeToChatMessage(consumer ChatConsumerInterface) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, c := range s.chatConsumers {
		if c == consumer {
			s.chatConsumers = append(s.chatConsumers[:i], s.chatConsumers[i+1:]...)
			log.Println("Consumer unsubscribed")
			break
		}
	}
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

	// Publish message to topic: chat-message
	err = s.kafkaClient.PublishMessage(chatMessage.ReceiverUserID, string(messageJson))
	if err != nil {
		return err
	}
	log.Println("Message published successfully with event id", chatMessage.EventID)
	return nil
}
