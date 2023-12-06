package controllers

import (
	"context"
	"errors"
	"time"

	"github.com/qthuy2k1/product-management/internal/repositories"
)

type PCateInput struct {
	Name        string
	Description string
}

// CreateProductCategory adds a new category to the database
func (c *Controller) CreateProductCategory(ctx context.Context, pCateInput PCateInput) error {
	pCateRepo := repositories.ProductCategory{
		Name:        pCateInput.Name,
		Description: pCateInput.Description,
	}

	return c.Repository.CreateProductCategory(ctx, pCateRepo)
}

type PCateOutput struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetProductCategoryByName retrieves a product category by name
func (c *Controller) GetProductCategoryByName(ctx context.Context, name string) (PCateOutput, error) {
	pCate, err := c.Repository.GetProductCategoryByName(ctx, name)
	if err != nil {
		if errors.Is(err, repositories.ErrProductCategoryNotFound) {
			return PCateOutput{}, ErrProductCategoryNotFound
		}
		return PCateOutput{}, err
	}
	pCateResponse := PCateOutput{
		ID:          pCate.ID,
		Name:        pCate.Name,
		Description: pCate.Description,
		CreatedAt:   pCate.CreatedAt,
		UpdatedAt:   pCate.UpdatedAt,
	}
	return pCateResponse, nil
}
