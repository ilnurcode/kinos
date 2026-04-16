package models

import "time"

type CartItem struct{
	ProductID uint64
	ProductName string
	Quantity uint32
	Price float64
	AddedAt time.Time
}

type Cart struct{
	UserID uint64
	Items []*CartItem
	Total float64
	UpdatedAt time.Time
}

