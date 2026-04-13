package service

import (
	"blog-backend/internal/model"
	"blog-backend/internal/repository"
	"blog-backend/internal/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
}

// NewAuthService 创建认证服务实例
func NewAuthService() *AuthService {
	return &AuthService{}
}

// SendLoginCode 发送登录验证码
func (s *AuthService) SendLoginCode(phone string, clientIP string) error {
	// 检查手机号是否注册
	var user model.User
	result := repository.GetDB().Where("phone = ?", phone).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return errors.New("该手机号未注册")
	}
	if result.Error != nil {
		return result.Error
	}

	// 生成验证码
	code, err := utils.GenerateSMSCode(6)
	if err != nil {
		return err
	}

	// 保存验证码到数据库
	expiresAt := time.Now().Add(5 * time.Minute)
	smsCode := model.SMSCode{
		Phone:     phone,
		Code:      code,
		Purpose:   "login",
		ExpiresAt: expiresAt,
		Used:      0,
		IP:        clientIP,
	}

	if err := repository.GetDB().Create(&smsCode).Error; err != nil {
		return err
	}

	// 发送短信
	if err := utils.SendLoginSMS(phone); err != nil {
		// 短信发送失败，但验证码已保存，仍然返回成功（开发环境）
		println("短信发送失败，但验证码已保存")
	}

	return nil
}

// VerifyLoginCode 验证登录验证码
func (s *AuthService) VerifyLoginCode(phone, code string) (*model.User, string, error) {
	// 查询未使用的验证码
	var smsCode model.SMSCode
	result := repository.GetDB().Where("phone = ? AND code = ? AND purpose = ? AND used = ?", 
		phone, code, "login", 0).Order("created_at DESC").First(&smsCode)
	
	if result.Error == gorm.ErrRecordNotFound {
		return nil, "", errors.New("验证码错误或已使用")
	}
	if result.Error != nil {
		return nil, "", result.Error
	}

	// 检查是否过期
	if time.Now().After(smsCode.ExpiresAt) {
		return nil, "", errors.New("验证码已过期")
	}

	// 查询用户
	var user model.User
	if err := repository.GetDB().Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, "", errors.New("用户不存在")
	}

	// 标记验证码为已使用
	repository.GetDB().Model(&smsCode).Update("used", 1)

	// 生成JWT令牌
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, 
		"blog-jwt-secret-key-2024", 24*time.Hour)
	if err != nil {
		return nil, "", err
	}

	// 更新最后登录时间
	repository.GetDB().Model(&user).Updates(map[string]interface{}{
		"last_login":     time.Now(),
		"last_login_ip":  utils.GetIPLocation(""),
	})

	return &user, token, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(userID uint64) (*model.User, error) {
	var user model.User
	if err := repository.GetDB().First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile 更新用户资料
func (s *AuthService) UpdateUserProfile(userID uint64, updates map[string]interface{}) error {
	return repository.GetDB().Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}
