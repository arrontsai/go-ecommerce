package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/arrontsai/ecommerce/pkg/models"
	"github.com/arrontsai/ecommerce/pkg/database"
)

// CartRepository 定義購物車儲存庫的介面
type CartRepository interface {
	AddToCart(userID, productID string, quantity int) error
	GetCart(userID string) (*models.Cart, error)
	ClearCart(userID string) error
}

// MongoCartRepository 實現基於MongoDB的購物車儲存庫
type MongoCartRepository struct {
	collection *mongo.Collection
}

// NewMongoCartRepository 創建一個新的MongoDB購物車儲存庫
func NewMongoCartRepository(client *database.MongoClient) *MongoCartRepository {
	collection := client.Collection("carts")
	return &MongoCartRepository{collection: collection}
}

// AddToCart 添加商品到購物車
func (r *MongoCartRepository) AddToCart(userID, productID string, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 檢查購物車是否存在
	var cart models.Cart
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	
	if err == mongo.ErrNoDocuments {
		// 創建新購物車
		cart = models.Cart{
			ID:        uuid.New().String(),
			UserID:    userID,
			Items:     []models.CartItem{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	} else if err != nil {
		return err
	}

	// 檢查商品是否已在購物車中
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity += quantity
			found = true
			break
		}
	}

	// 如果商品不在購物車中，添加新商品
	if !found {
		cart.Items = append(cart.Items, models.CartItem{
			ProductID: productID,
			Quantity:  quantity,
			AddedAt:   time.Now(),
		})
	}

	cart.UpdatedAt = time.Now()

	// 更新或插入購物車
	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": cart},
		opts,
	)

	return err
}

// GetCart 獲取用戶的購物車
func (r *MongoCartRepository) GetCart(userID string) (*models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart models.Cart
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

// ClearCart 清空用戶的購物車
func (r *MongoCartRepository) ClearCart(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}
