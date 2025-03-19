package repository

import (
	"context"
	"errors"
	"time"

	"github.com/arrontsai/ecommerce/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductRepository defines the interface for product repository operations
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	FindByID(ctx context.Context, id string) (*models.Product, error)
	FindBySKU(ctx context.Context, sku string) (*models.Product, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*models.Product, error)
	FindByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id string) error
	CountAll(ctx context.Context) (int64, error)
	CountByCategory(ctx context.Context, categoryID string) (int64, error)
}

// MongoProductRepository implements ProductRepository using MongoDB
type MongoProductRepository struct {
	collection *mongo.Collection
}

// NewMongoProductRepository creates a new MongoProductRepository
func NewMongoProductRepository(db *mongo.Database) ProductRepository {
	return &MongoProductRepository{
		collection: db.Collection("products"),
	}
}

// Create creates a new product in the database
func (r *MongoProductRepository) Create(ctx context.Context, product *models.Product) error {
	// Check if a product with the same SKU already exists
	existingProduct, err := r.FindBySKU(ctx, product.SKU)
	if err != nil {
		return err
	}
	if existingProduct != nil {
		return errors.New("產品 SKU 已存在")
	}

	// Insert the product
	_, err = r.collection.InsertOne(ctx, product)
	return err
}

// FindByID finds a product by ID
func (r *MongoProductRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// FindBySKU finds a product by SKU
func (r *MongoProductRepository) FindBySKU(ctx context.Context, sku string) (*models.Product, error) {
	var product models.Product
	err := r.collection.FindOne(ctx, bson.M{"sku": sku}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// FindAll finds all products with pagination
func (r *MongoProductRepository) FindAll(ctx context.Context, page, pageSize int) ([]*models.Product, error) {
	// Calculate skip
	skip := (page - 1) * pageSize

	// Set options
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Find products
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode products
	var products []*models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// FindByCategory finds products by category with pagination
func (r *MongoProductRepository) FindByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*models.Product, error) {
	// Calculate skip
	skip := (page - 1) * pageSize

	// Set options
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Find products
	cursor, err := r.collection.Find(ctx, bson.M{"category_id": categoryID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode products
	var products []*models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// Update updates a product in the database
func (r *MongoProductRepository) Update(ctx context.Context, product *models.Product) error {
	product.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": product.ID}, product)
	return err
}

// Delete deletes a product from the database
func (r *MongoProductRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// CountAll counts all products
func (r *MongoProductRepository) CountAll(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByCategory counts products by category
func (r *MongoProductRepository) CountByCategory(ctx context.Context, categoryID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"category_id": categoryID})
}

