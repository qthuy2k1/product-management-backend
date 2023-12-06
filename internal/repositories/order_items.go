package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type OrderItem struct {
	ID        int             `redis:"id"`
	OrderID   int             `redis:"order_id"`
	ProductID int             `redis:"product_id"`
	Quantity  int             `redis:"quantity"`
	Price     decimal.Decimal `redis:"price"`
}

// CreateOrderItem creates an order item in db given by order item model in parameter
func (r *Repository) CreateOrderItem(ctx context.Context, tx *sql.Tx, oiReq []OrderItem, order models.Order) error {
	ctxExec := boil.GetContextDB()
	if tx != nil {
		ctxExec = tx
	}

	var orderItems []*models.OrderItem
	for _, oi := range oiReq {
		orderItems = append(orderItems, &models.OrderItem{
			ProductID: oi.ProductID,
			Quantity:  oi.Quantity,
			Price:     oi.Price,
		})
	}

	if err := order.AddOrderItems(ctx, ctxExec, true, orderItems...); err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

// UpdateOrderItem updates an order item in db given by order item model inparameter
func (r *Repository) UpdateOrderItem(ctx context.Context, tx *sql.Tx, id int, oiReq OrderItem) error {
	ctxExec := boil.GetContextDB()
	if tx != nil {
		ctxExec = tx
	}

	orderItem := models.OrderItem{
		ID:        id,
		OrderID:   oiReq.OrderID,
		ProductID: oiReq.ProductID,
		Quantity:  oiReq.Quantity,
		Price:     oiReq.Price,
	}

	if _, err := orderItem.Update(ctx, ctxExec, boil.Blacklist("id", "created_at")); err != nil {
		return pkgerrors.WithStack(err)
	}

	if err := r.Redis.HSet(ctx, fmt.Sprintf("orderItem:%d", oiReq.ID), map[string]interface{}{
		"id":         oiReq.ID,
		"order_id":   oiReq.OrderID,
		"product_id": oiReq.ProductID,
		"quantity":   oiReq.Quantity,
		"price":      oiReq.Price.String(),
	}); err != nil {
		return err.Err()
	}

	return nil
}

// GetOrderItem retrieves an order item in db by ID
func (r *Repository) GetOrderItem(ctx context.Context, id int) (models.OrderItem, error) {
	res := r.Redis.HGetAll(ctx, fmt.Sprintf("orderItem:%d", id))
	if len(res.Val()) == 0 {
		orderItem, err := models.FindOrderItem(ctx, boil.GetContextDB(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.OrderItem{}, ErrOrderItemNotFound
			}
			return models.OrderItem{}, err
		}

		if errCache := r.Redis.HSet(ctx, fmt.Sprintf("orderItem:%d", id), map[string]interface{}{
			"id":         orderItem.ID,
			"order_id":   orderItem.OrderID,
			"product_id": orderItem.ProductID,
			"quantity":   orderItem.Quantity,
			"price":      orderItem.Price.String(),
		}); errCache.Err() != nil {
			return models.OrderItem{}, errCache.Err()
		}

		return *orderItem, nil
	}

	var orderItemScan OrderItem
	if err := res.Scan(&orderItemScan); err != nil {
		return models.OrderItem{}, err
	}

	return models.OrderItem{
		ID:        orderItemScan.ID,
		OrderID:   orderItemScan.OrderID,
		ProductID: orderItemScan.ProductID,
		Price:     orderItemScan.Price,
		Quantity:  orderItemScan.Quantity,
	}, nil
}
