package models

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // 密码不返回给客户端
	Email    string `gorm:"unique;not null" json:"email"`
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title   string    `gorm:"not null" json:"title"`
	Content string    `gorm:"not null" json:"content"`
	UserID  uint      `json:"user_id"`
	User    User      `json:"user,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `json:"user,omitempty"`
	PostID  uint   `json:"post_id"`
	Post    Post   `json:"post,omitempty"`
}

// 用户注册请求结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// 用户登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 文章创建/更新请求结构体
type PostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// 评论创建请求结构体
type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}