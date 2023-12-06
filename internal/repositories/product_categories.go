package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type ProductCategory struct {
	ID          int       `redis:"id"`
	Name        string    `redis:"name"`
	Description string    `redis:"description"`
	CreatedAt   time.Time `redis:"created_at"`
	UpdatedAt   time.Time `redis:"updated_at"`
}

// CreateProductCategory adds a new category to the database
func (r *Repository) CreateProductCategory(ctx context.Context, pcResp ProductCategory) error {
	pc := models.ProductCategory{
		Name:        pcResp.Name,
		Description: pcResp.Description,
	}
	if err := pc.Insert(ctx, boil.GetContextDB(), boil.Infer()); err != nil {
		return err
	}
	return nil
}

// GetProductCategory gets a product category from db by product cateogory id
func (r *Repository) GetProductCategory(ctx context.Context, id int) (models.ProductCategory, error) {
	res := r.Redis.HGetAll(ctx, fmt.Sprintf("pCate:%d", id))
	if len(res.Val()) == 0 {
		productCategory, err := models.FindProductCategory(ctx, boil.GetContextDB(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.ProductCategory{}, ErrProductCategoryNotFound
			}
			return models.ProductCategory{}, err
		}

		pCateCache := ProductCategory{
			ID:          productCategory.ID,
			Name:        productCategory.Name,
			Description: productCategory.Description,
			CreatedAt:   productCategory.CreatedAt,
			UpdatedAt:   productCategory.UpdatedAt,
		}

		// set cache
		if errCache := r.Redis.HSet(ctx, fmt.Sprintf("pCate:%d", id), pCateCache); errCache.Err() != nil {
			return models.ProductCategory{}, errCache.Err()
		}

		return *productCategory, nil
	}

	var pCateScan ProductCategory
	if err := res.Scan(&pCateScan); err != nil {
		return models.ProductCategory{}, err
	}

	return models.ProductCategory{
		ID:          pCateScan.ID,
		Name:        pCateScan.Name,
		Description: pCateScan.Description,
		CreatedAt:   pCateScan.CreatedAt,
		UpdatedAt:   pCateScan.UpdatedAt,
	}, nil
}

// GetProductCategoryByName retrieves a product category from db by name
func (r *Repository) GetProductCategoryByName(ctx context.Context, name string) (models.ProductCategory, error) {
	res := r.Redis.HGetAll(ctx, fmt.Sprintf("pCate:%s", name))
	if len(res.Val()) == 0 {
		productCategory, err := models.ProductCategories(qm.Where("name = ?", name)).One(ctx, boil.GetContextDB())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.ProductCategory{}, ErrProductCategoryNotFound
			}
			return models.ProductCategory{}, err
		}

		pCateCache := ProductCategory{
			ID:          productCategory.ID,
			Name:        productCategory.Name,
			Description: productCategory.Description,
			CreatedAt:   productCategory.CreatedAt,
			UpdatedAt:   productCategory.UpdatedAt,
		}

		// set cache
		if errCache := r.Redis.HSet(ctx, fmt.Sprintf("pCate:%s", name), pCateCache); errCache.Err() != nil {
			return models.ProductCategory{}, errCache.Err()
		}

		return *productCategory, nil
	}

	var pCateScan ProductCategory
	if err := res.Scan(&pCateScan); err != nil {
		return models.ProductCategory{}, err
	}

	return models.ProductCategory{
		ID:          pCateScan.ID,
		Name:        pCateScan.Name,
		Description: pCateScan.Description,
		CreatedAt:   pCateScan.CreatedAt,
		UpdatedAt:   pCateScan.UpdatedAt,
	}, nil
}

const (
	pCateIDDefault   = 1
	pCateNameDefault = "Uncategorized"
	pCateDescDefault = "Product is uncategorized"
)

// GetProductCategoryDefault gets the id of default user
func (r *Repository) GetProductCategoryDefault(ctx context.Context) (models.ProductCategory, error) {
	res := r.Redis.HGetAll(ctx, fmt.Sprintf("pCate:%d", pCateIDDefault))
	if len(res.Val()) == 0 {
		pCate, err := models.FindProductCategory(ctx, boil.GetContextDB(), pCateIDDefault)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				pcDefault := models.ProductCategory{
					ID:          pCateIDDefault,
					Name:        pCateNameDefault,
					Description: pCateDescDefault,
				}
				if errQuery := pcDefault.Insert(ctx, boil.GetContextDB(), boil.Infer()); errQuery != nil {
					return models.ProductCategory{}, errQuery
				}

				pCateCache := ProductCategory{
					ID:          pcDefault.ID,
					Name:        pcDefault.Name,
					Description: pcDefault.Description,
					CreatedAt:   pcDefault.CreatedAt,
					UpdatedAt:   pcDefault.UpdatedAt,
				}

				// set cache
				if errCache := r.Redis.HSet(ctx, fmt.Sprintf("pCate:%d", pcDefault.ID), pCateCache); errCache.Err() != nil {
					return models.ProductCategory{}, errCache.Err()
				}

				return pcDefault, nil
			}
			return models.ProductCategory{}, err
		}

		pCateCache := ProductCategory{
			ID:          pCate.ID,
			Name:        pCate.Name,
			Description: pCate.Description,
			CreatedAt:   pCate.CreatedAt,
			UpdatedAt:   pCate.UpdatedAt,
		}

		// set cache
		if errCache := r.Redis.HSet(ctx, fmt.Sprintf("pCate:%d", pCate.ID), pCateCache); errCache.Err() != nil {
			return models.ProductCategory{}, errCache.Err()
		}

		return *pCate, nil
	}

	var pCateScan ProductCategory
	if err := res.Scan(&pCateScan); err != nil {
		return models.ProductCategory{}, err
	}

	return models.ProductCategory{
		ID:          pCateScan.ID,
		Name:        pCateScan.Name,
		Description: pCateScan.Description,
		CreatedAt:   pCateScan.CreatedAt,
		UpdatedAt:   pCateScan.UpdatedAt,
	}, nil
}
