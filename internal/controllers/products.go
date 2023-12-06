package controllers

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/qthuy2k1/product-management/internal/repositories"
	"github.com/qthuy2k1/product-management/internal/utils/email"
	"github.com/shopspring/decimal"
)

type ProductInput struct {
	ID           int
	Name         string
	Description  string
	Price        decimal.Decimal
	Quantity     int
	AuthorID     int
	CategoryName string
}

// CreateProduct creates a product in db given by product model in parameter
func (c *Controller) CreateProduct(ctx context.Context, pInput ProductInput) error {
	// check user exists
	if _, err := c.Repository.GetUser(ctx, pInput.AuthorID); err != nil {
		return err
	}

	// check product category exists
	pCate, err := c.Repository.GetProductCategoryByName(ctx, pInput.CategoryName)
	if err != nil {
		return err
	}

	product := repositories.Product{
		Name:        pInput.Name,
		Description: pInput.Description,
		Price:       pInput.Price,
		Quantity:    pInput.Quantity,
		AuthorID:    pInput.AuthorID,
		CategoryID:  pCate.ID,
	}

	return c.Repository.CreateProduct(ctx, product)
}

// UpdateProduct updates a product in db given by product model in parameter
func (c *Controller) UpdateProduct(ctx context.Context, pInput ProductInput) error {
	// check product exists
	product, err := c.Repository.GetProduct(ctx, pInput.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return ErrProductNotFound
		}
		return err
	}

	// check user exists
	if _, err := c.Repository.GetUser(ctx, pInput.AuthorID); err != nil {
		return err
	}

	// check product category exists
	pCate, err := c.Repository.GetProductCategoryByName(ctx, pInput.CategoryName)
	if err != nil {
		return err
	}

	product.Name = pInput.Name
	product.Description = pInput.Description
	product.Price = pInput.Price
	product.Quantity = pInput.Quantity
	product.AuthorID = pInput.AuthorID
	product.CategoryID = pCate.ID

	return c.Repository.UpdateProduct(ctx, nil, product)
}

