package graph

import (
	"errors"

	"github.com/qthuy2k1/product-management/internal/controllers"
)

var (
	ErrInvalidProductID                = errors.New("invalid product id")
	ErrInvalidQuantity                 = errors.New("quantity must be non-negative")
	ErrInvalidUserID                   = errors.New("invalid user id")
	ErrMissingOrderStatus              = errors.New("missing order status")
	ErrInvalidOrderTotal               = errors.New("invalid order total price")
	ErrInsufficientQuantity            = errors.New("insufficient quantity")
	ErrProductNotFound                 = errors.New("product not found")
	ErrProductCategoryNotFound         = errors.New("product category not found")
	ErrUserNotFound                    = errors.New("user not found")
	ErrMissingName                     = errors.New("name cannot be blank")
	ErrNameTooLong                     = errors.New("name too long")
	ErrMissingDesc                     = errors.New("description cannot be blank")
	ErrMissingCategoryName             = errors.New("category cannot be blank")
	ErrInvalidPrice                    = errors.New("price must be greater than 0 and less than 15 digits")
	ErrInvalidAuthorID                 = errors.New("invalid author id")
	ErrInternalServer                  = errors.New("internal server error")
	ErrDateBadRequest                  = errors.New("invalid date format, dates must follow the format yyyy-mm-dd")
	ErrDateExpired                     = errors.New("the date has expired")
	ErrInvalidPaymentID                = errors.New("invalid payment id")
	ErrInvalidOrderID                  = errors.New("invalid order id")
	ErrOrderPaidFull                   = errors.New("the order has been paid in full")
	ErrPriceInputGreaterThanOrderPrice = errors.New("the price input is greater than the order total price")
	ErrMissingPaymentMethod            = errors.New("missing payment method")
	ErrInvalidPaymentMethod            = errors.New("invalid payment method")
	ErrOrderNotFound                   = errors.New("order not found")
	ErrOrderItemNotFound               = errors.New("order item not found")
	ErrStartDateAfterEndDate           = errors.New("start date must not be after end date")
	ErrDateAfterCurrentDate            = errors.New("date must not be after the current date")
)

// convertCtrlError compares the error return with the error in controller and returns the corresponding ErrorResponse
func convertCtrlError(err error) error {
	switch err {
	case controllers.ErrUserNotFound:
		return ErrUserNotFound
	case controllers.ErrProductCategoryNotFound:
		return ErrProductCategoryNotFound
	case controllers.ErrProductNotFound:
		return ErrProductNotFound
	case controllers.ErrInvalidPrice:
		return ErrInvalidPrice
	case controllers.ErrInvalidQuantity:
		return ErrInvalidQuantity
	case controllers.ErrInvalidAuthorID:
		return ErrInvalidAuthorID
	case controllers.ErrMissingName:
		return ErrMissingName
	case controllers.ErrMissingDesc:
		return ErrMissingDesc
	case controllers.ErrMissingCategoryName:
		return ErrMissingCategoryName
	case controllers.ErrOrderPaidFull:
		return ErrOrderPaidFull
	case controllers.ErrPriceInputGreaterThanOrderPrice:
		return ErrPriceInputGreaterThanOrderPrice
	case controllers.ErrOrderNotFound:
		return ErrOrderNotFound
	case controllers.ErrOrderItemNotFound:
		return ErrOrderItemNotFound
	case controllers.ErrInsufficientQuantity:
		return ErrInsufficientQuantity
	default:
		return ErrInternalServer
	}
}
