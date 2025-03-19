package service

import (
	"context"
	"errors"

	"github.com/arrontsai/ecommerce/pkg/models"
	"github.com/arrontsai/ecommerce/services/product/repository"
)

// ProductService defines the interface for product service operations
type ProductService interface {
	CreateProduct(ctx context.Context, req models.ProductRequest) (*models.Product, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetProducts(ctx context.Context, page, pageSize int) ([]*models.Product, int64, error)
	GetProductsByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*models.Product, int64, error)
	UpdateProduct(ctx context.Context, id string, req models.ProductRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

// DefaultProductService implements ProductService
type DefaultProductService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService creates a new ProductService
func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &DefaultProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateProduct creates a new product
func (s *DefaultProductService) CreateProduct(ctx context.Context, req models.ProductRequest) (*models.Product, error) {
	// Check if the category exists
	category, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("類別不存在")
	}

	// Create the product
	product := models.NewProduct(req.Name, req.Description, req.Price, req.SKU, req.CategoryID, req.Inventory, req.Images)

	// Save the product
	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID gets a product by ID
func (s *DefaultProductService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	// Find the product
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("產品不存在")
	}

	return product, nil
}

// GetProducts gets all products with pagination
func (s *DefaultProductService) GetProducts(ctx context.Context, page, pageSize int) ([]*models.Product, int64, error) {
	// Get the total count
	total, err := s.productRepo.CountAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get the products
	products, err := s.productRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetProductsByCategory gets products by category with pagination
func (s *DefaultProductService) GetProductsByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*models.Product, int64, error) {
	// Check if the category exists
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, 0, err
	}
	if category == nil {
		return nil, 0, errors.New("類別不存在")
	}

	// Get the total count
	total, err := s.productRepo.CountByCategory(ctx, categoryID)
	if err != nil {
		return nil, 0, err
	}

	// Get the products
	products, err := s.productRepo.FindByCategory(ctx, categoryID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// UpdateProduct updates a product
func (s *DefaultProductService) UpdateProduct(ctx context.Context, id string, req models.ProductRequest) (*models.Product, error) {
	// Check if the product exists
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("產品不存在")
	}

	// Check if the category exists
	category, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("類別不存在")
	}

	// Update the product
	product.UpdateProduct(req)

	// Save the product
	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product
func (s *DefaultProductService) DeleteProduct(ctx context.Context, id string) error {
	// Check if the product exists
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("產品不存在")
	}

	// Delete the product
	return s.productRepo.Delete(ctx, id)
}

