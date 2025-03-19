package repository

type CartRepository struct {
	// MongoDB 連接實作
}

func NewCartRepo(client *mongo.Client) *CartRepository {
	return &CartRepository{}
}
