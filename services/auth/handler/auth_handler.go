package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ecommerce/pkg/models"
	"github.com/yourusername/ecommerce/services/auth/service"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求數據: " + err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "註冊失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "用戶註冊成功", "user": user})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求數據: " + err.Error()})
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "登入失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登入成功",
		"token":   token,
		"user":    user,
	})
}

// GetMe handles getting the current user
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授權"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取用戶信息失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// RegisterRoutes registers the authentication routes
func (h *AuthHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.GET("/me", authMiddleware, h.GetMe)
	}
}
