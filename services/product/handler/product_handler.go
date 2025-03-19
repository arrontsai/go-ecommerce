package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/arrontsai/ecommerce/pkg/models"
	"github.com/arrontsai/ecommerce/services/product/service"
)

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	productService service.ProductService
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求數據: " + err.Error()})
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "創建產品失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "產品創建成功", "product": product})
}

// GetProduct handles getting a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "產品 ID 不能為空"})
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取產品失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// GetProducts handles getting all products with pagination
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Parse query parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get products
	products, total, err := h.productService.GetProducts(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取產品失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"metadata": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetProductsByCategory handles getting products by category with pagination
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	// Get category ID
	categoryID := c.Param("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "類別 ID 不能為空"})
		return
	}

	// Parse query parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get products
	products, total, err := h.productService.GetProductsByCategory(c.Request.Context(), categoryID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取產品失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"metadata": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateProduct handles updating a product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Get product ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "產品 ID 不能為空"})
		return
	}

	// Parse request body
	var req models.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求數據: " + err.Error()})
		return
	}

	// Update product
	product, err := h.productService.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新產品失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "產品更新成功", "product": product})
}

// DeleteProduct handles deleting a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Get product ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "產品 ID 不能為空"})
		return
	}

	// Delete product
	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除產品失敗: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "產品刪除成功"})
}

// RegisterRoutes registers the product routes
func (h *ProductHandler) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	products := router.Group("/api/products")
	{
		products.GET("", h.GetProducts)
		products.GET("/:id", h.GetProduct)
		products.GET("/category/:category_id", h.GetProductsByCategory)
		
		// Protected routes
		products.POST("", authMiddleware, h.CreateProduct)
		products.PUT("/:id", authMiddleware, h.UpdateProduct)
		products.DELETE("/:id", authMiddleware, h.DeleteProduct)
	}
}

