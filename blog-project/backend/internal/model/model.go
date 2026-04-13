package model

import (
	"time"
)

// User 用户模型（管理员）
type User struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Phone        string     `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	Email        string     `gorm:"size:100" json:"email"`
	Avatar       string     `gorm:"size:255" json:"avatar"`
	Role         int8       `gorm:"default:1" json:"role"` // 1-超级管理员，2-普通管理员
	Status       int8       `gorm:"default:1" json:"status"` // 0-禁用，1-启用
	LastLogin    *time.Time `json:"last_login"`
	LastLoginIP  string     `gorm:"size:45" json:"last_login_ip"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// SMSCode 短信验证码模型
type SMSCode struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone     string    `gorm:"size:20;not null;index:idx_phone_code" json:"phone"`
	Code      string    `gorm:"size:10;not null;index:idx_phone_code" json:"code"`
	Purpose   string    `gorm:"size:50;not null" json:"purpose"` // login, register, reset
	ExpiresAt time.Time `gorm:"not null;index:idx_expires" json:"expires_at"`
	Used      int8      `gorm:"default:0" json:"used"` // 0-未使用，1-已使用
	CreatedAt time.Time `json:"created_at"`
	IP        string    `gorm:"size:45" json:"ip"`
}

// Category 文章分类模型
type Category struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Slug        string    `gorm:"size:50;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"size:255" json:"description"`
	ParentID    uint64    `gorm:"default:0;index:idx_parent" json:"parent_id"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	ArticleCount int      `gorm:"default:0" json:"article_count"`
	Status      int8      `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag 文章标签模型
type Tag struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Slug         string    `gorm:"size:50;uniqueIndex;not null" json:"slug"`
	ArticleCount int       `gorm:"default:0" json:"article_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Article 文章模型
type Article struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title          string     `gorm:"size:200;not null" json:"title"`
	Slug           string     `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	Summary        string     `gorm:"size:500" json:"summary"`
	Content        string     `gorm:"type:longtext;not null" json:"content"`
	ContentHTML    string     `gorm:"type:longtext" json:"content_html"`
	CoverImage     string     `gorm:"size:255" json:"cover_image"`
	AuthorID       uint64     `gorm:"not null;index:idx_author" json:"author_id"`
	CategoryID     uint64     `gorm:"index:idx_category" json:"category_id"`
	ViewCount      uint64     `gorm:"default:0;index:idx_view_count" json:"view_count"`
	LikeCount      uint64     `gorm:"default:0" json:"like_count"`
	ShareCount     uint64     `gorm:"default:0" json:"share_count"`
	CommentCount   uint64     `gorm:"default:0" json:"comment_count"`
	Status         int8       `gorm:"default:0;index:idx_status" json:"status"` // 0-草稿，1-发布，2-下架
	IsTop          int8       `gorm:"default:0;index:idx_is_top" json:"is_top"` // 0-否，1-是
	SEOTitle       string     `gorm:"size:200" json:"seo_title"`
	SEOKeywords    string     `gorm:"size:500" json:"seo_keywords"`
	SEODescription string     `gorm:"size:500" json:"seo_description"`
	PublishedAt    *time.Time `gorm:"index:idx_published" json:"published_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at"`
	
	// 关联字段
	Category   *Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Author     *User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags       []Tag       `gorm:"many2many:article_tags;" json:"tags,omitempty"`
	Attachments []ArticleAttachment `gorm:"foreignKey:ArticleID" json:"attachments,omitempty"`
}

