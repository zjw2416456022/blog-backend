package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有API路由
func SetupRoutes(router *gin.Engine) {
	// API路由组
	api := router.Group("/api/v1")
	{
		// 设置用户相关路由
		setupUserRoutes(api)
		
		// 设置文章相关路由
		setupPostRoutes(api)
		
		// 设置评论相关路由
		setupCommentRoutes(api)
	}
}