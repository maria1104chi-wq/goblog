package utils

import (
	"blog-backend/internal/config"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// SMSService 短信服务
type SMSService struct {
	client *dysmsapi.Client
	config *config.SMSConfig
}

// NewSMSService 创建短信服务实例
func NewSMSService(cfg *config.SMSConfig) (*SMSService, error) {
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		return nil, fmt.Errorf("阿里云短信配置不完整")
	}

	client, err := dysmsapi.NewClientWithAccessKey(
		"cn-hangzhou",
		cfg.AccessKeyID,
		cfg.AccessKeySecret,
	)
	if err != nil {
		return nil, fmt.Errorf("创建短信客户端失败：%w", err)
	}

	return &SMSService{
		client: client,
		config: cfg,
	}, nil
}

// SendVerificationCode 发送验证码短信
func (s *SMSService) SendVerificationCode(phone, code string) error {
	request := dysmsapi.CreateSendSmsRequest()
	request.PhoneNumbers = phone
	request.SignName = s.config.SignName
	request.TemplateCode = s.config.TemplateCode
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)

	response, err := s.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("发送短信失败：%w", err)
	}

	if response.Code != "OK" {
		return fmt.Errorf("短信发送失败：%s - %s", response.Code, response.Message)
	}

	return nil
}

// SendLoginCode 发送登录验证码
func (s *SMSService) SendLoginCode(phone string) (string, error) {
	code, err := GenerateSMSCode(s.config.CodeLength)
	if err != nil {
		return "", err
	}

	if err := s.SendVerificationCode(phone, code); err != nil {
		return "", err
	}

	return code, nil
}

// MockSMSService 模拟短信服务（用于开发测试）
type MockSMSService struct {
	codes map[string]string // phone -> code
}

// NewMockSMSService 创建模拟短信服务
func NewMockSMSService() *MockSMSService {
	return &MockSMSService{
		codes: make(map[string]string),
	}
}

// SendVerificationCode 发送验证码（模拟）
func (s *MockSMSService) SendVerificationCode(phone, code string) error {
	s.codes[phone] = code
	fmt.Printf("[模拟短信] 发送到 %s 的验证码：%s\n", phone, code)
	return nil
}

// SendLoginCode 发送登录验证码（模拟）
func (s *MockSMSService) SendLoginCode(phone string) (string, error) {
	code, err := GenerateSMSCode(6)
	if err != nil {
		return "", err
	}

	if err := s.SendVerificationCode(phone, code); err != nil {
		return "", err
	}

	return code, nil
}

// GetCode 获取验证码（仅用于测试）
func (s *MockSMSService) GetCode(phone string) string {
	return s.codes[phone]
}

// ClearCode 清除验证码
func (s *MockSMSService) ClearCode(phone string) {
	delete(s.codes, phone)
}

// GlobalSMSService 全局短信服务实例
var GlobalSMSService interface {
	SendLoginCode(phone string) (string, error)
	SendVerificationCode(phone, code string) error
}

// InitSMSService 初始化短信服务
func InitSMSService(cfg *config.SMSConfig) error {
	// 如果配置不完整，使用模拟服务
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" {
		fmt.Println("未配置阿里云短信参数，使用模拟短信服务")
		GlobalSMSService = NewMockSMSService()
		return nil
	}

	service, err := NewSMSService(cfg)
	if err != nil {
		fmt.Printf("初始化阿里云短信服务失败：%v，使用模拟服务\n", err)
		GlobalSMSService = NewMockSMSService()
		return nil
	}

	GlobalSMSService = service
	fmt.Println("阿里云短信服务初始化成功")
	return nil
}

// SendLoginSMS 发送登录验证码短信
func SendLoginSMS(phone string) (string, error) {
	if GlobalSMSService == nil {
		return "", fmt.Errorf("短信服务未初始化")
	}

	code, err := GlobalSMSService.SendLoginCode(phone)
	if err != nil {
		return "", err
	}

	return code, nil
}

// VerifySMSCode 验证短信验证码
func VerifySMSCode(phone, code string, expireMinutes int) bool {
	// 这里需要结合数据库中的sms_codes表进行验证
	// 实际实现应该在service层完成
	if GlobalSMSService == nil {
		return false
	}

	// 如果是模拟服务，直接验证
	if mockService, ok := GlobalSMSService.(*MockSMSService); ok {
		storedCode := mockService.GetCode(phone)
		if storedCode == "" {
			return false
		}
		
		// 简单验证，实际应该检查过期时间
		mockService.ClearCode(phone)
		return storedCode == code
	}

	// 真实服务需要通过数据库验证
	return false
}

// GetIPLocation 获取IP地理位置（简化实现）
func GetIPLocation(ip string) string {
	// 实际项目中可以调用第三方API（如淘宝IP地址库、高德地图等）
	// 这里返回一个默认值
	if ip == "127.0.0.1" || ip == "::1" {
		return "本地"
	}
	
	// 简化处理，返回"未知"
	// 实际应该调用IP地理位置查询API
	return "未知地区"
}

// CalculateExpiresAt 计算过期时间
func CalculateExpiresAt(minutes int) time.Time {
	return time.Now().Add(time.Duration(minutes) * time.Minute)
}
