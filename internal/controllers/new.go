package controllers

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/qthuy2k1/product-management/internal/models"

	"github.com/qthuy2k1/product-management/internal/repositories"
)

type IController interface {
	// CreateProduct creates a product in db given by product model in parameter
	CreateProduct(ctx context.Context, productInput ProductInput) error
	// UpdateProduct updates a product in db given by product model in parameter
	UpdateProduct(ctx context.Context, pInput ProductInput) error
	// DeleteProduct deletes a product in db by ID
	DeleteProduct(ctx context.Context, id int) error
	// GetProducts retrieves all the products in db
	GetProducts(ctx context.Context, filter ProductCtrlFilter) ([]ProductOutput, error)
	// ImportProductsFromCSV imports list of products data from a CSV file
	ImportProductsFromCSV(ctx context.Context, file multipart.File) error
	// GetProductsGraph retrieves all the products and their author, category
	GetProductsGraph(ctx context.Context, pFilter ProductCtrlFilter) ([]ProductOutputGraph, error)
	// SendEmailProduct send an email to a list of user gmail with a list of products csv file attachment
	SendEmailProduct(emailToList []string, reader io.Reader) error
	// ExportProductsToCSV exports all the products in db to a csv file
	ExportProductsToCSV(ctx context.Context, filter ProductCtrlFilter) (io.Reader, error)

	// CreateUser adds a user to database
	CreateUser(ctx context.Context, user UserInput) error
	// GetUser gets a user from database by ID
	GetUser(ctx context.Context, id int) (UserOutput, error)

	// CreateProductCategory adds a new category to the database
	CreateProductCategory(ctx context.Context, pCateInput PCateInput) error
	// GetProductCategoryByName retrieves a product category by name
	GetProductCategoryByName(ctx context.Context, name string) (PCateOutput, error)

	// CreateOrder creates an order in db given by order model in parameter
	CreateOrder(ctx context.Context, orderInput OrderInput, orderItemsInput []OrderItemInput) error
	// SendEmailOrder sends an email that contains the order detail to the user
	SendEmailOrder(ctx context.Context, emailTo string, order models.Order, orderItem []repositories.OrderItem) error
	// UpdateOrder updates an order in db given by order model in parameter
	UpdateOrder(ctx context.Context, orderID int, orderInput OrderInput) error
	// GetOrders retrieves all the orders in db
	GetOrders(ctx context.Context, filter OrderFilterCtrl) ([]OrderOutputGraph, int64, error)
}

type Controller struct {
	Repository repositories.IRepository
}

// NewController creates new Controller
func NewController(repository repositories.IRepository) IController {
	return &Controller{Repository: repository}
}
