package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus 訂單狀態枚舉
type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"   // 待處理
	StatusPaid      OrderStatus = "PAID"      // 已支付
	StatusShipped   OrderStatus = "SHIPPED"   // 已發貨
	StatusDelivered OrderStatus = "DELIVERED" // 已送達
	StatusCancelled OrderStatus = "CANCELLED" // 已取消
)

// Order 訂單模型
type Order struct {
	ID         string      `json:"id" bson:"_id"`
	UserID     string      `json:"user_id" bson:"user_id"`
	Items      []OrderItem `json:"items" bson:"items"`
	TotalPrice float64     `json:"total_price" bson:"total_price"`
	Status     OrderStatus `json:"status" bson:"status"`
	CreatedAt  time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" bson:"updated_at"`
}

// OrderItem 訂單商品項目
type OrderItem struct {
	ProductID   string  `json:"product_id" bson:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	UnitPrice   float64 `json:"unit_price" bson:"unit_price"`
	Subtotal    float64 `json:"subtotal" bson:"subtotal"`
}

// NewOrder 創建新訂單
func NewOrder(userID string, items []OrderItem) *Order {
	totalPrice := 0.0
	for _, item := range items {
		totalPrice += item.Subtotal
	}

	return &Order{
		ID:         uuid.NewString(),
		UserID:     userID,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
