package api

import (
	"blog-backend/middleware"
	"blog-backend/controller"

	"github.com/gin-gonic/gin"
)

// setupPostRoutes 配置文章相关路由
func setupPostRoutes(api *gin.RouterGroup) {
	// 文章相关路由
	posts := api.Group("/posts")
	{
		// 获取文章列表和详情（无需认证）
		posts.GET("", middleware.OptionalAuthMiddleware(), controller.GetPosts)
		posts.GET("/:id", middleware.OptionalAuthMiddleware(), controller.GetPost)

		// 创建、更新、删除文章（需要认证）
		authPosts := posts.Group("/")
		authPosts.Use(middleware.AuthMiddleware())
		{
			authPosts.POST("", controller.CreatePost)
			authPosts.PUT("/:id", controller.UpdatePost)
			authPosts.DELETE("/:id", controller.DeletePost)
		}
	}
}