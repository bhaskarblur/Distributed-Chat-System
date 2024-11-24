package services

import (
	"distributed-chat-system/internal/apis/dtos"
	"log"
	"sync"
)

type ChatConsumerInterface interface {
	Notify(senderUserID string, message dtos.ChatMessageDto) error
}

type ChatMessageService struct {
	// Mutex to ensure thread-safe operations
	mutex         sync.RWMutex
	chatConsumers []ChatConsumerInterface
}

func NewChatMessageService() *ChatMessageService {
	//Connect to Kafka Consumr here!
	return &ChatMessageService{
		chatConsumers: make([]ChatConsumerInterface, 0),
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

func (s *ChatMessageService) SendMessageToUser(senderUserID string, message dtos.ChatMessageDto) error {
	// Here convert the message to string and publish to topic: chat-message
}
