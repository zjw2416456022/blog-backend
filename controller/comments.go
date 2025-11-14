package controller

import (
	"blog-backend/models"
	"blog-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建服务实例
var commentService = services.NewCommentService()

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"error":   "Unauthorized",
		})
		return
	}

	// 获取文章ID
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid post ID",
			"error":   "Invalid post ID",
		})
		return
	}

	// 绑定评论内容
	var req models.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   "Invalid request data",
		})
		return
	}

	// 调用服务层创建评论
	comment, err := commentService.CreateComment(req.Content, userID.(uint), uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Post not found",
				"error":   "Post not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   err.Error(),
			})
		}
		return
	}

	// 直接返回评论对象，匹配测试期望的格式
	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// GetComments 获取文章评论列表
func GetComments(c *gin.Context) {
	// 获取文章ID
	postIDStr := c.Param("postId")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid post ID",
			"error":   "Invalid post ID",
		})
		return
	}

	// 调用服务层获取评论列表
	comments, _, err := commentService.GetComments(uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Post not found",
				"error":   "Post not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, comments)
}

// UpdateComment 更新评论（只有评论作者可以更新）
func UpdateComment(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"error":   "Unauthorized",
		})
		return
	}

	// 获取评论ID
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid comment ID",
			"error":   "Invalid comment ID",
		})
		return
	}

	// 绑定评论内容
	var req models.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   "Invalid request data",
		})
		return
	}

	// 调用服务层更新评论
	comment, err := commentService.UpdateComment(uint(commentID), req.Content, userID.(uint))
	if err != nil {
		if err.Error() == "comment not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Comment not found",
				"error":   "Comment not found",
			})
		} else if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You don't have permission to update this comment",
				"error":   "You don't have permission to update this comment",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": comment,
	})
}

// DeleteComment 删除评论（评论作者或文章作者可以删除）
func DeleteComment(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"error":   "Unauthorized",
		})
		return
	}

	// 获取评论ID
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid comment ID",
			"error":   "Invalid comment ID",
		})
		return
	}

	// 调用服务层删除评论
	err = commentService.DeleteComment(uint(commentID), userID.(uint))
	if err != nil {
		if err.Error() == "comment not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Comment not found",
				"error":   "Comment not found",
			})
		} else if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You don't have permission to delete this comment",
				"error":   "You don't have permission to delete this comment",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}