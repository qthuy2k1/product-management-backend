package graph

import (
	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/qthuy2k1/product-management/internal/handlers/graph/model"
)

// validateAndConvertOrderItem validates the order item from body request and returns order item struct in controller layer
func validateAndConvertOrderItem(oiReq model.OrderItemRequest) (controllers.OrderItemInput, error) {
	if oiReq.ProductID <= 0 {
		return controllers.OrderItemInput{}, ErrInvalidProductID
	}

	if oiReq.Quantity <= 0 {
		return controllers.OrderItemInput{}, ErrInvalidQuantity
	}

	oiInput := controllers.OrderItemInput{
		ProductID: oiReq.ProductID,
		Quantity:  oiReq.Quantity,
	}

	if oiReq.ID != nil {
		if *oiReq.ID <= 0 {
			return controllers.OrderItemInput{}, ErrInvalidOrderID
		}
		oiInput.ID = *oiReq.ID
	}

	return oiInput, nil
}
