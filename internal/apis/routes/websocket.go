package routes

import (
	"distributed-chat-system/internal/apis/handlers"

	"github.com/gin-gonic/gin"
)

// SetupWebSocketRoutes sets up WebSocket routes
func SetupWebSocket(router *gin.RouterGroup) {
	socketHandler := handlers.InitWebSocketHandler()
	router.GET("/ws/user/:user_id", socketHandler.InitWebSocket)
}
