package controller

import (
	"blog-backend/models"
	"blog-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建服务实例
var postService = services.NewPostService()

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"error":   "Unauthorized",
		})
		return
	}

	var req models.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   "Invalid request data",
		})
		return
	}

	// 调用服务层创建文章
	post, err := postService.CreatePost(req.Title, req.Content, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create post",
			"error":   "Failed to create post",
		})
		return
	}

	// 直接返回文章对象，匹配测试期望的格式
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

// GetPosts 获取文章列表
func GetPosts(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 调用服务层获取文章列表
	posts, total, err := postService.GetPosts(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch posts",
			"error":   "Failed to fetch posts",
		})
		return
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	
	// 构建分页信息，直接返回符合测试期望的格式
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetPost 获取单个文章详情
func GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid post ID",
			"error":   "Invalid post ID",
		})
		return
	}

	// 调用服务层获取文章详情
	post, err := postService.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Post not found",
			"error":   "Post not found",
		})
		return
	}

	// 直接返回文章对象
	c.JSON(http.StatusOK, post)
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"error":   "Unauthorized",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid post ID",
			"error":   "Invalid post ID",
		})
		return
	}

	// 调用服务层更新文章
	var req models.PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   "Invalid request data",
		})
		return
	}

	post, err := postService.UpdatePost(uint(id), req.Title, req.Content, userID.(uint))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Post not found",
				"error":   "Post not found",
			})
		} else if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You don't have permission to update this post",
				"error":   "You don't have permission to update this post",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to update post",
				"error":   "Failed to update post",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    post,
	})
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
			"error":   "Unauthorized",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid post ID",
			"error":   "Invalid post ID",
		})
		return
	}

	// 调用服务层删除文章
	err = postService.DeletePost(uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Post not found",
				"error":   "Post not found",
			})
		} else if err.Error() == "permission denied" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "You don't have permission to delete this post",
				"error":   "You don't have permission to delete this post",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to delete post",
				"error":   "Failed to delete post",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}