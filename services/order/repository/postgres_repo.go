package repository

import (
	"fmt"
	
	"github.com/arrontsai/ecommerce/services/order/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(userID string, items []model.OrderItem) error {
	// PostgreSQL交易實作
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("開始交易失敗: %w", err)
	}
	
	// 創建訂單
	orderID := uuid.New().String()
	totalPrice := 0.0
	for _, item := range items {
		totalPrice += item.Subtotal
	}
	
	_, err = tx.Exec(
		`INSERT INTO orders (order_id, user_id, total_price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		orderID, userID, totalPrice, "PENDING",
	)
	
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("插入訂單失敗: %w", err)
	}
	
	// 插入訂單項目
	for _, item := range items {
		_, err = tx.Exec(
			`INSERT INTO order_items (order_id, product_id, product_name, quantity, unit_price, subtotal)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			orderID, item.ProductID, item.ProductName, item.Quantity, item.UnitPrice, item.Subtotal,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("插入訂單項目失敗: %w", err)
		}
	}
	
	// 提交交易
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交交易失敗: %w", err)
	}
	
	return nil
}
