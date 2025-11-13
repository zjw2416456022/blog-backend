package config

import (
	"blog-backend/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接并进行自动迁移
func InitDB() {
	var err error
	
	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	log.Println("Database connected successfully")
	
	// 自动迁移数据库表结构
	err = DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)
	
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	log.Println("Database migration completed")
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}