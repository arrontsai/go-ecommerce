package models

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the system
type Product struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	SKU         string    `json:"sku" bson:"sku"`
	CategoryID  string    `json:"category_id" bson:"category_id"`
	Inventory   int       `json:"inventory" bson:"inventory"`
	Images      []string  `json:"images" bson:"images"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Category represents a product category
type Category struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// ProductRequest represents the data needed to create or update a product
type ProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	SKU         string   `json:"sku" binding:"required"`
	CategoryID  string   `json:"category_id" binding:"required"`
	Inventory   int      `json:"inventory" binding:"required,gte=0"`
	Images      []string `json:"images"`
}

// CategoryRequest represents the data needed to create or update a category
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// NewProduct creates a new product
func NewProduct(name, description string, price float64, sku, categoryID string, inventory int, images []string) *Product {
	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		SKU:         sku,
		CategoryID:  categoryID,
		Inventory:   inventory,
		Images:      images,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewCategory creates a new category
func NewCategory(name, description string) *Category {
	return &Category{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// UpdateProduct updates a product with the provided data
func (p *Product) UpdateProduct(req ProductRequest) {
	p.Name = req.Name
	p.Description = req.Description
	p.Price = req.Price
	p.SKU = req.SKU
	p.CategoryID = req.CategoryID
	p.Inventory = req.Inventory
	p.Images = req.Images
	p.UpdatedAt = time.Now()
}

// UpdateCategory updates a category with the provided data
func (c *Category) UpdateCategory(req CategoryRequest) {
	c.Name = req.Name
	c.Description = req.Description
	c.UpdatedAt = time.Now()
}
