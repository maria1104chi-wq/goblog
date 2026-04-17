package middleware

import (
	"blog-backend/internal/config"
	"blog-backend/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 提取Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证令牌",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "未认证",
			})
			c.Abort()
			return
		}

		if role.(int8) < 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path

		println(f"[日志] {method} {path} {statusCode} {latency} {clientIP}")
	}
}

// RateLimitMiddleware 简单的限流中间件（简化实现）
func RateLimitMiddleware(limit int, duration time.Duration) gin.HandlerFunc {
	// 实际项目中应该使用redis等实现分布式限流
	requestCount := make(map[string]int)
	lastReset := time.Now()

	return func(c *gin.Context) {
		// 定期重置计数器
		if time.Since(lastReset) > duration {
			requestCount = make(map[string]int)
			lastReset = time.Now()
		}

		clientIP := c.ClientIP()
		count := requestCount[clientIP]

		if count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		requestCount[clientIP] = count + 1
		c.Next()
	}
}

// SecurityMiddleware 安全中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 添加安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// GetClientIP 获取客户端真实IP
func GetClientIP(c *gin.Context) string {
	// 尝试从X-Forwarded-For获取
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For可能有多个IP，取第一个
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 尝试从X-Real-IP获取
	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 直接获取
	return c.ClientIP()
}
