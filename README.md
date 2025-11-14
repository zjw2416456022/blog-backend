# 博客系统后端 API

## 项目简介
这是一个基于Go语言和Gin框架开发的博客系统后端API，提供用户认证、文章管理和评论功能。

## 技术栈
- Go 1.24.0
- Gin Web框架
- GORM ORM框架
- SQLite数据库
- JWT认证

## 功能特性
- 用户注册和登录
- JWT认证保护API
- 文章CRUD操作
- 评论功能
- CORS跨域支持
- 健康检查接口

## 测试用例
项目包含完整的单元测试，测试文件位于`tests/api_test.go`。测试覆盖以下功能：

### 已实现的测试用例
1. **用户注册测试** (`TestRegister`)：验证用户注册功能是否正常工作
2. **用户登录测试** (`TestLogin`)：验证用户登录功能和JWT令牌生成
3. **文章创建测试** (`TestCreatePost`)：验证创建文章功能
4. **文章列表获取测试** (`TestGetPosts`)：验证获取文章列表功能
5. **评论创建测试** (`TestCreateComment`)：验证创建评论功能
6. **文章更新测试** (`TestUpdatePost`)：验证更新文章功能
7. **文章删除测试** (`TestDeletePost`)：验证删除文章功能

### 测试结果

#### 单元测试结果
所有测试用例均已成功通过：

```
=== RUN   TestRegister
--- PASS: TestRegister (0.13s)
=== RUN   TestLogin
--- PASS: TestLogin (0.12s)
=== RUN   TestCreatePost
--- PASS: TestCreatePost (0.13s)
=== RUN   TestGetPosts
--- PASS: TestGetPosts (0.13s)
=== RUN   TestCreateComment
--- PASS: TestCreateComment (0.13s)
=== RUN   TestUpdatePost
--- PASS: TestUpdatePost (0.14s)
=== RUN   TestDeletePost
--- PASS: TestDeletePost (0.13s)
PASS
ok      blog-backend/tests      1.426s
```

#### 手动API测试结果
通过手动测试，确认以下功能正常工作：

- ✅ 服务器启动成功，监听在8000端口
- ✅ 健康检查接口：`GET /health` 返回状态正常
- ✅ 用户注册：`POST /api/v1/auth/register` 成功创建用户
- ✅ 用户登录：`POST /api/v1/auth/login` 成功返回JWT令牌
- ✅ 创建文章：`POST /api/v1/posts/` 成功创建文章
- ✅ 获取文章列表：`GET /api/v1/posts` 成功返回文章列表

## API接口说明

