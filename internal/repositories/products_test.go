package repositories

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Test_ProductRepository_CreateProduct(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		input Product
		err   error
	}{
		"success": {
			input: Product{
				Name:        "iPhone 14",
				Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:       decimal.New(1500, 0),
				Quantity:    20,
				CategoryID:  1,
				AuthorID:    1,
			},
		},
		"create product with price greater than 15 digits": {
			input: Product{
				Name:        "iPhone 14",
				Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:       decimal.New(150000000000000, 0),
				Quantity:    20,
				CategoryID:  1,
				AuthorID:    1,
			},
			err: errors.New("models: unable to insert into products: pq: numeric field overflow"),
		},
		"create product with name longer than 255 characters": {
			input: Product{
				Name:        "iPhone 14 Pro MaxxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro Maxxx",
				Description: "An Apple cellphone with A16 Bionic chip, 6GB RAM and 128GB storage",
				Price:       decimal.New(1500, 0),
				Quantity:    20,
				CategoryID:  1,
				AuthorID:    1,
			},
			err: errors.New("models: unable to insert into products: pq: value too long for type character varying(255)"),
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, err := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, err)
			boil.SetDB(db)
			defer db.Exec("DELETE FROM products;")

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			err = repo.CreateProduct(context.Background(), tc.input)

			// Then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}

func Test_ProductRepository_GetProduct(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		input     int
		expResult models.Product
		err       error
	}{
		"success": {
			input: 1001,
			expResult: models.Product{
				ID:          1001,
				Name:        "Macbook Air M1 16GB",
				Description: "A macbook air with Apple M1 chip, 16GB of RAM, 512GB SSD",
				Price:       decimal.New(1500, 0),
				Quantity:    10,
				CategoryID:  2002,
				AuthorID:    1001,
			},
		},
		"product not found": {
			input: 999,
			err:   ErrProductNotFound,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, dbErr := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, dbErr)
			boil.SetDB(db)

			err := runSQLTest(db, "./datatest/products/insert_product.sql")
			assert.NoError(t, err)

			defer runSQLTest(db, "./datatest/products/rollback_insert_product.sql")

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			result, err := repo.GetProduct(context.Background(), tc.input)

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

func Test_ProductRepository_UpdateProduct(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		input models.Product
		err   error
	}{
		"success": {
			input: models.Product{
				Name:        "iPhone 14",
				Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:       decimal.New(1500, 0),
				Quantity:    20,
				CategoryID:  1,
				AuthorID:    1,
			},
		},
		"update product with price greater than 15 digits": {
			input: models.Product{
				Name:        "iPhone 14",
				Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:       decimal.New(1500000000000000, 0),
				Quantity:    20,
				CategoryID:  1,
				AuthorID:    1,
			},
			err: errors.New("models: unable to update products row: pq: numeric field overflow"),
		},
		"update product with name longer than 255 characters": {
			input: models.Product{
				Name:        "iPhone 14 Pro MaxxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro MaxxxiPhone 14 Pro Maxxx",
				Description: "An Apple cellphone with A16 Bionic chip, 6GB RAM and 128GB storage",
				Price:       decimal.New(1500, 0),
				Quantity:    20,
				CategoryID:  1,
				AuthorID:    1,
			},
			err: errors.New("models: unable to update products row: pq: value too long for type character varying(255)"),
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, err := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, err)
			boil.SetDB(db)

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			err = repo.UpdateProduct(context.Background(), nil, tc.input)

			// Then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}

func Test_ProductRepository_DeleteProduct(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		input int
		err   error
	}{
		"success": {
			input: 115,
		},
		"product not found": {
			input: -1,
			err:   ErrProductNotFound,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, dbErr := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, dbErr)
			boil.SetDB(db)

			err := runSQLTest(db, "./datatest/products/insert_product.sql")
			assert.NoError(t, err)

			defer runSQLTest(db, "./datatest/products/rollback_insert_product.sql")

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			err = repo.DeleteProduct(context.Background(), tc.input)

			// Then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}

func Test_ProductRepository_GetProducts(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		expResult []ProductOutput
		filter    ProductRepoFilter
		err       error
	}{
		"get products successfully": {
			expResult: []ProductOutput{
				{
					ID:           1001,
					Name:         "Macbook Air M1 16GB",
					Description:  "A macbook air with Apple M1 chip, 16GB of RAM, 512GB SSD",
					Price:        decimal.New(1500, 0),
					Quantity:     10,
					CategoryName: "Laptop",
					AuthorID:     1001,
				},
				{
					ID:           1002,
					Name:         "iPhone 13",
					Description:  "An iPhone with Apple A15 chipset, 4GB RAM, 128GB storage",
					Price:        decimal.New(1500, 0),
					Quantity:     10,
					CategoryName: "Smartphone",
					AuthorID:     1002,
				},
			},
		},
		"get products with filter successfully": {
			expResult: []ProductOutput{
				{
					ID:           1001,
					Name:         "Macbook Air M1 16GB",
					Description:  "A macbook air with Apple M1 chip, 16GB of RAM, 512GB SSD",
					Price:        decimal.New(1500, 0),
					Quantity:     10,
					CategoryName: "Laptop",
					AuthorID:     1001,
				},
			},
			filter: ProductRepoFilter{
				Name: "Macbook",
			},
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, dbErr := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, dbErr)
			boil.SetDB(db)

			err := runSQLTest(db, "./datatest/products/insert_product.sql")
			assert.NoError(t, err)

			defer runSQLTest(db, "./datatest/products/rollback_insert_product.sql")

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			repo := NewRepository(db, redis)

			// When
			result, err := repo.GetProducts(context.Background(), tc.filter)

			// Then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)

				for i := range result {
					result[i].CreatedAt = time.Time{}
					result[i].UpdatedAt = time.Time{}
				}

				require.Equal(t, tc.expResult, result)
			}
		})
	}
}
