package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	// OrderStatusPending represents a pending order
	OrderStatusPending OrderStatus = "PENDING"
	// OrderStatusPaid represents a paid order
	OrderStatusPaid OrderStatus = "PAID"
	// OrderStatusShipped represents a shipped order
	OrderStatusShipped OrderStatus = "SHIPPED"
	// OrderStatusDelivered represents a delivered order
	OrderStatusDelivered OrderStatus = "DELIVERED"
	// OrderStatusCancelled represents a cancelled order
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	// PaymentStatusPending represents a pending payment
	PaymentStatusPending PaymentStatus = "pending"
	// PaymentStatusCompleted represents a completed payment
	PaymentStatusCompleted PaymentStatus = "completed"
	// PaymentStatusFailed represents a failed payment
	PaymentStatusFailed PaymentStatus = "failed"
	// PaymentStatusRefunded represents a refunded payment
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// Order represents an order in the system
type Order struct {
	ID            string       `json:"id" gorm:"primaryKey"`
	UserID        string       `json:"user_id" gorm:"index"`
	TotalAmount   float64      `json:"total_amount"`
	Status        OrderStatus  `json:"status"`
	Items         []OrderItem  `json:"items" gorm:"foreignKey:OrderID"`
	ShippingInfo  ShippingInfo `json:"shipping_info" gorm:"embedded"`
	PaymentInfo   PaymentInfo  `json:"payment_info" gorm:"embedded"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        string  `json:"id" gorm:"primaryKey"`
	OrderID   string  `json:"order_id" gorm:"index"`
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

// ShippingInfo represents shipping information for an order
type ShippingInfo struct {
	FullName     string `json:"full_name"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	PhoneNumber  string `json:"phone_number"`
}

// PaymentInfo represents payment information for an order
type PaymentInfo struct {
	PaymentMethod string    `json:"payment_method"`
	PaymentID     string    `json:"payment_id"`
	PaidAt        time.Time `json:"paid_at"`
}

// OrderRequest represents the data needed to create an order
type OrderRequest struct {
	Items        []OrderItemRequest `json:"items" binding:"required,dive"`
	ShippingInfo ShippingInfo       `json:"shipping_info" binding:"required"`
}

// OrderItemRequest represents the data needed to create an order item
type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// NewOrder creates a new order
func NewOrder(userID string, items []OrderItem, shippingInfo ShippingInfo, totalAmount float64) *Order {
	return &Order{
		ID:            uuid.New().String(),
		UserID:        userID,
		Status:        OrderStatusPending,
		TotalAmount:   totalAmount,
		Items:         items,
		ShippingInfo:  shippingInfo,
		PaymentStatus: PaymentStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// NewOrderItem creates a new order item
func NewOrderItem(orderID, productID, name string, price float64, quantity int) OrderItem {
	return OrderItem{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		ProductID: productID,
		Name:      name,
		Price:     price,
		Quantity:  quantity,
		Subtotal:  price * float64(quantity),
	}
}

// UpdateStatus updates the status of an order
func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
}

// UpdatePaymentStatus updates the payment status of an order
func (o *Order) UpdatePaymentStatus(status PaymentStatus) {
	o.PaymentStatus = status
	o.UpdatedAt = time.Now()
}

// SetPaymentInfo sets the payment information for an order
func (o *Order) SetPaymentInfo(method, paymentID string) {
	o.PaymentInfo = PaymentInfo{
		PaymentMethod: method,
		PaymentID:     paymentID,
		PaidAt:        time.Now(),
	}
	o.UpdatedAt = time.Now()
}