// ArticleTag 文章标签关联模型
type ArticleTag struct {
	ArticleID uint64    `gorm:"primaryKey;not null" json:"article_id"`
	TagID     uint64    `gorm:"primaryKey;not null;index:idx_tag" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ArticleAttachment 文章附件模型
type ArticleAttachment struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID uint64    `gorm:"not null;index:idx_article" json:"article_id"`
	FileType  string    `gorm:"size:20;not null;index:idx_file_type" json:"file_type"` // image, pdf, video
	FileName  string    `gorm:"size:255;not null" json:"file_name"`
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	FileURL   string    `gorm:"size:500;not null" json:"file_url"`
	FileSize  int64     `json:"file_size"`
	MimeType  string    `gorm:"size:100" json:"mime_type"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Duration  int       `json:"duration"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

// SensitiveWord 敏感词模型
type SensitiveWord struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Word      string    `gorm:"size:100;uniqueIndex;not null" json:"word"`
	Level     int8      `gorm:"default:1" json:"level"` // 1-低，2-中，3-高
	Category  string    `gorm:"size:50" json:"category"`
	Status    int8      `gorm:"default:1;index:idx_status" json:"status"` // 0-禁用，1-启用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Comment 评论模型
type Comment struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID uint64     `gorm:"not null;index:idx_article" json:"article_id"`
	ParentID  uint64     `gorm:"default:0;index:idx_parent" json:"parent_id"`
	UserID    *uint64    `gorm:"index" json:"user_id"`
	Nickname  string     `gorm:"size:50" json:"nickname"`
	Avatar    string     `gorm:"size:255" json:"avatar"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	ContentHTML string   `gorm:"type:text" json:"content_html"`
	IPAddress string     `gorm:"size:45" json:"ip_address"`
	IPLocation string    `gorm:"size:100" json:"ip_location"`
	UserAgent string     `gorm:"size:500" json:"user_agent"`
	Status    int8       `gorm:"default:0;index:idx_status" json:"status"` // 0-待审核，1-通过，2-拒绝
	LikeCount int        `gorm:"default:0" json:"like_count"`
	IsAdmin   int8       `gorm:"default:0" json:"is_admin"`
	CreatedAt time.Time  `gorm:"index:idx_created" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
	
	// 关联字段
	Article *Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	Replies []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// CommentLike 评论点赞记录模型
type CommentLike struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CommentID uint64    `gorm:"not null;uniqueIndex:uk_comment_ip" json:"comment_id"`
	IPAddress string    `gorm:"size:45;not null;uniqueIndex:uk_comment_ip" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

// ArticleLike 文章点赞记录模型
type ArticleLike struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID uint64    `gorm:"not null;uniqueIndex:uk_article_ip" json:"article_id"`
	IPAddress string    `gorm:"size:45;not null;uniqueIndex:uk_article_ip" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

// ArticleShare 文章分享记录模型
type ArticleShare struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID uint64    `gorm:"not null;index:idx_article" json:"article_id"`
	Platform  string    `gorm:"size:20;index:idx_platform" json:"platform"` // wechat, weibo, qq, link
	IPAddress string    `gorm:"size:45" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string    `gorm:"size:100;uniqueIndex;not null" json:"config_key"`
	ConfigValue string    `gorm:"type:text" json:"config_value"`
	ConfigType  string    `gorm:"size:20;default:string" json:"config_type"`
	GroupName   string    `gorm:"size:50;index:idx_group" json:"group_name"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OperationLog 操作日志模型
type OperationLog struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       *uint64   `gorm:"index:idx_user" json:"user_id"`
	Action       string    `gorm:"size:100;not null;index:idx_action" json:"action"`
	Module       string    `gorm:"size:50" json:"module"`
	RequestMethod string   `gorm:"size:10" json:"request_method"`
	RequestURL   string    `gorm:"size:500" json:"request_url"`
	RequestParams string  `gorm:"type:text" json:"request_params"`
	ResponseCode int       `json:"response_code"`
	IPAddress    string    `gorm:"size:45" json:"ip_address"`
	UserAgent    string    `gorm:"size:500" json:"user_agent"`
	Duration     int       `json:"duration"`
	Remark       string    `gorm:"size:500" json:"remark"`
	CreatedAt    time.Time `gorm:"index:idx_created" json:"created_at"`
}

// DailyStatistics 每日统计模型
type DailyStatistics struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	StatDate       time.Time `gorm:"type:date;uniqueIndex:idx_date;not null" json:"stat_date"`
	PageViews      uint64    `gorm:"default:0" json:"page_views"`
	UniqueVisitors uint64    `gorm:"default:0" json:"unique_visitors"`
	NewArticles    int       `gorm:"default:0" json:"new_articles"`
	NewComments    int       `gorm:"default:0" json:"new_comments"`
	TotalArticles  int       `gorm:"default:0" json:"total_articles"`
	TotalComments  int       `gorm:"default:0" json:"total_comments"`
	TotalUsers     int       `gorm:"default:0" json:"total_users"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BackupRecord 备份记录模型
type BackupRecord struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BackupType  string     `gorm:"size:20;not null;index:idx_type" json:"backup_type"` // database, file, full
	FilePath    string     `gorm:"size:500;not null" json:"file_path"`
	FileSize    int64      `json:"file_size"`
	Status      int8       `gorm:"default:0;index:idx_status" json:"status"` // 0-进行中，1-完成，2-失败
	Remark      string     `gorm:"size:500" json:"remark"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (SMSCode) TableName() string {
	return "sms_codes"
}

func (Category) TableName() string {
	return "categories"
}

func (Tag) TableName() string {
	return "tags"
}

func (Article) TableName() string {
	return "articles"
}

func (ArticleTag) TableName() string {
	return "article_tags"
}

func (ArticleAttachment) TableName() string {
	return "article_attachments"
}

func (SensitiveWord) TableName() string {
	return "sensitive_words"
}

func (Comment) TableName() string {
	return "comments"
}

func (CommentLike) TableName() string {
	return "comment_likes"
}

func (ArticleLike) TableName() string {
	return "article_likes"
}

func (ArticleShare) TableName() string {
	return "article_shares"
}

func (SystemConfig) TableName() string {
	return "system_configs"
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

func (DailyStatistics) TableName() string {
	return "daily_statistics"
}

func (BackupRecord) TableName() string {
	return "backup_records"
}
