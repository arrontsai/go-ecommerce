package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgresClient represents a PostgreSQL client
type PostgresClient struct {
	DB *sqlx.DB
}

// NewPostgresClient creates a new PostgreSQL client
func NewPostgresClient(connStr string) (*PostgresClient, error) {
	db, err := sqlx.ConnectContext(context.Background(), "postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("無法連接PostgreSQL: %v", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	return &PostgresClient{
		DB: db,
	}, nil
}

// Close closes the PostgreSQL connection
func (p *PostgresClient) Close() error {
	return p.DB.Close()
}

// AutoMigrate is not supported by sqlx
// func (p *PostgresClient) AutoMigrate(models ...interface{}) error {
// 	return p.DB.AutoMigrate(models...)
// }
