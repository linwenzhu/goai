package router

import (
	"github.com/gin-gonic/gin"
	"goAi/internal/handler"
	"goAi/internal/middleware"
	"goAi/pkg/response"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		response.Success(c, "pong")
	})

	// 用户相关路由
	userHandler := handler.NewUserHandler()
	chatHandler := handler.NewChatHandler()

	v1 := r.Group("/api/v1")
	{
		// 不需要鉴权的路由
		user := v1.Group("/user")
		{
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)
		}

		// 需要鉴权的路由，后续 AI 对话接口加在这里
		auth := v1.Group("/")
		auth.Use(middleware.Auth())
		{
			// 测试鉴权是否生效
			auth.GET("/profile", func(c *gin.Context) {
				userID, _ := c.Get("userID")
				username, _ := c.Get("username")
				response.Success(c, gin.H{
					"user_id":  userID,
					"username": username,
				})
			})
			// AI 对话
			auth.POST("/chat", chatHandler.Chat)
			auth.POST("/chat/stream", chatHandler.ChatStream)
		}

		return r
	}
}
