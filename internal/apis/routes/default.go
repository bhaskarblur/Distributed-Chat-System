package routes

import (
	"distributed-chat-system/internal/apis/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	router.GET("/", handlers.HealthCheck)

	wsGroup := router.Group("/ws")
	SetupWebSocket(wsGroup)
}
