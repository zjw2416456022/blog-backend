package services

import (
	"blog-backend/config"
	"blog-backend/models"
	"errors"
)

// CommentService 评论服务接口
type CommentService interface {
	CreateComment(content string, userID uint, postID uint) (*models.Comment, error)
	GetComments(postID uint) ([]models.Comment, int, error)
	UpdateComment(commentID uint, content string, userID uint) (*models.Comment, error)
	DeleteComment(commentID uint, userID uint) error
}

// commentService 评论服务实现
type commentService struct{}

// NewCommentService 创建评论服务实例
func NewCommentService() CommentService {
	return &commentService{}
}

// CreateComment 创建评论
func (s *commentService) CreateComment(content string, userID uint, postID uint) (*models.Comment, error) {
	db := config.GetDB()
	
	// 检查文章是否存在
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		return nil, errors.New("post not found")
	}
	
	// 创建评论
	comment := models.Comment{
		Content: content,
		UserID:  userID,
		PostID:  postID,
	}
	
	if err := db.Create(&comment).Error; err != nil {
		return nil, err
	}
	
	// 重新查询以获取关联信息
	db.Preload("User").First(&comment, comment.ID)
	
	return &comment, nil
}

// GetComments 获取文章的所有评论
func (s *commentService) GetComments(postID uint) ([]models.Comment, int, error) {
	db := config.GetDB()
	
	// 检查文章是否存在
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		return nil, 0, errors.New("post not found")
	}
	
	// 获取评论列表（按创建时间倒序）
	var comments []models.Comment
	if err := db.Where("post_id = ?", postID).Preload("User").Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	
	return comments, len(comments), nil
}

// UpdateComment 更新评论（只有评论作者可以更新）
func (s *commentService) UpdateComment(commentID uint, content string, userID uint) (*models.Comment, error) {
	db := config.GetDB()
	
	// 查找评论
	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		return nil, errors.New("comment not found")
	}
	
	// 检查权限：只有评论作者可以更新
	if comment.UserID != userID {
		return nil, errors.New("permission denied")
	}
	
	// 更新评论内容
	comment.Content = content
	if err := db.Save(&comment).Error; err != nil {
		return nil, err
	}
	
	// 重新查询以获取关联信息
	db.Preload("User").First(&comment, comment.ID)
	
	return &comment, nil
}

// DeleteComment 删除评论（评论作者或文章作者可以删除）
func (s *commentService) DeleteComment(commentID uint, userID uint) error {
	db := config.GetDB()
	
	// 查找评论
	var comment models.Comment
	if err := db.Preload("Post").First(&comment, commentID).Error; err != nil {
		return errors.New("comment not found")
	}
	
	// 检查权限：只有评论作者或文章作者可以删除
	if comment.UserID != userID && comment.Post.UserID != userID {
		return errors.New("permission denied")
	}
	
	// 删除评论
	if err := db.Delete(&comment).Error; err != nil {
		return err
	}
	
	return nil
}