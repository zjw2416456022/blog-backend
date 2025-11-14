package api

import (
	"blog-backend/middleware"
	"blog-backend/controller"

	"github.com/gin-gonic/gin"
)

// setupCommentRoutes 配置评论相关路由
func setupCommentRoutes(api *gin.RouterGroup) {
	// 评论相关路由
	comments := api.Group("/posts/:id/comments")
	{
		// 获取评论列表（无需认证）
		comments.GET("", controller.GetComments)

		// 创建评论（需要认证）
		comments.POST("", middleware.AuthMiddleware(), controller.CreateComment)
	}

	// 更新评论（需要认证）
	api.PUT("/comments/:id", middleware.AuthMiddleware(), controller.UpdateComment)

	// 删除评论（需要认证）
	api.DELETE("/comments/:id", middleware.AuthMiddleware(), controller.DeleteComment)
}