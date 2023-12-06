package controllers

import "errors"

var (
	ErrUserNotFound                    = errors.New("user not found")
	ErrProductNotFound                 = errors.New("product not found")
	ErrProductCategoryNotFound         = errors.New("product category not found")
	ErrNotEnoughColumns                = errors.New("not enough columns")
	ErrIncorrectColumnNames            = errors.New("incorrect column names, the file must have columns named: Name, Description, Price, Quantity, AuthorID, Category")
	ErrCSVFileFormat                   = errors.New("the file is not a valid CSV file")
	ErrInvalidPrice                    = errors.New("price must be greater than 0 and less than 15 digits")
	ErrInvalidQuantity                 = errors.New("quantity must be greater than 0")
	ErrInvalidAuthorID                 = errors.New("invalid author id")
	ErrMissingName                     = errors.New("name cannot be blank")
	ErrNameTooLong                     = errors.New("name too long")
	ErrMissingDesc                     = errors.New("description cannot be blank")
	ErrMissingCategoryName             = errors.New("category cannot be blank")
	ErrInsufficientQuantity            = errors.New("insufficient quantity")
	ErrOrderPaidFull                   = errors.New("the order has been paid in full")
	ErrPriceInputGreaterThanOrderPrice = errors.New("the price input is greater than the order total price")
	ErrOrderNotFound                   = errors.New("order not found")
	ErrOrderItemNotFound               = errors.New("order item not found")
	ErrInvalidOrderID                  = errors.New("invalid order id")
)
