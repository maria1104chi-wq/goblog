package service

import (
	"blog-backend/internal/model"
	"blog-backend/internal/repository"
	"blog-backend/internal/utils"
	"time"

	"gorm.io/gorm"
)

// CommentService 评论服务
type CommentService struct {
}

// NewCommentService 创建评论服务实例
func NewCommentService() *CommentService {
	return &CommentService{}
}

// GetCommentList 获取评论列表（按文章）
func (s *CommentService) GetCommentList(articleID uint64, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	db := repository.GetDB().Model(&model.Comment{}).
		Where("article_id = ? AND status = ? AND parent_id = ?", articleID, 1, 0)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 加载回复
	for i := range comments {
		var replies []model.Comment
		repository.GetDB().Where("parent_id = ? AND status = ?", comments[i].ID, 1).
			Order("created_at ASC").Find(&replies)
		comments[i].Replies = replies
	}

	return comments, total, nil
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(comment *model.Comment) error {
	// 敏感词过滤
	comment.Content = utils.FilterText(comment.Content)
	
	// 转换为HTML（简化处理，实际应该更复杂）
	comment.ContentHTML = comment.Content

	// 获取IP地理位置
	comment.IPLocation = utils.GetIPLocation(comment.IPAddress)

	return repository.WithTransaction(func(tx *gorm.DB) error {
		// 创建评论
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		// 增加文章评论数
		tx.Model(&model.Article{}).Where("id = ?", comment.ArticleID).
			UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))

		return nil
	})
}

// LikeComment 点赞评论
func (s *CommentService) LikeComment(commentID uint64, ip string) error {
	return repository.WithTransaction(func(tx *gorm.DB) error {
		// 检查是否已点赞
		var like model.CommentLike
		if tx.Where("comment_id = ? AND ip_address = ?", commentID, ip).First(&like).Error == nil {
			return nil // 已点赞过，不返回错误
		}

		// 创建点赞记录
		if err := tx.Create(&model.CommentLike{
			CommentID: commentID,
			IPAddress: ip,
		}).Error; err != nil {
			return err
		}

		// 增加点赞数
		tx.Model(&model.Comment{}).Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))

		return nil
	})
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(id uint64) error {
	now := time.Now()
	return repository.GetDB().Model(&model.Comment{}).Where("id = ?", id).
		Update("deleted_at", &now).Error
}

// AuditComment 审核评论
func (s *CommentService) AuditComment(id uint64, status int8) error {
	return repository.GetDB().Model(&model.Comment{}).Where("id = ?", id).
		Update("status", status).Error
}

// GetCommentCount 获取评论总数
func (s *CommentService) GetCommentCount(articleID uint64) (int64, error) {
	var count int64
	err := repository.GetDB().Model(&model.Comment{}).
		Where("article_id = ? AND status = ?", articleID, 1).
		Count(&count).Error
	return count, err
}
