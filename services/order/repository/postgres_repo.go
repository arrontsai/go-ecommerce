package repository

import "github.com/jmoiron/sqlx"

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(userID string, items []OrderItem) error {
	// PostgreSQL交易實作
}
