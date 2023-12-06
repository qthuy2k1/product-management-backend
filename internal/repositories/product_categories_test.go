package repositories

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tanimutomo/sqlfile"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func runSQLTest(dbConn *sql.DB, files ...string) error {
	s := sqlfile.New()
	if err := s.Files(files...); err != nil {
		return err
	}
	if _, err := s.Exec(dbConn); err != nil {
		return err
	}
	return nil
}

func Test_ProductCategoryRepository_CreateProductCategory(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		input ProductCategory
		err   error
	}{
		"success": {
			input: ProductCategory{
				Name:        "Smartphone",
				Description: "A mobile device",
			},
		},
		"create product category with name longer than 255 characters": {
			input: ProductCategory{
				Name:        "cellphone" + strings.Repeat("a", 251),
				Description: "An Apple cellphone with A16 Bionic chip, 6GB RAM and 128GB storage",
			},
			err: errors.New("models: unable to insert into product_categories: pq: value too long for type character varying(255)"),
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, err := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, err)
			boil.SetDB(db)
			defer db.Exec("DELETE FROM product_categories;")

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			err = repo.CreateProductCategory(context.Background(), tc.input)

			// then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}

func Test_UserRepository_GetProductCategory(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		input     int
		expResult models.ProductCategory
		err       error
	}{
		"success": {
			input: 1001,
			expResult: models.ProductCategory{
				ID:          1001,
				Name:        "Smartwatch",
				Description: "A wearable device that functions as a digital watch",
			},
		},
		"product category not found": {
			input: 99,
			err:   ErrProductCategoryNotFound,
		},
	}
	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, dbErr := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, dbErr)
			boil.SetDB(db)

			err := runSQLTest(db, "./datatest/product_categories/insert_product_category.sql")
			assert.NoError(t, err)

			defer runSQLTest(db, "./datatest/product_categories/rollback_insert_product_category.sql")

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			result, err := repo.GetProductCategory(context.Background(), tc.input)

			result.CreatedAt = time.Time{}
			result.UpdatedAt = time.Time{}

			// Then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func Test_UserRepository_GetProductCategoryByName(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		input     string
		expResult models.ProductCategory
		err       error
	}{
		"success": {
			input: "Smartwatch",
			expResult: models.ProductCategory{
				ID:          1001,
				Name:        "Smartwatch",
				Description: "A wearable device that functions as a digital watch",
			},
		},
		"success 2": {
			input: "Monitor",
			expResult: models.ProductCategory{
				ID:          1002,
				Name:        "Monitor",
				Description: "A display screen that allows you to view images, videos, and other visual content from a computer or other device",
			},
		},
		"product category not found": {
			input: "Earphone",
			err:   ErrProductCategoryNotFound,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, dbErr := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, dbErr)
			boil.SetDB(db)

			err := runSQLTest(db, "./datatest/product_categories/insert_product_category.sql")
			assert.NoError(t, err)

			defer runSQLTest(db, "./datatest/product_categories/rollback_insert_product_category.sql")

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			result, err := repo.GetProductCategoryByName(context.Background(), tc.input)

			// Then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)

				result.CreatedAt = time.Time{}
				result.UpdatedAt = time.Time{}

				require.Equal(t, tc.expResult, result)
			}
		})
	}
}
