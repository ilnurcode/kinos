package models

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipping  OrderStatus = "shipping"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID              uint64
	UserID          uint64
	Items           []OrderItem
	Total           float64
	Status          OrderStatus
	DeliveryAddress string
	Phone           string
	Comment         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ID          uint64
	OrderID     uint64
	ProductID   uint64
	ProductName string
	Quantity    uint32
	Price       float64
	Subtotal    float64
}
