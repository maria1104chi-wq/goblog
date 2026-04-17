package config

import (
	"os"
	"strconv"
	"time"
)

// Config 应用配置结构体
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	SMS      SMSConfig
	Upload   UploadConfig
	App      AppConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string
	Mode         string // debug, release, test
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// SMSConfig 短信配置
type SMSConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
	CodeExpireMin   int
	CodeLength      int
}

// UploadConfig 上传配置
type UploadConfig struct {
	MaxImageSize int64
	MaxPDFSize   int64
	MaxVideoSize int64
	AllowedTypes []string
	UploadPath   string
}

// AppConfig 应用配置
type AppConfig struct {
	JWTSecret       string
	JWTExpire       time.Duration
	EnableSEOAuto   bool
	CommentNeedAudit bool
	AllowAnonymous  bool
	HotArticleLimit int
	PageSize        int
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Mode:         getEnv("GIN_MODE", "release"),
			ReadTimeout:  time.Duration(getEnvInt("READ_TIMEOUT", 30)) * time.Second,
			WriteTimeout: time.Duration(getEnvInt("WRITE_TIMEOUT", 30)) * time.Second,
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "mysql"),
			Port:            getEnv("DB_PORT", "3306"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", "root123"),
			DBName:          getEnv("DB_NAME", "blog_db"),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE", 10),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN", 100),
			ConnMaxLifetime: time.Duration(getEnvInt("DB_MAX_LIFETIME", 3600)) * time.Second,
		},
		SMS: SMSConfig{
			AccessKeyID:     getEnv("ALIYUN_SMS_ACCESS_KEY", ""),
			AccessKeySecret: getEnv("ALIYUN_SMS_ACCESS_SECRET", ""),
			SignName:        getEnv("ALIYUN_SMS_SIGN_NAME", "我的博客"),
			TemplateCode:    getEnv("ALIYUN_SMS_TEMPLATE_CODE", "SMS_123456789"),
			CodeExpireMin:   getEnvInt("SMS_CODE_EXPIRE_MIN", 5),
			CodeLength:      getEnvInt("SMS_CODE_LENGTH", 6),
		},
		Upload: UploadConfig{
			MaxImageSize: getEnvInt64("MAX_IMAGE_SIZE", 5*1024*1024),
			MaxPDFSize:   getEnvInt64("MAX_PDF_SIZE", 10*1024*1024),
			MaxVideoSize: getEnvInt64("MAX_VIDEO_SIZE", 50*1024*1024),
			AllowedTypes: []string{"jpg", "jpeg", "png", "gif", "webp", "pdf", "mp4", "avi", "mov"},
			UploadPath:   getEnv("UPLOAD_PATH", "/app/backend/static/uploads"),
		},
		App: AppConfig{
			JWTSecret:        getEnv("JWT_SECRET", "blog-jwt-secret-key-2024"),
			JWTExpire:        time.Duration(getEnvInt("JWT_EXPIRE_HOUR", 24)) * time.Hour,
			EnableSEOAuto:    getEnvBool("ENABLE_SEO_AUTO", true),
			CommentNeedAudit: getEnvBool("COMMENT_NEED_AUDIT", false),
			AllowAnonymous:   getEnvBool("ALLOW_ANONYMOUS", true),
			HotArticleLimit:  getEnvInt("HOT_ARTICLE_LIMIT", 10),
			PageSize:         getEnvInt("PAGE_SIZE", 10),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt 获取整数类型环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvInt64 获取int64类型环境变量
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool 获取布尔类型环境变量
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
