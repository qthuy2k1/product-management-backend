package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.32

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/qthuy2k1/product-management/internal/handlers/graph/model"
	"github.com/shopspring/decimal"
)

// CreateProduct is the resolver for the createProduct field.
func (r *mutationResolver) CreateProduct(ctx context.Context, input model.ProductRequest) (bool, error) {
	product, errResp := validateAndConvertProduct(input)
	if errResp != nil {
		return false, errResp
	}

	if err := r.Controller.CreateProduct(ctx, product); err != nil {
		log.Println(err)
		return false, convertCtrlError(err)
	}

	return true, nil
}

// GetProducts is the resolver for the getProducts field.
func (r *queryResolver) GetProducts(ctx context.Context, queryName string, date string) ([]*model.Product, error) {
	pCtrlFilter, err := validateAndConvertProductFilter(queryName, date)
	if err != nil {
		return nil, err
	}

	products, err := r.Controller.GetProductsGraph(ctx, pCtrlFilter)
	if err != nil {
		log.Println(err)
		return nil, convertCtrlError(err)
	}

	var productsResp []*model.Product
	for _, p := range products {
		productsResp = append(productsResp, &model.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price.Abs().InexactFloat64(),
			Quantity:    p.Quantity,
			Author: &model.User{
				ID:        p.Author.ID,
				Name:      p.Author.Name,
				Email:     p.Author.Email,
				Role:      p.Author.Role,
				Status:    p.Author.Status,
				CreatedAt: p.Author.CreatedAt.String(),
				UpdatedAt: p.Author.UpdatedAt.String(),
			},
			Category: &model.ProductCategory{
				ID:          p.Category.ID,
				Name:        p.Category.Name,
				Description: p.Category.Description,
				CreatedAt:   p.Category.CreatedAt.String(),
				UpdatedAt:   p.Category.UpdatedAt.String(),
			},
			CreatedAt: p.CreatedAt.String(),
			UpdatedAt: p.UpdatedAt.String(),
		})
	}

	return productsResp, nil
}

// validateAndConvertProduct validates the product from body request and return product struct in controller layer
func validateAndConvertProduct(pReq model.ProductRequest) (controllers.ProductInput, error) {
	if len(strings.TrimSpace(pReq.Name)) == 0 {
		return controllers.ProductInput{}, ErrMissingName
	}

	if len(strings.TrimSpace(pReq.Name)) > 255 {
		return controllers.ProductInput{}, ErrNameTooLong
	}

	if len(strings.TrimSpace(pReq.Description)) == 0 {
		return controllers.ProductInput{}, ErrMissingDesc
	}

	if pReq.Price <= 0 && pReq.Price > 15 {
		return controllers.ProductInput{}, ErrInvalidPrice
	}

	if pReq.Quantity < 0 {
		return controllers.ProductInput{}, ErrInvalidQuantity
	}

	if pReq.AuthorID <= 0 {
		return controllers.ProductInput{}, ErrInvalidAuthorID
	}

	if len(strings.TrimSpace(pReq.CategoryName)) == 0 {
		return controllers.ProductInput{}, ErrMissingCategoryName
	}

	return controllers.ProductInput{
		Name:         strings.TrimSpace(pReq.Name),
		Description:  strings.TrimSpace(pReq.Description),
		Price:        decimal.NewFromFloat(pReq.Price),
		Quantity:     pReq.Quantity,
		AuthorID:     pReq.AuthorID,
		CategoryName: strings.TrimSpace(pReq.CategoryName),
	}, nil
}

func validateAndConvertProductFilter(queryName, date string) (controllers.ProductCtrlFilter, error) {
	var pCtrlFilter controllers.ProductCtrlFilter
	if len(strings.TrimSpace(queryName)) != 0 {
		pCtrlFilter.Name = queryName
	}
	if len(strings.TrimSpace(date)) != 0 {
		if _, err := time.Parse("2006-01-02", date); err != nil {
			return controllers.ProductCtrlFilter{}, ErrDateBadRequest
		}
		pCtrlFilter.Date = date
	}

	return pCtrlFilter, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
