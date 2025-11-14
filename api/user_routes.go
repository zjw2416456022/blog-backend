package api

import (
	"blog-backend/middleware"
	"blog-backend/controller"

	"github.com/gin-gonic/gin"
)

// setupUserRoutes 配置用户相关路由
func setupUserRoutes(api *gin.RouterGroup) {
	// 用户相关路由（无需认证）
	api.POST("/auth/register", controller.Register)
	api.POST("/auth/login", controller.Login)

	// 用户相关路由（需要认证）
	user := api.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/profile", controller.GetProfile)
	}
}