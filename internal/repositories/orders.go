package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Order struct {
	ID         int
	UserID     int
	Status     string
	TotalPrice decimal.NullDecimal
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CreateOrder creates an order in db given by order model in parameter
func (r *Repository) CreateOrder(ctx context.Context, tx *sql.Tx, oReq Order) (models.Order, error) {
	ctxExec := boil.GetContextDB()
	if tx != nil {
		ctxExec = tx
	}

	order := models.Order{
		UserID:     oReq.UserID,
		Status:     oReq.Status,
		TotalPrice: oReq.TotalPrice,
	}
	if err := order.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

// UpdateOrder updates an order in db given by product model in parameter
func (r *Repository) UpdateOrder(ctx context.Context, tx *sql.Tx, oReq models.Order) error {
	ctxExec := boil.GetContextDB()
	if tx != nil {
		ctxExec = tx
	}

	if _, err := oReq.Update(ctx, ctxExec, boil.Blacklist("id", "created_at")); err != nil {
		return err
	}

	if err := r.Redis.HSet(ctx, fmt.Sprintf("order:%d", oReq.ID), map[string]interface{}{
		"id":          oReq.ID,
		"user_id":     oReq.UserID,
		"status":      oReq.Status,
		"total_price": oReq.TotalPrice.Decimal.String(),
		"created_at":  oReq.CreatedAt,
		"updated_at":  oReq.UpdatedAt,
	}); err != nil {
		return err.Err()
	}

	return nil
}

// GetOrder retrieves an order in db by id
func (r *Repository) GetOrder(ctx context.Context, orderID int) (models.Order, error) {
	res := r.Redis.HGetAll(ctx, fmt.Sprintf("order:%d", orderID))
	if len(res.Val()) == 0 {
		order, err := models.FindOrder(ctx, boil.GetContextDB(), orderID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.Order{}, ErrOrderNotFound
			}
			return models.Order{}, err
		}

		if errCache := r.Redis.HSet(ctx, fmt.Sprintf("order:%d", orderID), map[string]interface{}{
			"id":          order.ID,
			"user_id":     order.UserID,
			"status":      order.Status,
			"total_price": order.TotalPrice.Decimal.String(),
			"created_at":  order.CreatedAt,
			"updated_at":  order.UpdatedAt,
		}); errCache.Err() != nil {
			return models.Order{}, errCache.Err()
		}

		return *order, nil
	}

	var orderScan Order
	if err := res.Scan(&orderScan); err != nil {
		return models.Order{}, err
	}

	return models.Order{
		ID:         orderScan.ID,
		UserID:     orderScan.UserID,
		Status:     orderScan.Status,
		TotalPrice: orderScan.TotalPrice,
		CreatedAt:  orderScan.CreatedAt,
		UpdatedAt:  orderScan.UpdatedAt,
	}, nil
}

type OrderOutputGraph struct {
	ID          int             `boil:"orders.id"`
	UserName    string          `boil:"users.name"`
	UserEmail   string          `boil:"users.email"`
	Status      string          `boil:"orders.status"`
	TotalPrice  decimal.Decimal `boil:"orders.total_price"`
	CreatedAt   time.Time       `boil:"orders.created_at"`
	ItemID      string          `boil:"items_id"`
	ProductName string          `boil:"products_name"`
	Quantity    string          `boil:"items_quantity"`
	ItemPrice   string          `boil:"items_price"`
}

type OrderFilterRepo struct {
	Sorting    []Sorting
	Pagination Pagination
	FilterDate FilterDate
}

type Pagination struct {
	Limit int
	Page  int
}

type Sorting struct {
	ColumnName string
	SortOrder  string
}

type FilterDate struct {
	StartDate string
	EndDate   string
}

// GetOrders retrieves all of order in db
func (r *Repository) GetOrders(ctx context.Context, filter OrderFilterRepo) ([]OrderOutputGraph, int64, error) {
	orderTable := models.TableNames.Orders
	orderItemTable := models.TableNames.OrderItems
	productTable := models.TableNames.Products
	userTable := models.TableNames.Users

	queryMod := []qm.QueryMod{
		// SELECT
		qm.Select(
			fmt.Sprintf("%s.id", orderTable),
			fmt.Sprintf("%s.status", orderTable),
			fmt.Sprintf("%s.total_price", orderTable),
			fmt.Sprintf("%s.created_at", orderTable),
			fmt.Sprintf("%s.name", userTable),
			fmt.Sprintf("%s.email", userTable),
			fmt.Sprintf(`array_agg(%s.id) as "items_id"`, orderItemTable),
			fmt.Sprintf(`array_agg(%s.name) as "products_name"`, productTable),
			fmt.Sprintf(`array_agg(%s.price) as "items_price"`, orderItemTable),
			fmt.Sprintf(`array_agg(%s.quantity) as "items_quantity"`, orderItemTable),
		),
		// GROUP BY
		qm.GroupBy(fmt.Sprintf("%s.id, %s.name, %s.email", orderTable, userTable, userTable)),

		// INNER JOIN orders and order_items
		qm.InnerJoin(fmt.Sprintf("%s ON %s.%s=%s.%s", orderItemTable, orderItemTable, models.OrderItemColumns.OrderID, orderTable, models.OrderColumns.ID)),

		// INNER JOIN order_items and products
		qm.InnerJoin(fmt.Sprintf("%s ON %s.%s=%s.%s", productTable, orderItemTable, models.OrderItemColumns.ProductID, productTable, models.ProductColumns.ID)),

		// INNER JOIN orders and users
		qm.InnerJoin(fmt.Sprintf("%s ON %s.%s=%s.%s", userTable, orderTable, models.OrderColumns.UserID, userTable, models.UserColumns.ID)),
	}

	// filter the date range by start and end
	if filter.FilterDate.StartDate != "" && filter.FilterDate.EndDate != "" {
		queryMod = append(queryMod, qm.Where(fmt.Sprintf("%s.created_at BETWEEN TO_TIMESTAMP('%s 00:00:00', 'DD MM YYYY HH24:MI:SS') AND TO_TIMESTAMP('%s 23:59:59', 'DD MM YYYY HH24:MI:SS')", orderTable, filter.FilterDate.StartDate, filter.FilterDate.EndDate)))
	}

	// sorting
	var orderByQuery []string
	for _, s := range filter.Sorting {
		orderByQuery = append(orderByQuery, fmt.Sprintf("%s.%s %s", orderTable, s.ColumnName, s.SortOrder))
	}
	queryMod = append(queryMod, qm.OrderBy(strings.Join(orderByQuery, ",")))

	// pagination
	queryMod = append(queryMod, qm.Limit(filter.Pagination.Limit))
	queryMod = append(queryMod, qm.Offset((filter.Pagination.Page-1)*filter.Pagination.Limit))

	var orders []OrderOutputGraph
	if err := models.Orders(queryMod...).Bind(ctx, boil.GetContextDB(), &orders); err != nil {
		return nil, 0, err
	}

	totalCount, err := models.Orders().Count(ctx, boil.GetContextDB())
	if err != nil {
		return nil, 0, err
	}

	return orders, totalCount, nil
}

// BeginTx begins a transaction with the current global database handle
func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return &sql.Tx{}, err
	}

	return tx, nil
}

// RollbackTx aborts the transaction
func (r *Repository) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

// CommitTx commits the transaction
func (r *Repository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}