// DeleteProduct deletes a product in db by ID
func (c *Controller) DeleteProduct(ctx context.Context, id int) error {
	if err := c.Repository.DeleteProduct(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

type ProductCtrlFilter struct {
	Name  string
	Date  string
	Email []string
}

type ProductOutput struct {
	ID           int
	Name         string
	Description  string
	Price        decimal.Decimal
	Quantity     int
	AuthorID     int
	CategoryName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// GetProducts retrieves all the products in db
func (c *Controller) GetProducts(ctx context.Context, filter ProductCtrlFilter) ([]ProductOutput, error) {
	pRepoFilter := repositories.ProductRepoFilter{
		Name: filter.Name,
		Date: filter.Date,
	}

	products, err := c.Repository.GetProducts(ctx, pRepoFilter)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	var pListResp []ProductOutput
	for _, product := range products {
		pListResp = append(pListResp, ProductOutput{
			ID:           product.ID,
			Name:         product.Name,
			Description:  product.Description,
			Price:        product.Price,
			Quantity:     product.Quantity,
			AuthorID:     product.AuthorID,
			CategoryName: product.CategoryName,
			CreatedAt:    product.CreatedAt,
			UpdatedAt:    product.UpdatedAt,
		})
	}
	return pListResp, nil
}

type ProductIndexHeader struct {
	Name        int
	Description int
	Price       int
	Quantity    int
	AuthorID    int
	Category    int
}

type ProductCSVInput struct {
	Name         string
	Description  string
	Price        string
	Quantity     string
	AuthorID     string
	CategoryName string
}

// ImportProductsFromCSV imports list of products data from a CSV file
func (c *Controller) ImportProductsFromCSV(ctx context.Context, file multipart.File) error {
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return ErrCSVFileFormat
	}
	// check if the file doesn't have enough columns
	if (reflect.TypeOf(ProductInput{}).NumField()-1 > len(records[0])) {
		return ErrNotEnoughColumns
	}

	var pHeader ProductIndexHeader
	for i, r := range records[0] {
		switch r {
		case "Name":
			pHeader.Name = i
		case "Description":
			pHeader.Description = i
		case "Price":
			pHeader.Price = i
		case "Quantity":
			pHeader.Quantity = i
		case "AuthorID":
			pHeader.AuthorID = i
		case "Category":
			pHeader.Category = i
		}
	}

	// get the user default
	userDefault, err := c.Repository.GetUserDefault(ctx)
	if err != nil {
		return err
	}

	// get product category default
	pCateDefault, err := c.Repository.GetProductCategoryDefault(ctx)
	if err != nil {
		return err
	}

	// skip the first index, which is the header
	records = records[1:]

	// split records into chunks of 1000 rows each
	chunkSize := 1000
	numChunks := (len(records) + chunkSize - 1) / chunkSize // round up division

	productsInput := []repositories.Product{}
	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(records) {
			end = len(records)
		}

		// iterate the records and append it to product controller input
		for _, record := range records[start:end] {
			product, err := c.validateAndConvertProductCSV(ctx, ProductCSVInput{
				Name:         record[pHeader.Name],
				Description:  record[pHeader.Description],
				Price:        record[pHeader.Price],
				Quantity:     record[pHeader.Quantity],
				AuthorID:     record[pHeader.AuthorID],
				CategoryName: record[pHeader.Category],
			}, userDefault.ID, pCateDefault.ID)
			if err != nil {
				log.Println(err, "at product:", record[pHeader.Name])
				continue
			}

			productsInput = append(productsInput, product)
		}
	}

	if err := c.Repository.UpsertProducts(ctx, productsInput); err != nil {
		return err
	}

	return nil
}

func (c *Controller) validateAndConvertProductCSV(ctx context.Context, product ProductCSVInput, userIDDefault int, pCateIDDefault int) (repositories.Product, error) {
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)

	if len(product.Name) == 0 {
		return repositories.Product{}, ErrMissingName
	}

	if len(product.Name) > 255 {
		return repositories.Product{}, ErrNameTooLong
	}

	if len(product.Description) == 0 {
		return repositories.Product{}, ErrMissingDesc
	}

	price, err := decimal.NewFromString(product.Price)
	if err != nil {
		return repositories.Product{}, err
	}
	if price.IsZero() && price.IsNegative() && len(product.Price) > 15 {
		return repositories.Product{}, ErrInvalidPrice
	}

	quantity, err := strconv.Atoi(product.Quantity)
	if err != nil {
		return repositories.Product{}, ErrInvalidQuantity
	}

	if quantity < 0 {
		return repositories.Product{}, ErrInvalidQuantity
	}

	authorID, err := strconv.Atoi(product.AuthorID)
	if err != nil {
		return repositories.Product{}, ErrInvalidAuthorID
	}

	if authorID <= 0 {
		return repositories.Product{}, ErrInvalidAuthorID
	}

	if len(strings.TrimSpace(product.CategoryName)) == 0 {
		return repositories.Product{}, ErrMissingCategoryName
	}

	// check user exists in cache
	if _, err := c.Repository.GetUser(ctx, authorID); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			authorID = userIDDefault
		} else {
			return repositories.Product{}, err
		}
	}

	// check product category exists in cache
	var pCateID int
	pCate, err := c.Repository.GetProductCategoryByName(ctx, product.CategoryName)
	if err != nil {
		if errors.Is(err, repositories.ErrProductCategoryNotFound) {
			pCateID = pCateIDDefault
		} else {
			return repositories.Product{}, err
		}
	} else {
		pCateID = pCate.ID
	}

	// remove any unexpected characters that may exist in the subStr slice
	subStr := []string{`<`, `>`, `^`, `\`, `""`, `\"`}
	for _, x := range subStr {
		if strings.Contains(product.Name, x) {
			product.Name = strings.ReplaceAll(product.Name, x, "")
		}

		if strings.Contains(product.Description, x) {
			product.Description = strings.ReplaceAll(product.Description, x, "")
		}
	}

	return repositories.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       price,
		Quantity:    quantity,
		AuthorID:    authorID,
		CategoryID:  pCateID,
	}, nil
}

