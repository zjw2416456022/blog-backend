package services

import (
	"blog-backend/config"
	"blog-backend/models"
	"errors"
)

// PostService 定义文章相关的业务逻辑接口
type PostService interface {
	// CreatePost 创建文章
	CreatePost(title, content string, userID uint) (*models.Post, error)
	// GetPosts 获取文章列表（支持分页）
	GetPosts(page, pageSize int) ([]models.Post, int64, error)
	// GetPostByID 根据ID获取文章详情
	GetPostByID(id uint) (*models.Post, error)
	// UpdatePost 更新文章
	UpdatePost(id uint, title, content string, userID uint) (*models.Post, error)
	// DeletePost 删除文章
	DeletePost(id, userID uint) error
}

// postService 是PostService接口的实现
type postService struct{}

// NewPostService 创建一个新的PostService实例
func NewPostService() PostService {
	return &postService{}
}

// CreatePost 创建文章实现
func (s *postService) CreatePost(title, content string, userID uint) (*models.Post, error) {
	// 创建文章
	post := models.Post{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	db := config.GetDB()
	if err := db.Create(&post).Error; err != nil {
		return nil, errors.New("failed to create post")
	}

	// 重新查询以获取关联的用户信息
	db.Preload("User").First(&post, post.ID)

	return &post, nil
}

// GetPosts 获取文章列表实现
func (s *postService) GetPosts(page, pageSize int) ([]models.Post, int64, error) {
	// 参数验证和调整
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询文章列表
	var posts []models.Post
	db := config.GetDB()
	total := int64(0)
	
	// 统计总数
	db.Model(&models.Post{}).Count(&total)
	
	// 查询带分页的文章，预加载用户信息
	if err := db.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, 0, errors.New("failed to fetch posts")
	}

	return posts, total, nil
}

// GetPostByID 根据ID获取文章详情实现
func (s *postService) GetPostByID(id uint) (*models.Post, error) {
	var post models.Post
	db := config.GetDB()
	// 预加载用户和评论信息
	if err := db.Preload("User").Preload("Comments").Preload("Comments.User").First(&post, id).Error; err != nil {
		return nil, errors.New("post not found")
	}

	return &post, nil
}

// UpdatePost 更新文章实现
func (s *postService) UpdatePost(id uint, title, content string, userID uint) (*models.Post, error) {
	db := config.GetDB()
	var post models.Post
	
	// 查找文章
	if err := db.First(&post, id).Error; err != nil {
		return nil, errors.New("post not found")
	}
	
	// 检查是否是文章作者
	if post.UserID != userID {
		return nil, errors.New("permission denied")
	}

	// 更新文章内容
	post.Title = title
	post.Content = content

	if err := db.Save(&post).Error; err != nil {
		return nil, errors.New("failed to update post")
	}

	// 重新查询以获取关联信息
	db.Preload("User").First(&post, id)

	return &post, nil
}

// DeletePost 删除文章实现
func (s *postService) DeletePost(id, userID uint) error {
	db := config.GetDB()
	var post models.Post
	
	// 查找文章
	if err := db.First(&post, id).Error; err != nil {
		return errors.New("post not found")
	}
	
	// 检查是否是文章作者
	if post.UserID != userID {
		return errors.New("permission denied")
	}

	// 删除文章
	if err := db.Delete(&post).Error; err != nil {
		return errors.New("failed to delete post")
	}

	return nil
}