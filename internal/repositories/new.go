package repositories

import (
	"context"
	"database/sql"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/redis/go-redis/v9"
)

type IRepository interface {
	// CreateProduct creates a product using given product model in parameter
	CreateProduct(ctx context.Context, productRequest Product) error
	// GetProduct retrieves a product in db by ID
	GetProduct(ctx context.Context, id int) (models.Product, error)
	// UpdateProduct updates a product in db given by product model in parameter
	UpdateProduct(ctx context.Context, tx *sql.Tx, pReq models.Product) error
	// DeleteProduct deletes a product in db by ID
	DeleteProduct(ctx context.Context, id int) error
	// GetProducts retrieves all the products in db
	GetProducts(ctx context.Context, filter ProductRepoFilter) ([]ProductOutput, error)
	// UpsertProducts updates list of products. If a product already exists in db, updates only changed values instead
	UpsertProducts(ctx context.Context, products []Product) error
	// GetProductsGraph retrieves all products in db and is used in GraphQL.
	GetProductsGraph(ctx context.Context, filter ProductRepoFilter) ([]GetProductsGraph, error)

	// CreateUser creates a user using the given user model in parameter
	CreateUser(ctx context.Context, user User) error
	// GetUser retrieves a user by user id
	GetUser(ctx context.Context, id int) (models.User, error)
	// GetUserDefault gets the id of default user
	GetUserDefault(ctx context.Context) (models.User, error)

	// CreateProductCategory creates a product category using given product category model in parameter
	CreateProductCategory(ctx context.Context, productCategory ProductCategory) error
	// GetProductCategory gets a product category from db by product category id
	GetProductCategory(ctx context.Context, id int) (models.ProductCategory, error)
	// GetProductCategoryByName retrieves a product category from db by name
	GetProductCategoryByName(ctx context.Context, name string) (models.ProductCategory, error)
	// GetProductCategoryDefault gets the id of default user
	GetProductCategoryDefault(ctx context.Context) (models.ProductCategory, error)

	// CreateOrderItem creates an order item in db given by order item model in parameter
	CreateOrderItem(ctx context.Context, tx *sql.Tx, oiReq []OrderItem, order models.Order) error
	// UpdateOrderItem updates an order item in db given by order item model inparameter
	UpdateOrderItem(ctx context.Context, tx *sql.Tx, id int, oiReq OrderItem) error
	// GetOrderItem retrieves an order item in db by ID
	GetOrderItem(ctx context.Context, id int) (models.OrderItem, error)

	// CreateOrder creates an order in db given by order model in parameter
	CreateOrder(ctx context.Context, tx *sql.Tx, oReq Order) (models.Order, error)
	// UpdateOrder updates an order in db given by product model in parameter
	UpdateOrder(ctx context.Context, tx *sql.Tx, orderReq models.Order) error
	// GetOrder retrieves an order in db by id
	GetOrder(ctx context.Context, orderID int) (models.Order, error)
	// GetOrders retrieves all of order in db
	GetOrders(ctx context.Context, filter OrderFilterRepo) ([]OrderOutputGraph, int64, error)

	// BeginTx begins a transaction with the current global database handle
	BeginTx(ctx context.Context) (*sql.Tx, error)
	// RollbackTx aborts the transaction
	RollbackTx(tx *sql.Tx) error
	// CommitTx commits the transaction
	CommitTx(tx *sql.Tx) error
}

type Repository struct {
	Database *sql.DB
	Redis    *redis.Client
}

func NewRepository(db *sql.DB, redis *redis.Client) IRepository {
	return &Repository{Database: db, Redis: redis}
}
