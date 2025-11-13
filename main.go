package main

import (
	"blog-backend/config"
	"blog-backend/middleware"
	"blog-backend/routes"
	"blog-backend/utils"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	utils.InitLogger()
	utils.Info("Starting blog backend server...")

	// 初始化数据库
	config.InitDB()

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户相关路由（无需认证）
		api.POST("/auth/register", routes.Register)
		api.POST("/auth/login", routes.Login)

		// 用户相关路由（需要认证）
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/profile", routes.GetProfile)
		}

		// 文章相关路由
		posts := api.Group("/posts")
		{
			// 获取文章列表和详情（无需认证）
			posts.GET("", middleware.OptionalAuthMiddleware(), routes.GetPosts)
			posts.GET("/:id", middleware.OptionalAuthMiddleware(), routes.GetPost)

			// 创建、更新、删除文章（需要认证）
			authPosts := posts.Group("/")
			authPosts.Use(middleware.AuthMiddleware())
			{
				authPosts.POST("", routes.CreatePost)
				authPosts.PUT("/:id", routes.UpdatePost)
				authPosts.DELETE("/:id", routes.DeletePost)
			}
		}

		// 评论相关路由
		comments := api.Group("/posts/:id/comments")
		{
			// 获取评论列表（无需认证）
			comments.GET("", routes.GetComments)

			// 创建评论（需要认证）
			comments.POST("", middleware.AuthMiddleware(), routes.CreateComment)
		}

		// 删除评论（需要认证）
		api.DELETE("/comments/:commentId", middleware.AuthMiddleware(), routes.DeleteComment)
	}

	// 启动服务器
	port := ":8000"
	utils.Info("Server is running on port %s", port)
	utils.Info("API documentation available at http://localhost%s/api/v1/health", port)

	if err := r.Run(port); err != nil {
		utils.Error("Failed to start server: %v", err)
	}
}