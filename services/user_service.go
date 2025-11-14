package services

import (
	"blog-backend/config"
	"blog-backend/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService 定义用户相关的业务逻辑接口
type UserService interface {
	// Register 用户注册
	Register(req *models.RegisterRequest) error
	// Login 用户登录，返回用户信息和JWT令牌
	Login(req *models.LoginRequest) (*models.User, string, error)
	// GetUserByID 根据ID获取用户信息
	GetUserByID(id uint) (*models.User, error)
	// GetUserByUsername 根据用户名获取用户信息
	GetUserByUsername(username string) (*models.User, error)
}

// userService 是UserService接口的实现
type userService struct{}

// NewUserService 创建一个新的UserService实例
func NewUserService() UserService {
	return &userService{}
}

// Register 用户注册实现
func (s *userService) Register(req *models.RegisterRequest) error {
	// 检查用户名是否已存在
	var existingUser models.User
	db := config.GetDB()
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := db.Create(&user).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

// Login 用户登录实现
func (s *userService) Login(req *models.LoginRequest) (*models.User, string, error) {
	// 查找用户
	var user models.User
	db := config.GetDB()
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// 生成JWT
	jwtConfig := config.GetJWTConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(jwtConfig.ExpiresIn).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtConfig.SecretKey))
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return &user, tokenString, nil
}

// GetUserByID 根据ID获取用户信息实现
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	db := config.GetDB()
	if err := db.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户信息实现
func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	db := config.GetDB()
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}