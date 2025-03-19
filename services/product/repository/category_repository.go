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

// CategoryRepository defines the interface for category repository operations
type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	FindByID(ctx context.Context, id string) (*models.Category, error)
	FindByName(ctx context.Context, name string) (*models.Category, error)
	FindAll(ctx context.Context) ([]*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id string) error
}

// MongoCategoryRepository implements CategoryRepository using MongoDB
type MongoCategoryRepository struct {
	collection *mongo.Collection
}

// NewMongoCategoryRepository creates a new MongoCategoryRepository
func NewMongoCategoryRepository(db *mongo.Database) CategoryRepository {
	return &MongoCategoryRepository{
		collection: db.Collection("categories"),
	}
}

// Create creates a new category in the database
func (r *MongoCategoryRepository) Create(ctx context.Context, category *models.Category) error {
	// Check if a category with the same name already exists
	existingCategory, err := r.FindByName(ctx, category.Name)
	if err != nil {
		return err
	}
	if existingCategory != nil {
		return errors.New("類別名稱已存在")
	}

	// Insert the category
	_, err = r.collection.InsertOne(ctx, category)
	return err
}

// FindByID finds a category by ID
func (r *MongoCategoryRepository) FindByID(ctx context.Context, id string) (*models.Category, error) {
	var category models.Category
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindByName finds a category by name
func (r *MongoCategoryRepository) FindByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindAll finds all categories
func (r *MongoCategoryRepository) FindAll(ctx context.Context) ([]*models.Category, error) {
	// Set options
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	// Find categories
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode categories
	var categories []*models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// Update updates a category in the database
func (r *MongoCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	category.UpdatedAt = time.Now()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": category.ID}, category)
	return err
}

// Delete deletes a category from the database
func (r *MongoCategoryRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