### 用户相关接口
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/user/profile` - 获取用户信息（需要认证）

### 文章相关接口
- `GET /api/v1/posts` - 获取文章列表
- `GET /api/v1/posts/:id` - 获取文章详情
- `POST /api/v1/posts/` - 创建文章（需要认证）
- `PUT /api/v1/posts/:id` - 更新文章（需要认证）
- `DELETE /api/v1/posts/:id` - 删除文章（需要认证）

### 评论相关接口
- `GET /api/v1/posts/:id/comments` - 获取文章评论
- `POST /api/v1/posts/:id/comments` - 创建评论（需要认证）
- `DELETE /api/v1/comments/:commentId` - 删除评论（需要认证）

### 健康检查接口
- `GET /health` - 健康检查

## 运行测试
要运行测试，执行以下命令：
```bash
cd /Volumes/F/project/go/blog-backend
go test ./tests/... -v
```

这是一个使用 Go 语言、Gin 框架和 GORM 开发的个人博客系统后端服务。实现了用户认证、文章管理和评论功能的 RESTful API。

## 技术栈

- **编程语言**: Go 1.23+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: SQLite（支持切换到 MySQL/PostgreSQL）
- **认证**: JWT (JSON Web Token)
- **密码加密**: bcrypt
- **跨域支持**: CORS

## 功能特性

- **用户管理**
  - 用户注册
  - 用户登录
  - 获取用户个人信息

- **文章管理**
  - 创建文章（需认证）
  - 获取文章列表（支持分页）
  - 获取文章详情
  - 更新文章（仅作者可操作）
  - 删除文章（仅作者可操作）

- **评论功能**
  - 创建评论（需认证）
  - 获取文章评论列表
  - 删除评论（评论作者或文章作者可操作）

- **系统特性**
  - JWT 认证机制
  - 统一错误处理
  - 日志记录
  - CORS 支持
  - 健康检查接口

## 项目结构

```
blog-backend/
├── config/         # 配置文件
│   ├── db.go       # 数据库配置
│   └── jwt.go      # JWT 配置
├── middleware/     # 中间件
│   └── auth.go     # JWT 认证中间件
├── models/         # 数据模型
│   └── models.go   # 数据库模型定义
├── routes/         # 路由处理
│   ├── comments.go # 评论相关路由
│   ├── posts.go    # 文章相关路由
│   └── users.go    # 用户相关路由
├── utils/          # 工具函数
│   ├── logger.go   # 日志工具
│   └── response.go # 响应格式化工具
├── logs/           # 日志文件目录
├── go.mod          # Go 模块文件
├── go.sum          # 依赖版本锁定文件
├── main.go         # 应用入口
├── blog.db         # SQLite 数据库文件
└── README.md       # 项目说明文档
```

## 快速开始

### 环境要求

- Go 1.23 或更高版本
- Git

### 安装步骤

1. **克隆项目**

```bash
git clone https://your-repository-url/blog-backend.git
cd blog-backend
```

2. **安装依赖**

```bash
go mod tidy
```

3. **配置 JWT 密钥**

编辑 `config/jwt.go` 文件，修改 `SecretKey` 为安全的随机字符串：

```go
func GetJWTConfig() JWTConfig {
	return JWTConfig{
		SecretKey: "your_secure_random_secret_key_here", // 修改为安全的随机密钥
		ExpiresIn: 24 * time.Hour,
	}
}
```

4. **运行项目**

```bash
go run main.go
```

服务将在 `http://localhost:8000` 启动。

## API 接口文档

### 健康检查

```
GET /health
```

响应：
```json
{
  "status": "ok",
  "timestamp": 1621548000
}
```

### 用户认证

#### 注册

```
POST /api/v1/auth/register
```

请求体：
```json
{
  "username": "user123",
  "password": "password123",
  "email": "user@example.com"
}
```

响应：
```json
{
  "success": true,
  "message": "User registered successfully"
}
```

#### 登录

```
POST /api/v1/auth/login
```

请求体：
```json
{
  "username": "user123",
  "password": "password123"
}
```

响应：
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "user123",
      "email": "user@example.com"
    }
  }
}
```

### 用户信息

#### 获取个人信息（需认证）

```
GET /api/v1/user/profile
```

请求头：
```
Authorization: Bearer <token>
```

响应：
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "username": "user123",
      "email": "user@example.com",
      "created_at": "2023-01-01T00:00:00Z"
    }
  }
}
```

### 文章管理

#### 获取文章列表

```
GET /api/v1/posts?page=1&page_size=10
```

