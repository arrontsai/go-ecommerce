package service

import (
	"context"
	"errors"

	"github.com/yourusername/ecommerce/pkg/models"
	"github.com/yourusername/ecommerce/services/product/repository"
)

// CategoryService defines the interface for category service operations
type CategoryService interface {
	CreateCategory(ctx context.Context, req models.CategoryRequest) (*models.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*models.Category, error)
	GetCategories(ctx context.Context) ([]*models.Category, error)
	UpdateCategory(ctx context.Context, id string, req models.CategoryRequest) (*models.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

// DefaultCategoryService implements CategoryService
type DefaultCategoryService struct {
	categoryRepo repository.CategoryRepository
	productRepo  repository.ProductRepository
}

// NewCategoryService creates a new CategoryService
func NewCategoryService(categoryRepo repository.CategoryRepository, productRepo repository.ProductRepository) CategoryService {
	return &DefaultCategoryService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}

// CreateCategory creates a new category
func (s *DefaultCategoryService) CreateCategory(ctx context.Context, req models.CategoryRequest) (*models.Category, error) {
	// Check if the category already exists
	existingCategory, err := s.categoryRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingCategory != nil {
		return nil, errors.New("類別名稱已存在")
	}

	// Create the category
	category := models.NewCategory(req.Name, req.Description)

	// Save the category
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID gets a category by ID
func (s *DefaultCategoryService) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	// Find the category
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("類別不存在")
	}

	return category, nil
}

// GetCategories gets all categories
func (s *DefaultCategoryService) GetCategories(ctx context.Context) ([]*models.Category, error) {
	return s.categoryRepo.FindAll(ctx)
}

// UpdateCategory updates a category
func (s *DefaultCategoryService) UpdateCategory(ctx context.Context, id string, req models.CategoryRequest) (*models.Category, error) {
	// Check if the category exists
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("類別不存在")
	}

	// Check if the name is already taken by another category
	if category.Name != req.Name {
		existingCategory, err := s.categoryRepo.FindByName(ctx, req.Name)
		if err != nil {
			return nil, err
		}
		if existingCategory != nil && existingCategory.ID != id {
			return nil, errors.New("類別名稱已存在")
		}
	}

	// Update the category
	category.UpdateCategory(req)

	// Save the category
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category
func (s *DefaultCategoryService) DeleteCategory(ctx context.Context, id string) error {
	// Check if the category exists
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("類別不存在")
	}

	// Check if there are products in this category
	count, err := s.productRepo.CountByCategory(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("無法刪除類別，因為有產品屬於此類別")
	}

	// Delete the category
	return s.categoryRepo.Delete(ctx, id)
}