type ProductOutputGraph struct {
	ID          int
	Name        string
	Description string
	Price       decimal.Decimal
	Quantity    int
	Author      UserOutput
	Category    PCateOutput
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetProductsGraph retrieves all the products and  author, category
func (c *Controller) GetProductsGraph(ctx context.Context, pFilter ProductCtrlFilter) ([]ProductOutputGraph, error) {
	products, err := c.Repository.GetProductsGraph(ctx, repositories.ProductRepoFilter{
		Name: pFilter.Name,
		Date: pFilter.Date,
	})
	if err != nil {
		return nil, err
	}

	var pResp []ProductOutputGraph
	for _, p := range products {
		userResp := UserOutput{
			ID:        p.User.ID,
			Name:      p.User.Name,
			Email:     p.User.Email,
			Role:      p.User.Role,
			CreatedAt: p.User.CreatedAt,
			UpdatedAt: p.User.UpdatedAt,
			Status:    p.User.Status,
		}

		pCateResp := PCateOutput{
			ID:          p.PCate.ID,
			Name:        p.PCate.Name,
			Description: p.PCate.Description,
			CreatedAt:   p.PCate.CreatedAt,
			UpdatedAt:   p.PCate.UpdatedAt,
		}

		pResp = append(pResp, ProductOutputGraph{
			ID:          p.Product.ID,
			Name:        p.Product.Name,
			Description: p.Product.Description,
			Price:       p.Product.Price,
			Quantity:    p.Product.Quantity,
			Category:    pCateResp,
			Author:      userResp,
			CreatedAt:   p.Product.CreatedAt,
			UpdatedAt:   p.Product.UpdatedAt,
		})
	}

	return pResp, nil
}

// ExportProductsToCSV exports all the products in db to a csv file
func (c *Controller) ExportProductsToCSV(ctx context.Context, filter ProductCtrlFilter) (io.Reader, error) {
	pipeReader, pipeWriter := io.Pipe()

	// get all products
	products, err := c.Repository.GetProducts(ctx, repositories.ProductRepoFilter{
		Name: filter.Name,
		Date: filter.Date,
	})
	if err != nil {
		return nil, err
	}

	// write the products to the pipe
	go func() {
		csvWriter := csv.NewWriter(pipeWriter)

		// Write CSV header
		csvWriter.Write([]string{"ID", "Name", "Description", "Price", "Quantity", "AuthorID", "Category", "CreatedAt", "UpdatedAt"})

		for _, p := range products {
			record := []string{
				strconv.Itoa(p.ID),
				p.Name,
				p.Description,
				p.Price.String(),
				strconv.Itoa(p.Quantity),
				strconv.Itoa(p.AuthorID),
				p.CategoryName,
				p.CreatedAt.Format(time.RFC3339),
				p.UpdatedAt.Format(time.RFC3339),
			}
			csvWriter.Write(record)
		}
		csvWriter.Flush()
		pipeWriter.Close()
	}()

	return pipeReader, nil
}

// SendEmailProduct send an email to a list of user gmail with a list of products csv file attachment
func (c *Controller) SendEmailProduct(emailToList []string, reader io.Reader) error {
	sender := email.NewEmailSender()

	// Write email subject and body
	m := email.NewMessage("List of products", "Hi users,\nThis is a list of products exported to a CSV file that is being attached below.\nThanks!")
	m.To = emailToList
	m.Attachments = reader

	if err := sender.Send(m); err != nil {
		return err
	}
	return nil
}
