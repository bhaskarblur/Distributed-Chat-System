package handlers

import (
	"distributed-chat-system/internal/apis/dtos"
	"distributed-chat-system/internal/models"
	"distributed-chat-system/internal/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader    websocket.Upgrader
	conns       map[string]*websocket.Conn
	chatService *services.ChatMessageService
}

// InitWebSocketHandler initializes the WebSocketHandler and subscribes it to the ChatMessageService
func InitWebSocketHandler(chatService *services.ChatMessageService) *WebSocketHandler {

	handler := &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for simplicity, modify for production
				return true
			},
		},
		conns:       make(map[string]*websocket.Conn),
		chatService: chatService,
	}

	// Subscribe to the ChatMessageService once
	chatService.SubscribeToChatMessage(handler)
	return handler
}

// InitWebSocket handles WebSocket connections and communication
func (h *WebSocketHandler) InitWebSocket(c *gin.Context) {
	// Extract user_id from path params
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Upgrade the connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	// Store the connection
	h.conns[userID] = conn
	defer func() {
		conn.Close()
		delete(h.conns, userID)
		log.Printf("WebSocket connection closed for user: %s", userID)
	}()

	log.Printf("WebSocket connection established for user: %s", userID)

	// WebSocket communication loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Parse the received JSON message
		var chatMessage dtos.ChatMessageDto
		err = json.Unmarshal(message, &chatMessage)
		if err != nil {
			log.Println("Invalid message format:", err)
			conn.WriteMessage(websocket.TextMessage, []byte(`{"error": "Invalid message format"}`))
			continue
		}

		// Log and send the message to the service
		log.Printf("Message received from user %s: %+v", userID, chatMessage)
		err = h.chatService.SendMessageToUser(userID, chatMessage)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}

// Notify sends a message to the connected WebSocket user
func (h *WebSocketHandler) Notify(senderUserID string, message models.ChatMessage) error {
	log.Println("Got Notified with event id:", message.EventID)
	conn, exists := h.conns[message.ReceiverUserID]
	if !exists {
		log.Printf("No active WebSocket connection for receiver_user: %s. Skipping ahead.", message.ReceiverUserID)
		return nil
	}

	response := gin.H{
		"chat_id":      message.ChatID,
		"sender":       senderUserID,
		"message_type": message.MessageType,
		"message":      message.Message,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling message for user %s: %v", message.ReceiverUserID, err)
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, responseJSON)
	if err != nil {
		log.Printf("Error sending message to user %s: %v", message.ReceiverUserID, err)
		return err
	}

	log.Printf("Message sent to user %s: %+v", message.ReceiverUserID, message)
	return nil
}
