package main

import (
	"blog-backend/internal/config"
	"blog-backend/internal/handler"
	"blog-backend/internal/middleware"
	"blog-backend/internal/repository"
	"blog-backend/internal/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	if err := repository.InitDB(&cfg.Database); err != nil {
		log.Fatalf("数据库初始化失败：%v", err)
	}
	defer repository.CloseDB()

	// 自动迁移数据库表
	if err := repository.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败：%v", err)
	}

	// 初始化敏感词过滤器
	if err := utils.InitGlobalFilter(); err != nil {
		fmt.Printf("警告：敏感词过滤器初始化失败：%v\n", err)
	}

	// 初始化短信服务
	if err := utils.InitSMSService(&cfg.SMS); err != nil {
		fmt.Printf("警告：短信服务初始化失败：%v\n", err)
	}

	// 从文件导入敏感词（如果数据库中还没有）
	sensitiveFile := "static/Sensitive.txt"
	if _, err := os.Stat(sensitiveFile); err == nil {
		service := utils.NewSensitiveWordService()
		count, err := service.ImportSensitiveWordsFromFile(sensitiveFile)
		if err != nil {
			fmt.Printf("导入敏感词失败：%v\n", err)
		} else if count > 0 {
			fmt.Printf("从文件导入了 %d 个敏感词\n", count)
			// 重新加载敏感词到过滤器
			utils.GlobalFilter.LoadFromDatabase()
		}
	}

	// 创建Gin引擎
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.SecurityMiddleware())

	// 静态文件服务
	uploadPath := cfg.Upload.UploadPath
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Fatalf("创建上传目录失败：%v", err)
	}
	r.Static("/uploads", uploadPath)
	
	// 前端静态文件（生产环境）
	frontendDist := "../frontend/dist"
	if _, err := os.Stat(frontendDist); err == nil {
		r.Static("/assets", filepath.Join(frontendDist, "assets"))
		r.StaticFile("/", filepath.Join(frontendDist, "index.html"))
		r.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(frontendDist, "index.html"))
		})
	}

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关
		authHandler := handler.NewAuthHandler(&cfg.App)
		auth := api.Group("/auth")
		{
			auth.POST("/sms", authHandler.SendSMSCode)
			auth.POST("/login", authHandler.Login)
			auth.GET("/user", middleware.AuthMiddleware(&cfg.App), authHandler.GetUserInfo)
			auth.PUT("/profile", middleware.AuthMiddleware(&cfg.App), authHandler.UpdateProfile)
		}

		// 文章相关
		articleHandler := handler.NewArticleHandler(&cfg.App)
		articles := api.Group("/articles")
		{
			articles.GET("", articleHandler.GetArticleList)
			articles.GET("/hot", articleHandler.GetHotArticles)
			articles.GET("/search", articleHandler.SearchArticles)
			articles.GET("/:slug", articleHandler.GetArticleBySlug)
			
			// 需要认证的操作
			authArticles := articles.Group("")
			authArticles.Use(middleware.AuthMiddleware(&cfg.App))
			{
				authArticles.POST("", articleHandler.CreateArticle)
				authArticles.PUT("/:id", articleHandler.UpdateArticle)
				authArticles.DELETE("/:id", articleHandler.DeleteArticle)
				authArticles.POST("/:id/like", articleHandler.LikeArticle)
				authArticles.POST("/:id/share", articleHandler.ShareArticle)
			}
		}

		// 评论相关（待实现）
		// comments := api.Group("/comments")
		
		// 分类相关（待实现）
		// categories := api.Group("/categories")
		
		// 标签相关（待实现）
		// tags := api.Group("/tags")
		
		// 敏感词管理（待实现）
		// sensitive := api.Group("/sensitive")
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   gin.TimeFormat(gin.TimeNow()),
		})
	})

	// 启动服务器
	addr := ":" + cfg.Server.Port
	fmt.Printf("服务器启动在 %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败：%v", err)
	}
}
