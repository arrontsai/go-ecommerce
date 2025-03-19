package models

import (
	"time"

	"github.com/google/uuid"
)

// CartItem represents an item in a shopping cart
type CartItem struct {
	ProductID string    `json:"product_id" bson:"product_id"`
	Quantity  int       `json:"quantity" bson:"quantity"`
	AddedAt   time.Time `json:"added_at" bson:"added_at"`
}

// Cart represents a shopping cart in the system
type Cart struct {
	ID        string     `json:"id" bson:"_id"`
	UserID    string     `json:"user_id" bson:"user_id"`
	Items     []CartItem `json:"items" bson:"items"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
}

// CartItemRequest represents the data needed to add an item to a cart
type CartItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// UpdateCartItemRequest represents the data needed to update a cart item
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

// NewCart creates a new cart
func NewCart(userID string) *Cart {
	return &Cart{
		ID:        uuid.New().String(),
		UserID:    userID,
		Items:     []CartItem{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewCartItem creates a new cart item
func NewCartItem(productID string, quantity int) CartItem {
	return CartItem{
		ProductID: productID,
		Quantity:  quantity,
		AddedAt:   time.Now(),
	}
}

// AddItem adds an item to the cart
func (c *Cart) AddItem(item CartItem) {
	// Check if the item already exists in the cart
	for i, existingItem := range c.Items {
		if existingItem.ProductID == item.ProductID {
			// Update the quantity
			c.Items[i].Quantity += item.Quantity
			c.UpdatedAt = time.Now()
			return
		}
	}

	// Add the new item
	c.Items = append(c.Items, item)
	c.UpdatedAt = time.Now()
}

// UpdateItem updates an item in the cart
func (c *Cart) UpdateItem(productID string, quantity int) bool {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity = quantity
			c.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(productID string) bool {
	for i, item := range c.Items {
		if item.ProductID == productID {
			// Remove the item
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// Clear clears all items from the cart
func (c *Cart) Clear() {
	c.Items = []CartItem{}
	c.UpdatedAt = time.Now()
}

// TotalItems returns the total number of items in the cart
func (c *Cart) TotalItems() int {
	total := 0
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

// TotalAmount returns the total amount of the cart
func (c *Cart) TotalAmount() float64 {
	// This function is not implemented as the price of the item is not available
	// in the CartItem struct. You may need to add the price to the CartItem struct
	// or retrieve it from another source.
	return 0.0
}
