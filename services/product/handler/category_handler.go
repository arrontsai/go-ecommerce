package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/arrontsai/ecommerce/pkg/models"
	"github.com/arrontsai/ecommerce/services/product/service"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory handles category creation
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求數據: " + err.Error()})
		return
	}

	category, err := h.categoryService.CreateCategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "創建類別失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "類別創建成功", "category": category})
}

// GetCategory handles getting a category by ID
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "類別 ID 不能為空"})
		return
	}

	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取類別失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// GetCategories handles getting all categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取類別失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// UpdateCategory handles updating a category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	// Get category ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "類別 ID 不能為空"})
		return
	}

	// Parse request body
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求數據: " + err.Error()})
		return
	}

	// Update category
	category, err := h.categoryService.UpdateCategory(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新類別失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "類別更新成功", "category": category})
}

// DeleteCategory handles deleting a category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	// Get category ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "類別 ID 不能為空"})
		return
	}

	// Delete category
	if err := h.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除類別失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "類別刪除成功"})
}

// RegisterRoutes registers the category routes
func (h *CategoryHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	categories := router.Group("/api/categories")
	{
		categories.GET("", h.GetCategories)
		categories.GET("/:id", h.GetCategory)
		
		// Protected routes
		categories.POST("", authMiddleware, h.CreateCategory)
		categories.PUT("/:id", authMiddleware, h.UpdateCategory)
		categories.DELETE("/:id", authMiddleware, h.DeleteCategory)
	}
}

