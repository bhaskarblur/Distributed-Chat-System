package routes

import (
	"distributed-chat-system/internal/apis/handlers"
	"distributed-chat-system/internal/di"

	"log"

	"github.com/gin-gonic/gin"
)

// SetupWebSocketRoutes sets up WebSocket routes
func SetupWebSocket(router *gin.RouterGroup) {
	// Resolve the socketHandler from the DI container
	var socketHandler *handlers.WebSocketHandler
	err := di.Container.Invoke(func(h *handlers.WebSocketHandler) {
		socketHandler = h
	})
	if err != nil {
		log.Fatalf("Failed to resolve WebSocketHandler: %v", err)
	}

	// Define the WebSocket route
	router.GET("/ws/user/:user_id", socketHandler.InitWebSocket)
}
