package controllers

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderItemInput struct {
	ID        int
	ProductID int
	Quantity  int
}

type OrderItemOutput struct {
	ID          int
	ProductName string
	Quantity    int
	Price       decimal.Decimal
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
