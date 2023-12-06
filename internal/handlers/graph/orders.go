package graph

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/qthuy2k1/product-management/internal/handlers/graph/model"
)

// CreateOrder is the resolver for the createOrders field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input model.OrderRequest) (bool, error) {
	order, errResp := validateAndConvertOrder(input)
	if errResp != nil {
		return true, errResp
	}

	var orderItemsInput []controllers.OrderItemInput
	for _, oi := range input.Items {
		orderItem, errResp := validateAndConvertOrderItem(*oi)
		if errResp != nil {
			return false, errResp
		}
		orderItemsInput = append(orderItemsInput, orderItem)
	}

	if err := r.Controller.CreateOrder(ctx, order, orderItemsInput); err != nil {
		log.Println(err)
		return false, convertCtrlError(err)
	}

	return true, nil
}

// UpdateOrder is the resolver for the updateOrder field.
func (r *mutationResolver) UpdateOrder(ctx context.Context, orderID int, input model.OrderRequest) (bool, error) {
	if orderID <= 0 {
		return false, ErrInvalidOrderID
	}

	order, errResp := validateAndConvertOrder(input)
	if errResp != nil {
		return false, errResp
	}

	for _, oi := range input.Items {
		orderItem, errResp := validateAndConvertOrderItem(*oi)
		if errResp != nil {
			return false, errResp
		}
		order.OrderItem = append(order.OrderItem, orderItem)
	}

	if err := r.Controller.UpdateOrder(ctx, orderID, order); err != nil {
		log.Println(err)
		return false, convertCtrlError(err)
	}
	return true, nil
}

// validateAndConvertOrder validates the order from body request and returns order struct in controller layer
func validateAndConvertOrder(orderReq model.OrderRequest) (controllers.OrderInput, error) {
	if orderReq.UserID <= 0 {
		return controllers.OrderInput{}, ErrInvalidUserID
	}
	if len(strings.TrimSpace(orderReq.Status.String())) == 0 {
		return controllers.OrderInput{}, ErrMissingOrderStatus
	}

	return controllers.OrderInput{
		UserID: orderReq.UserID,
		Status: orderReq.Status.String(),
	}, nil
}

func (r *queryResolver) GetOrders(ctx context.Context, filter *model.FilterDate, sorting *model.SortingInput, pagination model.PaginationInput) (*model.OrderResponse, error) {
	orderFilter, err := validateAndConvertOrderFilter(filter, sorting, pagination)
	if err != nil {
		return nil, err
	}

	orders, count, err := r.Controller.GetOrders(ctx, orderFilter)
	if err != nil {
		return nil, err
	}

	orderResp := make([]*model.Order, 0, len(orders))

	for _, o := range orders {
		total := o.TotalPrice.InexactFloat64()
		createdAt := o.CreatedAt.Format("02-01-2006 15:04:05")
		order := &model.Order{
			ID: o.ID,
			User: &model.User{
				Name:  o.UserName,
				Email: o.UserEmail,
			},
			Status:    model.Status(o.Status),
			Total:     &total,
			CreatedAt: &createdAt,
		}

		for index := range o.Items {
			order.Items = append(order.Items, &model.OrderItem{
				ID: o.Items[index].ID,
				Product: &model.Product{
					Name: o.Items[index].ProductName,
				},
				Quantity: o.Items[index].Quantity,
				Price:    o.Items[index].Price.InexactFloat64(),
			})
		}

		orderResp = append(orderResp, order)
	}

	orderRespGraph := model.OrderResponse{
		Order:      orderResp,
		TotalCount: int(count),
	}

	return &orderRespGraph, nil
}

func validateAndConvertOrderFilter(filter *model.FilterDate, sorting *model.SortingInput, pagination model.PaginationInput) (controllers.OrderFilterCtrl, error) {
	var orderFilter controllers.OrderFilterCtrl
	// Sorting
	for _, f := range sorting.Column {
		orderFilter.Sorting = append(orderFilter.Sorting, controllers.Sorting{
			ColumnName: f.ColumnName,
			Desc:       f.Desc,
		})
	}

	// Pagination
	if pagination.Limit > 0 {
		orderFilter.Pagination.Limit = pagination.Limit
	}
	if pagination.Page > 0 {
		orderFilter.Pagination.Page = pagination.Page
	}

	// Filter
	startDateTrimmed := strings.TrimSpace(filter.StartDate)
	endDateTrimmed := strings.TrimSpace(filter.EndDate)
	if startDateTrimmed != "" && endDateTrimmed != "" {
		startDateParsed, err := time.Parse("02-01-2006", startDateTrimmed)
		if err != nil {
			return controllers.OrderFilterCtrl{}, ErrDateBadRequest
		}

		endDateParsed, err := time.Parse("02-01-2006", endDateTrimmed)
		if err != nil {
			return controllers.OrderFilterCtrl{}, ErrDateBadRequest
		}

		if startDateParsed.After(time.Now()) || endDateParsed.After(time.Now()) {
			return controllers.OrderFilterCtrl{}, ErrDateAfterCurrentDate
		}

		if startDateParsed.After(endDateParsed) {
			return controllers.OrderFilterCtrl{}, ErrStartDateAfterEndDate
		}
	}

	orderFilter.FilterDate.StartDate = startDateTrimmed
	orderFilter.FilterDate.EndDate = endDateTrimmed

	return orderFilter, nil
}
