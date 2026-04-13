package handler

import (
	"blog-backend/internal/config"
	"blog-backend/internal/middleware"
	"blog-backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	articleService *service.ArticleService
	config         *config.AppConfig
}

// NewArticleHandler 创建文章处理器
func NewArticleHandler(cfg *config.AppConfig) *ArticleHandler {
	return &ArticleHandler{
		articleService: service.NewArticleService(),
		config:         cfg,
	}
}

// GetArticleList 获取文章列表
// @Summary 获取文章列表
// @Tags 文章
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query int false "状态" default(1)
// @Success 200 {object} map[string]interface{}
// @Router /api/articles [get]
func (h *ArticleHandler) GetArticleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(h.config.PageSize)))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "1"))

	articles, total, err := h.articleService.GetArticleList(page, pageSize, int8(status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文章列表失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      articles,
			"total":     total,
			"page":      page,
			"pageSize":  pageSize,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetArticleBySlug 根据slug获取文章
// @Summary 获取文章详情
// @Tags 文章
// @Produce json
// @Param slug path string true "文章slug"
// @Success 200 {object} map[string]interface{}
// @Router /api/articles/:slug [get]
func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	slug := c.Param("slug")

	article, err := h.articleService.GetArticleBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文章不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    article,
	})
}

// GetHotArticles 获取热门文章
// @Summary 获取热门文章
// @Tags 文章
// @Produce json
// @Param limit query int false "数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/articles/hot [get]
func (h *ArticleHandler) GetHotArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(h.config.HotArticleLimit)))

	articles, err := h.articleService.GetHotArticles(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取热门文章失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    articles,
	})
}

// CreateArticle 创建文章
// @Summary 创建文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param article body object true "文章信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req struct {
		Title      string   `json:"title" binding:"required"`
		Summary    string   `json:"summary"`
		Content    string   `json:"content" binding:"required"`
		CoverImage string   `json:"cover_image"`
		CategoryID uint64   `json:"category_id"`
		TagIDs     []uint64 `json:"tag_ids"`
		Status     int8     `json:"status"`
		IsTop      int8     `json:"is_top"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")
	article := &model.Article{
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		CategoryID: req.CategoryID,
		AuthorID:   userID.(uint64),
		Status:     req.Status,
		IsTop:      req.IsTop,
	}

	if err := h.articleService.CreateArticle(article, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建文章失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    article,
	})
}

// UpdateArticle 更新文章
// @Summary 更新文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param article body object true "文章信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/articles/:id [put]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Title      string   `json:"title"`
		Summary    string   `json:"summary"`
		Content    string   `json:"content"`
		CoverImage string   `json:"cover_image"`
		CategoryID uint64   `json:"category_id"`
		TagIDs     []uint64 `json:"tag_ids"`
		Status     int8     `json:"status"`
		IsTop      int8     `json:"is_top"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	updates["cover_image"] = req.CoverImage
	updates["category_id"] = req.CategoryID
	updates["status"] = req.Status
	updates["is_top"] = req.IsTop

	if err := h.articleService.UpdateArticle(id, updates, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新文章失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Tags 文章
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/articles/:id [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.articleService.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除文章失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// LikeArticle 点赞文章
// @Summary 点赞文章
// @Tags 文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/articles/:id/like [post]
func (h *ArticleHandler) LikeArticle(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ip := middleware.GetClientIP(c)

	if err := h.articleService.LikeArticle(id, ip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "点赞成功",
	})
}

// ShareArticle 分享文章
// @Summary 分享文章
// @Tags 文章
// @Produce json
// @Param id path int true "文章ID"
// @Param platform query string false "分享平台" default(link)
// @Success 200 {object} map[string]interface{}
// @Router /api/articles/:id/share [post]
func (h *ArticleHandler) ShareArticle(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	platform := c.DefaultQuery("platform", "link")
	ip := middleware.GetClientIP(c)

	if err := h.articleService.ShareArticle(id, platform, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "分享失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "分享成功",
	})
}

// SearchArticles 搜索文章
// @Summary 搜索文章
// @Tags 文章
// @Produce json
// @Param keyword query string true "关键词"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/articles/search [get]
func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请输入搜索关键词",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(h.config.PageSize)))

	articles, total, err := h.articleService.SearchArticles(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      articles,
			"total":     total,
			"page":      page,
			"pageSize":  pageSize,
			"totalPage": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}
