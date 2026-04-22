# 个人博客系统

基于 Go + Gin + Vue3 的全栈个人博客系统，支持管理员手机号验证码登录、文章管理、评论管理、敏感词过滤等功能。

## 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **认证**: JWT
- **短信**: 阿里云短信服务

### 前端（待实现）
- **框架**: Vue 3
- **构建工具**: Vite
- **UI库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router

## 项目结构

```
blog-project/
├── backend/                 # Go后端
│   ├── cmd/                # 应用入口
│   │   └── main.go         # 主程序
│   ├── internal/           # 内部包
│   │   ├── config/         # 配置
│   │   ├── handler/        # HTTP处理器
│   │   ├── middleware/     # 中间件
│   │   ├── model/          # 数据模型
│   │   ├── repository/     # 数据访问层
│   │   ├── service/        # 业务逻辑层
│   │   └── utils/          # 工具函数
│   ├── static/             # 静态资源
│   │   ├── uploads/        # 上传文件
│   │   └── Sensitive.txt   # 敏感词库
│   ├── go.mod              # Go模块定义
│   └── Dockerfile          # Docker构建文件
├── frontend/               # Vue前端（待实现）
│   ├── src/
│   │   ├── assets/         # 静态资源
│   │   ├── components/     # 组件
│   │   ├── views/          # 页面
│   │   ├── router/         # 路由
│   │   ├── store/          # 状态管理
│   │   └── api/            # API调用
│   ├── public/             # 公共资源
│   ├── Dockerfile          # Docker构建文件
│   └── nginx.conf          # Nginx配置
├── db/                     # 数据库
│   └── schema.sql          # 数据库结构
├── docker-compose.yml      # Docker编排
└── Caddyfile              # Caddy配置
```

## 快速开始

### 环境要求
- Docker & Docker Compose
- 或本地安装：Go 1.21+, MySQL 8.0, Node.js 18+

### 使用Docker Compose启动

```bash
# 克隆项目
git clone <repository-url>
cd blog-project

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 本地开发

#### 1. 启动MySQL

```bash
# 使用Docker启动MySQL
docker run -d \
  --name blog_mysql \
  -e MYSQL_ROOT_PASSWORD=root123 \
  -e MYSQL_DATABASE=blog_db \
  -p 3306:3306 \
  mysql:8.0 \
  --character-set-server=utf8mb4 \
  --collation-server=utf8mb4_unicode_ci

# 导入数据库结构
mysql -h 127.0.0.1 -P 3306 -u root -proot123 blog_db < db/schema.sql
```

#### 2. 启动后端

```bash
cd backend

# 下载依赖
go mod download

# 运行
go run cmd/main.go
```

#### 3. 配置环境变量（可选）

创建 `.env` 文件：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root123
DB_NAME=blog_db

# 阿里云短信配置
ALIYUN_SMS_ACCESS_KEY=your_access_key
ALIYUN_SMS_ACCESS_SECRET=your_access_secret
ALIYUN_SMS_SIGN_NAME=你的签名
ALIYUN_SMS_TEMPLATE_CODE=SMS_123456789

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOUR=24

# 服务器配置
SERVER_PORT=8080
GIN_MODE=debug
```

## API接口

### 认证相关

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/auth/sms | 发送登录验证码 | 否 |
| POST | /api/auth/login | 管理员登录 | 否 |
| GET | /api/auth/user | 获取当前用户信息 | 是 |
| PUT | /api/auth/profile | 更新用户资料 | 是 |

### 文章相关

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/articles | 获取文章列表 | 否 |
| GET | /api/articles/:slug | 获取文章详情 | 否 |
| GET | /api/articles/hot | 获取热门文章 | 否 |
| GET | /api/articles/search | 搜索文章 | 否 |
| POST | /api/articles | 创建文章 | 是 |
| PUT | /api/articles/:id | 更新文章 | 是 |
| DELETE | /api/articles/:id | 删除文章 | 是 |
| POST | /api/articles/:id/like | 点赞文章 | 否 |
| POST | /api/articles/:id/share | 分享文章 | 否 |

## 主要功能

### 已实现
- ✅ 管理员手机号验证码登录
- ✅ JWT令牌认证
- ✅ 文章CRUD操作
- ✅ Markdown转HTML
- ✅ 敏感词过滤（Trie树算法）
- ✅ 文章点赞、分享
- ✅ 评论功能
- ✅ 浏览次数统计
- ✅ SEO优化（自动生成关键词和描述）
- ✅ Docker容器化部署
- ✅ 数据库自动迁移

### 待实现
- ⏳ 前端Vue3界面
- ⏳ 文章分类管理
- ⏳ 文章标签管理
- ⏳ 评论审核
- ⏳ 文件上传（图片/PDF/视频）
- ⏳ 数据统计面板
- ⏳ 定时备份
- ⏳ Redis缓存

## 数据库表

- `users` - 管理员用户
- `sms_codes` - 短信验证码
- `articles` - 文章
- `categories` - 分类
- `tags` - 标签
- `article_tags` - 文章标签关联
- `article_attachments` - 文章附件
- `comments` - 评论
- `comment_likes` - 评论点赞
- `article_likes` - 文章点赞
- `article_shares` - 文章分享
- `sensitive_words` - 敏感词库
- `system_configs` - 系统配置
- `operation_logs` - 操作日志
- `daily_statistics` - 每日统计
- `backup_records` - 备份记录

## 安全特性

- JWT令牌认证
- 密码SHA256哈希
- 敏感词自动过滤
- XSS防护
- CORS跨域控制
- SQL注入防护（GORM参数化查询）
- 请求限流
- 安全响应头

## 默认账户

- 手机号：13800138000
- 初始密码：admin123（首次登录后请修改）

## 许可证

MIT License

## 联系方式

如有问题请提交Issue或联系开发者。
