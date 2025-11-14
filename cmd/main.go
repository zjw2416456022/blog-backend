package main

import (
	"blog-backend/api"
	"blog-backend/config"
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

	// 使用api包中的路由配置
	api.SetupRoutes(r)

	// 启动服务器
	port := ":8000"
	utils.Info("Server is running on port %s", port)
	utils.Info("API documentation available at http://localhost%s/api/v1/health", port)

	if err := r.Run(port); err != nil {
		utils.Error("Failed to start server: %v", err)
	}
}