响应：
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "title": "我的第一篇博客",
        "content": "这是博客内容...",
        "user_id": 1,
        "user": {
          "id": 1,
          "username": "user123",
          "email": "user@example.com"
        },
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 50,
      "total_pages": 5
    }
  }
}
```

#### 获取文章详情

```
GET /api/v1/posts/:id
```

响应：
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "我的第一篇博客",
    "content": "这是博客内容...",
    "user_id": 1,
    "user": {
      "id": 1,
      "username": "user123",
      "email": "user@example.com"
    },
    "comments": [
      {
        "id": 1,
        "content": "很好的文章！",
        "user_id": 2,
        "user": {
          "id": 2,
          "username": "user456",
          "email": "user456@example.com"
        },
        "created_at": "2023-01-02T00:00:00Z"
      }
    ],
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

#### 创建文章（需认证）

```
POST /api/v1/posts
```

请求头：
```
Authorization: Bearer <token>
```

请求体：
```json
{
  "title": "新文章标题",
  "content": "文章内容..."
}
```

响应：
```json
{
  "success": true,
  "message": "Post created successfully",
  "data": {
    "id": 2,
    "title": "新文章标题",
    "content": "文章内容...",
    "user_id": 1,
    "user": {
      "id": 1,
      "username": "user123",
      "email": "user@example.com"
    },
    "created_at": "2023-01-03T00:00:00Z",
    "updated_at": "2023-01-03T00:00:00Z"
  }
}
```

#### 更新文章（需认证，仅作者可操作）

```
PUT /api/v1/posts/:id
```

请求头：
```
Authorization: Bearer <token>
```

请求体：
```json
{
  "title": "更新后的标题",
  "content": "更新后的内容..."
}
```

响应：
```json
{
  "success": true,
  "message": "Post updated successfully",
  "data": {
    "id": 1,
    "title": "更新后的标题",
    "content": "更新后的内容...",
    "user_id": 1,
    "user": {
      "id": 1,
      "username": "user123",
      "email": "user@example.com"
    },
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-04T00:00:00Z"
  }
}
```

#### 删除文章（需认证，仅作者可操作）

```
DELETE /api/v1/posts/:id
```

请求头：
```
Authorization: Bearer <token>
```

响应：
```json
{
  "success": true,
  "message": "Post deleted successfully"
}
```

### 评论管理

#### 获取文章评论

```
GET /api/v1/posts/:postId/comments
```

响应：
```json
{
  "success": true,
  "data": {
    "comments": [
      {
        "id": 1,
        "content": "很好的文章！",
        "user_id": 2,
        "user": {
          "id": 2,
          "username": "user456",
          "email": "user456@example.com"
        },
        "post_id": 1,
        "created_at": "2023-01-02T00:00:00Z"
      }
    ],
    "count": 1
  }
}
```

#### 创建评论（需认证）

```
POST /api/v1/posts/:postId/comments
```

请求头：
```
Authorization: Bearer <token>
```

请求体：
```json
{
  "content": "这是一条评论"
}
```

响应：
```json
{
  "success": true,
  "message": "Comment created successfully",
  "data": {
    "id": 2,
    "content": "这是一条评论",
    "user_id": 1,
    "user": {
      "id": 1,
      "username": "user123",
      "email": "user@example.com"
    },
    "post_id": 1,
    "created_at": "2023-01-05T00:00:00Z"
  }
}
```

#### 删除评论（需认证，评论作者或文章作者可操作）

```
DELETE /api/v1/comments/:commentId
```

请求头：
```
Authorization: Bearer <token>
```

响应：
```json
{
  "success": true,
  "message": "Comment deleted successfully"
}
```

## 测试

### 使用 Postman 测试

1. 下载并安装 [Postman](https://www.postman.com/)
2. 导入 API 集合（如果提供）
3. 按照上述 API 文档进行测试

### 运行服务器

启动服务器后，可以使用 curl 命令测试接口：

```bash
# 健康检查
curl http://localhost:8000/health

# 用户注册
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123","email":"test@example.com"}'

# 用户登录
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

## 部署

### 编译二进制文件

```bash
go build -o blog-backend main.go
```

### 运行二进制文件

```bash
./blog-backend
```

## 注意事项

1. 在生产环境中，请务必修改 JWT 密钥为安全的随机字符串
2. 考虑使用环境变量存储敏感配置信息
3. 对于高并发场景，建议切换到 MySQL 或 PostgreSQL 数据库
4. 建议添加请求频率限制和更完善的错误处理机制
5. 定期备份数据库文件

## 许可证

MIT License
