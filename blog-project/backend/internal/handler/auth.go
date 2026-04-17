package handler

import (
	"blog-backend/internal/config"
	"blog-backend/internal/middleware"
	"blog-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
	config      *config.AppConfig
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(cfg *config.AppConfig) *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
		config:      cfg,
	}
}

// SendSMSCode 发送短信验证码
// @Summary 发送登录短信验证码
// @Tags 认证
// @Accept json
// @Produce json
// @Param phone body string true "手机号"
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/sms [post]
func (h *AuthHandler) SendSMSCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required,len=11"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	clientIP := middleware.GetClientIP(c)
	if err := h.authService.SendLoginCode(req.Phone, clientIP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "验证码已发送",
	})
}

// Login 登录
// @Summary 管理员登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param phone body string true "手机号"
// @Param code body string true "验证码"
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required,len=11"`
		Code  string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	user, token, err := h.authService.VerifyLoginCode(req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/user [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	user, err := h.authService.GetUserByID(userID.(uint64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    user,
	})
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body object true "用户资料"
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if err := h.authService.UpdateUserProfile(userID.(uint64), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}
