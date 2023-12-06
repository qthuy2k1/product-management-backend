package repositories

import "errors"

var (
	ErrProductNotFound         = errors.New("product not found")
	ErrUserNotFound            = errors.New("user not found")
	ErrProductCategoryNotFound = errors.New("product category not found")
	ErrOrderNotFound           = errors.New("order not found")
	ErrOrderItemNotFound       = errors.New("order item not found")
	ErrNilCache                = errors.New("cache is nil")
)
