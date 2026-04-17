package service

import (
	"blog-backend/internal/model"
	"blog-backend/internal/repository"
	"blog-backend/internal/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

// ArticleService 文章服务
type ArticleService struct {
}

// NewArticleService 创建文章服务实例
func NewArticleService() *ArticleService {
	return &ArticleService{}
}

// GetArticleList 获取文章列表（分页）
func (s *ArticleService) GetArticleList(page, pageSize int, status int8) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	db := repository.GetDB().Model(&model.Article{})

	if status >= 0 {
		db = db.Where("status = ?", status)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Preload("Category").Preload("Tags").Preload("Author").
		Order("is_top DESC, published_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetArticleBySlug 根据slug获取文章
func (s *ArticleService) GetArticleBySlug(slug string) (*model.Article, error) {
	var article model.Article
	if err := repository.GetDB().Preload("Category").Preload("Tags").Preload("Author").
		Preload("Attachments").Where("slug = ? AND status = ?", slug, 1).First(&article).Error; err != nil {
		return nil, err
	}

	// 增加浏览次数
	repository.GetDB().Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	return &article, nil
}

// GetHotArticles 获取热门文章（按浏览量）
func (s *ArticleService) GetHotArticles(limit int) ([]model.Article, error) {
	var articles []model.Article
	if err := repository.GetDB().Where("status = ?", 1).
		Order("view_count DESC").Limit(limit).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(article *model.Article, tagIDs []uint64) error {
	return repository.WithTransaction(func(tx *gorm.DB) error {
		// 生成slug
		if article.Slug == "" {
			article.Slug = utils.GenerateSlug(article.Title)
		}

		// 检查slug是否已存在
		var existing model.Article
		if tx.Where("slug = ?", article.Slug).First(&existing).Error == nil {
			// slug已存在，添加时间戳
			article.Slug = article.Slug + "-" + time.Now().Format("20060102150405")
		}

		// 自动提取摘要
		if article.Summary == "" {
			article.Summary = utils.ExtractSummary(article.Content, 200)
		}

		// Markdown转HTML
		htmlContent, err := utils.MarkdownToHTML(article.Content)
		if err != nil {
			htmlContent = article.Content
		}
		article.ContentHTML = utils.SanitizeHTML(htmlContent)

		// 敏感词过滤
		article.Title = utils.FilterText(article.Title)
		article.Summary = utils.FilterText(article.Summary)
		article.ContentHTML = utils.FilterText(article.ContentHTML)

		// 自动SEO优化
		if utils.GlobalFilter != nil { // 简化判断
			if article.SEOTitle == "" {
				article.SEOTitle = article.Title
			}
			if article.SEOKeywords == "" {
				article.SEOKeywords = utils.GenerateSEOKeywords(article.Title, article.Summary, nil)
			}
			if article.SEODescription == "" {
				article.SEODescription = utils.GenerateSEODescription(article.Title, article.Summary, 200)
			}
		}

		// 创建文章
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		// 关联标签
		if len(tagIDs) > 0 {
			for _, tagID := range tagIDs {
				tx.Create(&model.ArticleTag{
					ArticleID: article.ID,
					TagID:     tagID,
				})
			}
		}

		return nil
	})
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(id uint64, updates map[string]interface{}, tagIDs []uint64) error {
	return repository.WithTransaction(func(tx *gorm.DB) error {
		// 敏感词过滤
		if title, ok := updates["title"].(string); ok {
			updates["title"] = utils.FilterText(title)
		}
		if summary, ok := updates["summary"].(string); ok {
			updates["summary"] = utils.FilterText(summary)
		}
		if content, ok := updates["content"].(string); ok {
			htmlContent, _ := utils.MarkdownToHTML(content)
			updates["content_html"] = utils.FilterText(utils.SanitizeHTML(htmlContent))
		}

		// 更新文章
		if err := tx.Model(&model.Article{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		// 更新标签关联
		if tagIDs != nil {
			tx.Where("article_id = ?", id).Delete(&model.ArticleTag{})
			for _, tagID := range tagIDs {
				tx.Create(&model.ArticleTag{
					ArticleID: id,
					TagID:     tagID,
				})
			}
		}

		return nil
	})
}

// DeleteArticle 删除文章（软删除）
func (s *ArticleService) DeleteArticle(id uint64) error {
	now := time.Now()
	return repository.GetDB().Model(&model.Article{}).Where("id = ?", id).
		Update("deleted_at", &now).Error
}

// LikeArticle 点赞文章
func (s *ArticleService) LikeArticle(articleID uint64, ip string) error {
	return repository.WithTransaction(func(tx *gorm.DB) error {
		// 检查是否已点赞
		var like model.ArticleLike
		if tx.Where("article_id = ? AND ip_address = ?", articleID, ip).First(&like).Error == nil {
			return errors.New("已点赞过")
		}

		// 创建点赞记录
		if err := tx.Create(&model.ArticleLike{
			ArticleID: articleID,
			IPAddress: ip,
		}).Error; err != nil {
			return err
		}

		// 增加点赞数
		tx.Model(&model.Article{}).Where("id = ?", articleID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))

		return nil
	})
}

// ShareArticle 分享文章
func (s *ArticleService) ShareArticle(articleID uint64, platform, ip string) error {
	// 创建分享记录
	record := model.ArticleShare{
		ArticleID: articleID,
		Platform:  platform,
		IPAddress: ip,
	}
	if err := repository.GetDB().Create(&record).Error; err != nil {
		return err
	}

	// 增加分享数
	return repository.GetDB().Model(&model.Article{}).Where("id = ?", articleID).
		UpdateColumn("share_count", gorm.Expr("share_count + ?", 1)).Error
}

// SearchArticles 搜索文章
func (s *ArticleService) SearchArticles(keyword string, page, pageSize int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	searchPattern := "%" + keyword + "%"

	db := repository.GetDB().Model(&model.Article{}).
		Where("status = ? AND (title LIKE ? OR summary LIKE ? OR content LIKE ?)", 
			1, searchPattern, searchPattern, searchPattern)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Category").Preload("Tags").
		Order("published_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
