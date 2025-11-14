package tests

import (
	"blog-backend/api"
	"blog-backend/config"
	"blog-backend/models"
	"blog-backend/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	r           *gin.Engine
	testDB      *gorm.DB
	testUserID  uint
	testPostID  uint
	testCommentID uint
	testToken   string
)

// setupTest 设置测试环境
func setupTest(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 初始化测试数据库
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 保存原始DB实例并替换为测试DB
	originalDB := config.DB
	config.DB = testDB

	// 自动迁移表结构
	err = testDB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	assert.NoError(t, err)

	// 创建Gin引擎
	r = gin.Default()

	// 使用api包中的路由配置
	api.SetupRoutes(r)

	// 清理函数：恢复原始DB
	t.Cleanup(func() {
		config.DB = originalDB
	})
}

// TestRegister 测试用户注册功能
func TestRegister(t *testing.T) {
	setupTest(t)

	// 准备测试数据
	registerData := models.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}
	data, _ := json.Marshal(registerData)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	// 记录响应
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// 验证用户是否创建成功
	var user models.User
	result := testDB.Where("username = ?", "testuser").First(&user)
	assert.NoError(t, result.Error)
	assert.Equal(t, "test@example.com", user.Email)
}

// TestLogin 测试用户登录功能
func TestLogin(t *testing.T) {
	setupTest(t)

	// 先注册一个用户
	registerData := models.RegisterRequest{
		Username: "testuser",
		Password: "password123",
		Email:    "test@example.com",
	}
	data, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 测试登录
	loginData := models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	data, _ = json.Marshal(loginData)
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 解析响应获取token
	var response struct {
		Message string `json:"message"`
		Token   string `json:"token"`
		User    struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "Login successful", response.Message)
	
	// 保存测试用户ID和token供后续测试使用
	testUserID = response.User.ID
	testToken = response.Token
}

// TestCreatePost 测试创建文章功能
func TestCreatePost(t *testing.T) {
	setupTest(t)
	TestLogin(t) // 先登录获取token

	// 准备测试数据
	postData := models.PostRequest{
		Title:   "Test Post",
		Content: "This is a test post content.",
	}
	data, _ := json.Marshal(postData)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/posts/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)

	// 记录响应
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusCreated, w.Code)
	
	// 解析响应获取文章ID
	var response struct {
		Message string      `json:"message"`
		Post    models.Post `json:"post"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Post created successfully", response.Message)
	assert.NotZero(t, response.Post.ID)
	
	// 保存测试文章ID供后续测试使用
	testPostID = response.Post.ID
}

// TestGetPosts 测试获取文章列表功能
func TestGetPosts(t *testing.T) {
	setupTest(t)
	TestCreatePost(t) // 先创建一篇文章

	// 创建请求
	req, _ := http.NewRequest("GET", "/api/v1/posts", nil)
	
	// 记录响应
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 解析响应
	var response struct {
		Posts      []models.Post `json:"posts"`
		Pagination struct {
			Page       int `json:"page"`
			PageSize   int `json:"page_size"`
			Total      int `json:"total"`
			TotalPages int `json:"total_pages"`
		} `json:"pagination"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// 验证文章列表不为空
	assert.NotEmpty(t, response.Posts)
	assert.Equal(t, 1, response.Pagination.Total)
}

// TestCreateComment 测试创建评论功能
func TestCreateComment(t *testing.T) {
	setupTest(t)
	TestCreatePost(t) // 先创建一篇文章

	// 准备测试数据
	commentData := models.CommentRequest{
		Content: "This is a test comment.",
	}
	data, _ := json.Marshal(commentData)

	// 创建请求
	req, _ := http.NewRequest("POST", "/api/v1/posts/"+strconv.Itoa(int(testPostID))+"/comments", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)

	// 记录响应
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusCreated, w.Code)
}

