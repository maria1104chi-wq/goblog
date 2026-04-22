package service

import (
	"blog-backend/internal/model"
	"blog-backend/internal/repository"
	"bufio"
	"os"
	"strings"
)

// SensitiveWordService 敏感词服务
type SensitiveWordService struct {
}

// NewSensitiveWordService 创建敏感词服务实例
func NewSensitiveWordService() *SensitiveWordService {
	return &SensitiveWordService{}
}

// GetSensitiveWordList 获取敏感词列表
func (s *SensitiveWordService) GetSensitiveWordList(page, pageSize int, status *int8) ([]model.SensitiveWord, int64, error) {
	var words []model.SensitiveWord
	var total int64

	db := repository.GetDB().Model(&model.SensitiveWord{})

	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&words).Error; err != nil {
		return nil, 0, err
	}

	return words, total, nil
}

// CreateSensitiveWord 创建敏感词
func (s *SensitiveWordService) CreateSensitiveWord(word *model.SensitiveWord) error {
	return repository.GetDB().Create(word).Error
}

// UpdateSensitiveWord 更新敏感词
func (s *SensitiveWordService) UpdateSensitiveWord(id uint64, updates map[string]interface{}) error {
	return repository.GetDB().Model(&model.SensitiveWord{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteSensitiveWord 删除敏感词
func (s *SensitiveWordService) DeleteSensitiveWord(id uint64) error {
	return repository.GetDB().Delete(&model.SensitiveWord{}, id).Error
}

// ImportSensitiveWordsFromFile 从文件导入敏感词
func (s *SensitiveWordService) ImportSensitiveWordsFromFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 检查是否已存在
		var existing model.SensitiveWord
		result := repository.GetDB().Where("word = ?", line).First(&existing)
		if result.Error == nil {
			// 已存在，跳过
			continue
		}

		// 创建新敏感词
		word := model.SensitiveWord{
			Word:   line,
			Level:  1,
			Status: 1,
		}
		if err := repository.GetDB().Create(&word).Error; err != nil {
			continue
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
}

// GetAllActiveWords 获取所有启用的敏感词
func (s *SensitiveWordService) GetAllActiveWords() ([]string, error) {
	var words []model.SensitiveWord
	if err := repository.GetDB().Where("status = ?", 1).Pluck("word", &words).Error; err != nil {
		return nil, err
	}
	return words, nil
}

// GetSensitiveWordCount 获取敏感词总数
func (s *SensitiveWordService) GetSensitiveWordCount() (int64, error) {
	var count int64
	err := repository.GetDB().Model(&model.SensitiveWord{}).Count(&count).Error
	return count, err
}
