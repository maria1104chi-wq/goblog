package repository

import (
	"blog-backend/internal/config"
	"blog-backend/internal/model"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库实例
var Database *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败：%w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败：%w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	Database = db
	return nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	if Database == nil {
		return fmt.Errorf("数据库未初始化")
	}

	models := []interface{}{
		&model.User{},
		&model.SMSCode{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		&model.ArticleTag{},
		&model.ArticleAttachment{},
		&model.SensitiveWord{},
		&model.Comment{},
		&model.CommentLike{},
		&model.ArticleLike{},
		&model.ArticleShare{},
		&model.SystemConfig{},
		&model.OperationLog{},
		&model.DailyStatistics{},
		&model.BackupRecord{},
	}

	for _, m := range models {
		if err := Database.AutoMigrate(m); err != nil {
			return fmt.Errorf("迁移表 %T 失败：%w", m, err)
		}
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return Database
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if Database != nil {
		sqlDB, err := Database.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// WithTransaction 执行事务
func WithTransaction(fn func(db *gorm.DB) error) error {
	tx := Database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// RecordOperationLog 记录操作日志
func RecordOperationLog(log *model.OperationLog) error {
	return Database.Create(log).Error
}

// UpdateDailyStatistics 更新每日统计
func UpdateDailyStatistics(date time.Time, updates map[string]interface{}) error {
	return Database.Model(&model.DailyStatistics{}).
		Where("stat_date = ?", date).
		Updates(updates).Error
}

// GetOrCreateDailyStatistics 获取或创建每日统计记录
func GetOrCreateDailyStatistics(date time.Time) (*model.DailyStatistics, error) {
	var stat model.DailyStatistics
	result := Database.Where("stat_date = ?", date).First(&stat)
	
	if result.Error == gorm.ErrRecordNotFound {
		stat.StatDate = date
		if err := Database.Create(&stat).Error; err != nil {
			return nil, err
		}
		return &stat, nil
	}
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return &stat, nil
}