// TestUpdatePost 测试更新文章功能
func TestUpdatePost(t *testing.T) {
	setupTest(t)
	TestCreatePost(t) // 先创建一篇文章

	// 准备测试数据
	updatedData := models.PostRequest{
		Title:   "Updated Test Post",
		Content: "This is updated test post content.",
	}
	data, _ := json.Marshal(updatedData)

	// 创建请求
	req, _ := http.NewRequest("PUT", "/api/v1/posts/"+strconv.Itoa(int(testPostID)), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)

	// 记录响应
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestDeletePost 测试删除文章功能
func TestDeletePost(t *testing.T) {
	setupTest(t)
	TestCreatePost(t) // 先创建一篇文章

	// 创建请求
	req, _ := http.NewRequest("DELETE", "/api/v1/posts/"+strconv.Itoa(int(testPostID)), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	// 记录响应
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUpdateComment 测试更新评论功能
func TestUpdateComment(t *testing.T) {
	setupTest(t)
	TestCreatePost(t) // 先创建一篇文章
	
	// 创建评论
	commentData := models.CommentRequest{
		Content: "This is a test comment.",
	}
	data, _ := json.Marshal(commentData)
	req, _ := http.NewRequest("POST", "/api/v1/posts/"+strconv.Itoa(int(testPostID))+"/comments", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	// 解析响应获取评论ID
	var response struct {
		Message string         `json:"message"`
		Comment models.Comment `json:"comment"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	testCommentID = response.Comment.ID
	
	// 准备更新数据
	updatedCommentData := models.CommentRequest{
		Content: "This is an updated test comment.",
	}
	data, _ = json.Marshal(updatedCommentData)
	
	// 创建更新请求
	updateReq, _ := http.NewRequest("PUT", "/api/v1/comments/"+strconv.Itoa(int(testCommentID)), bytes.NewBuffer(data))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+testToken)
	
	// 记录更新响应
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)
	
	// 验证响应
	assert.Equal(t, http.StatusOK, updateW.Code)
	
	// 验证数据库中的评论是否已更新
	var comment models.Comment
	result := testDB.First(&comment, testCommentID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "This is an updated test comment.", comment.Content)
}

// TestDeleteComment 测试删除评论功能
func TestDeleteComment(t *testing.T) {
	setupTest(t)
	TestLogin(t) // 确保先登录获取token
	
	// 直接在本函数中创建文章，确保testPostID被正确设置
	postData := models.PostRequest{
		Title:   "Test Post for Comment",
		Content: "This is a test post for comment deletion.",
	}
	data, _ := json.Marshal(postData)
	req, _ := http.NewRequest("POST", "/api/v1/posts/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	// 解析响应获取文章ID
	var postResponse struct {
		Message string      `json:"message"`
		Post    models.Post `json:"post"`
	}
	json.Unmarshal(w.Body.Bytes(), &postResponse)
	testPostID = postResponse.Post.ID
	
	// 创建评论
	commentData := models.CommentRequest{
		Content: "This is a test comment for deletion.",
	}
	data, _ = json.Marshal(commentData)
	req, _ = http.NewRequest("POST", "/api/v1/posts/"+strconv.Itoa(int(testPostID))+"/comments", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	// 解析响应获取评论ID
	var response struct {
		Message string         `json:"message"`
		Comment models.Comment `json:"comment"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	testCommentID = response.Comment.ID
	
	// 创建删除请求
	deleteReq, _ := http.NewRequest("DELETE", "/api/v1/comments/"+strconv.Itoa(int(testCommentID)), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+testToken)
	
	// 记录删除响应
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)
	
	// 验证响应
	assert.Equal(t, http.StatusOK, deleteW.Code)
	
	// 验证数据库中的评论是否已删除
	var deletedComment models.Comment
	result := testDB.First(&deletedComment, testCommentID)
	assert.Error(t, result.Error) // 应该返回错误，因为评论已被删除
}

// 主测试函数
func TestMain(m *testing.M) {
	// 初始化日志
	utils.InitLogger()
	
	// 运行测试
	exitVal := m.Run()
	
	// 清理资源
	os.Exit(exitVal)
}