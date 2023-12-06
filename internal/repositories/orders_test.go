package repositories

import (
	context "context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Test_OrderRepository_GetOrders(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		expResult  []OrderOutputGraph
		totalCount int64
		filter     OrderFilterRepo
		err        error
	}{
		"get orders successfully": {
			expResult: []OrderOutputGraph{
				{
					ID:          1000,
					UserName:    "Quang Thuy",
					UserEmail:   "qthuy1000@gmail.com",
					Status:      "Created",
					TotalPrice:  decimal.NewFromBigInt(big.NewInt(120000), -2),
					ItemID:      "{1000,1001}",
					ProductName: `{"iPhone 14 1000","iPhone 14 1001"}`,
					Quantity:    "{1,1}",
					ItemPrice:   "{1200.00,1200.00}",
				},
				{
					ID:          1001,
					UserName:    "Quang Thuy",
					UserEmail:   "qthuy1000@gmail.com",
					Status:      "Created",
					TotalPrice:  decimal.NewFromBigInt(big.NewInt(120000), -2),
					ItemID:      "{1002,1003}",
					ProductName: `{"iPhone 14 1000","iPhone 14 1001"}`,
					Quantity:    "{1,1}",
					ItemPrice:   "{1200.00,1200.00}",
				},
			},
			totalCount: 2,
			filter:     OrderFilterRepo{},
		},
		"get orders with sort order filter successfully": {
			expResult: []OrderOutputGraph{
				{
					ID:          1001,
					UserName:    "Quang Thuy",
					UserEmail:   "qthuy1000@gmail.com",
					Status:      "Updated",
					TotalPrice:  decimal.NewFromBigInt(big.NewInt(120000), -2),
					ItemID:      "{1002,1003}",
					ProductName: `{"iPhone 14 1000","iPhone 14 1001"}`,
					Quantity:    "{1,1}",
					ItemPrice:   "{1200.00,1200.00}",
				},
				{
					ID:          1000,
					UserName:    "Quang Thuy",
					UserEmail:   "qthuy1000@gmail.com",
					Status:      "Created",
					TotalPrice:  decimal.NewFromBigInt(big.NewInt(120000), -2),
					ItemID:      "{1000,1001}",
					ProductName: `{"iPhone 14 1000","iPhone 14 1001"}`,
					Quantity:    "{1,1}",
					ItemPrice:   "{1200.00,1200.00}",
				},
			},
			totalCount: 2,
			filter:     OrderFilterRepo{},
		},
		"get orders with page size is 1, page number is 1 and sort status order is desc": {
			expResult: []OrderOutputGraph{
				{
					ID:          1001,
					UserName:    "Quang Thuy",
					UserEmail:   "qthuy1000@gmail.com",
					Status:      "Updated",
					TotalPrice:  decimal.NewFromBigInt(big.NewInt(120000), -2),
					ItemID:      "{1002,1003}",
					ProductName: `{"iPhone 14 1000","iPhone 14 1001"}`,
					Quantity:    "{1,1}",
					ItemPrice:   "{1200.00,1200.00}",
				},
			},
			totalCount: 2,
			filter:     OrderFilterRepo{},
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			// Given
			db, dbErr := Initialize(os.Getenv("DB_URL"))
			assert.NoError(t, dbErr)
			boil.SetDB(db)

			redis := RedisInitialize(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS"))

			err := runSQLTest(db, "./datatest/orders/insert_order.sql")
			assert.NoError(t, err)

			defer runSQLTest(db, "./datatest/orders/rollback_insert_order.sql")

			repo := NewRepository(db, redis)

			// When
			result, totalCount, err := repo.GetOrders(context.Background(), tc.filter)

			// Then
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)

				for i := range result {
					result[i].CreatedAt = time.Time{}
				}

				require.Equal(t, tc.totalCount, totalCount)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}
