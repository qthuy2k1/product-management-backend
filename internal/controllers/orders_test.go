package controllers

import (
	context "context"
	"database/sql"
	"testing"
	"time"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/qthuy2k1/product-management/internal/repositories"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_OrderController_CreateOrder(t *testing.T) {
	type mockUserRepo struct {
		output models.User
		err    error
	}

	type mockCreateOrderRepo struct {
		input  repositories.Order
		output models.Order
		err    error
	}

	type mockUpdateOrderRepo struct {
		input models.Order
		err   error
	}

	type mockGetProductRepo struct {
		output models.Product
		err    error
	}

	type mockUpdateProductRepo struct {
		input models.Product
		err   error
	}

	type mockOrderItemRepo struct {
		orderItemInputList []repositories.OrderItem
		err                error
	}

	testCases := map[string]struct {
		expCall               bool
		orderInput            OrderInput
		orderItemInput        []OrderItemInput
		mockUserRepo          mockUserRepo
		mockCreateOrderRepo   mockCreateOrderRepo
		mockUpdateOrderRepo   mockUpdateOrderRepo
		mockGetProductRepo    []mockGetProductRepo
		mockUpdateProductRepo []mockUpdateProductRepo
		mockOrderItemRepo     mockOrderItemRepo
		expErr                error
	}{
		"create order successfully": {
			expCall: true,
			orderInput: OrderInput{
				UserID: 1,
				Status: "Created",
			},
			orderItemInput: []OrderItemInput{
				{
					ProductID: 1,
					Quantity:  1,
				},
				{
					ProductID: 2,
					Quantity:  1,
				},
			},
			mockUserRepo: mockUserRepo{
				output: models.User{
					ID:    1,
					Name:  "Thuy Nguyen",
					Email: "qthuy@gmail.com",
				},
			},
			mockCreateOrderRepo: mockCreateOrderRepo{
				input: repositories.Order{
					UserID: 1,
					Status: "Created",
				},
				output: models.Order{
					ID:     1,
					UserID: 1,
					Status: "Created",
				},
			},
			mockGetProductRepo: []mockGetProductRepo{
				{
					output: models.Product{
						ID:          1,
						Name:        "iPhone 14",
						Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
						Price:       decimal.New(1500, 0),
						Quantity:    20,
						CategoryID:  1,
						AuthorID:    1,
					},
				},
				{
					output: models.Product{
						ID:          2,
						Name:        "iPhone 13",
						Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
						Price:       decimal.New(1500, 0),
						Quantity:    20,
						CategoryID:  1,
						AuthorID:    1,
					},
				},
			},
			mockUpdateProductRepo: []mockUpdateProductRepo{
				{
					input: models.Product{
						ID:          1,
						Name:        "iPhone 14",
						Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
						Price:       decimal.New(1500, 0),
						Quantity:    19,
						CategoryID:  1,
						AuthorID:    1,
					},
				},
				{
					input: models.Product{
						ID:          2,
						Name:        "iPhone 13",
						Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
						Price:       decimal.New(1500, 0),
						Quantity:    19,
						CategoryID:  1,
						AuthorID:    1,
					},
				},
			},
			mockOrderItemRepo: mockOrderItemRepo{
				orderItemInputList: []repositories.OrderItem{
					{
						ProductID: 1,
						Quantity:  1,
						Price:     decimal.New(1500, 0),
					},
					{
						ProductID: 2,
						Quantity:  1,
						Price:     decimal.New(1500, 0),
					},
				},
			},
			mockUpdateOrderRepo: mockUpdateOrderRepo{
				input: models.Order{
					ID:     1,
					UserID: 1,
					Status: "Created",
					TotalPrice: decimal.NullDecimal{
						Decimal: decimal.New(3000, 0),
						Valid:   true,
					},
				},
			},
		},
		"user not found": {
			expCall: true,
			orderInput: OrderInput{
				UserID: 1,
				Status: "Created",
			},
			mockUserRepo: mockUserRepo{
				err: ErrUserNotFound,
			},
			expErr: ErrUserNotFound,
		},
		"product not found": {
			expCall: true,
			orderInput: OrderInput{
				UserID: 1,
				Status: "Created",
			},
			orderItemInput: []OrderItemInput{
				{
					ProductID: 1,
					Quantity:  1,
				},
			},
			mockUserRepo: mockUserRepo{
				output: models.User{
					ID:    1,
					Name:  "Thuy Nguyen",
					Email: "qthuy@gmail.com",
				},
			},
			mockCreateOrderRepo: mockCreateOrderRepo{
				input: repositories.Order{
					UserID: 1,
					Status: "Created",
				},
				output: models.Order{
					ID:     1,
					UserID: 1,
					Status: "Created",
				},
			},
			mockGetProductRepo: []mockGetProductRepo{
				{
					err: ErrProductNotFound,
				},
			},
			mockUpdateProductRepo: []mockUpdateProductRepo{
				{},
			},
			expErr: ErrProductNotFound,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockRepo := &repositories.MockIRepository{}
			controller := NewController(mockRepo)
			tx := sql.Tx{}

			if tc.expCall {
				mockRepo.On("GetUser", context.Background(), tc.orderInput.UserID).Return(tc.mockUserRepo.output, tc.mockUserRepo.err)
				mockRepo.On("BeginTx", context.Background()).Return(&tx, nil)
				mockRepo.On("RollbackTx", &tx).Return(nil)
				mockRepo.On("CommitTx", &tx).Return(nil)
				mockRepo.On("CreateOrder", context.Background(), tc.mockCreateOrderRepo.input).Return(tc.mockCreateOrderRepo.output, tc.mockCreateOrderRepo.err)
				for i, oi := range tc.orderItemInput {
					mockRepo.On("GetProduct", context.Background(), oi.ProductID).Return(tc.mockGetProductRepo[i].output, tc.mockGetProductRepo[i].err)
					mockRepo.On("UpdateProduct", context.Background(), tc.mockUpdateProductRepo[i].input).Return(tc.mockUpdateProductRepo[i].err)
				}
				mockRepo.On("CreateOrderItem", context.Background(), tc.mockOrderItemRepo.orderItemInputList, tc.mockUpdateOrderRepo.input).Return(tc.mockOrderItemRepo.err)
				mockRepo.On("UpdateOrder", context.Background(), tc.mockUpdateOrderRepo.input).Return(tc.mockUpdateOrderRepo.err)
			}

			if err := controller.CreateOrder(context.Background(), tc.orderInput, tc.orderItemInput); tc.expErr != nil {
				assert.EqualError(t, err, tc.expErr.Error())
			}
		})
	}
}

func Test_OrderController_GetOrders(t *testing.T) {
	myCreatedTime, err := time.Parse("2006-01-02 15:04:05", "2023-06-02 00:00:00")
	assert.NoError(t, err)
	type mockOrderRepo struct {
		filter repositories.OrderFilterRepo
		output []repositories.OrderOutputGraph
		err    error
	}

	testCases := map[string]struct {
		expCall       bool
		filter        OrderFilterCtrl
		mockOrderRepo mockOrderRepo
		orderOutput   []OrderOutputGraph
		expErr        error
	}{
		"get all orders successfully": {
			expCall: true,
			orderOutput: []OrderOutputGraph{
				{
					ID:         1,
					UserName:   "Quang Thuy",
					UserEmail:  "qthuy@gmail.com",
					TotalPrice: decimal.New(2400, 0),
					Status:     "Created",
					Items: []OrderItemOutput{
						{
							ID:          1,
							ProductName: "iPhone 14",
							Price:       decimal.NewFromFloat(1200),
							Quantity:    1,
						},
						{
							ID:          2,
							ProductName: "Macbook Air M1",
							Price:       decimal.NewFromFloat(1200),
							Quantity:    1,
						},
					},
					CreatedAt: myCreatedTime,
				},
				{
					ID:         2,
					UserName:   "Quang Thuy",
					UserEmail:  "qthuy@gmail.com",
					TotalPrice: decimal.New(2400, 0),
					Status:     "Created",
					Items: []OrderItemOutput{
						{
							ID:          1,
							ProductName: "iPhone 14",
							Price:       decimal.NewFromFloat(1200),
							Quantity:    1,
						},
						{
							ID:          2,
							ProductName: "Macbook Air M1",
							Price:       decimal.NewFromFloat(1200),
							Quantity:    1,
						},
					},
					CreatedAt: myCreatedTime,
				},
			},
			mockOrderRepo: mockOrderRepo{
				filter: repositories.OrderFilterRepo{},
				output: []repositories.OrderOutputGraph{
					{
						ID:          1,
						UserName:    "Quang Thuy",
						UserEmail:   "qthuy@gmail.com",
						TotalPrice:  decimal.New(2400, 0),
						Status:      "Created",
						ItemID:      "{1,2}",
						ProductName: "{iPhone 14,Macbook Air M1}",
						Quantity:    "{1,1}",
						ItemPrice:   "{1200,1200}",
						CreatedAt:   myCreatedTime,
					},
					{
						ID:          2,
						UserName:    "Quang Thuy",
						UserEmail:   "qthuy@gmail.com",
						TotalPrice:  decimal.New(2400, 0),
						Status:      "Created",
						ItemID:      "{1,2}",
						ProductName: "{iPhone 14,Macbook Air M1}",
						Quantity:    "{1,1}",
						ItemPrice:   "{1200,1200}",
						CreatedAt:   myCreatedTime,
					},
				},
			},
		},
		"get all orders successfully with filter": {
			expCall: true,
			filter:  OrderFilterCtrl{},
			orderOutput: []OrderOutputGraph{
				{
					ID:         2,
					UserName:   "Quang Thuy",
					UserEmail:  "qthuy@gmail.com",
					TotalPrice: decimal.New(2400, 0),
					Status:     "Updated",
					Items: []OrderItemOutput{
						{
							ID:          1,
							ProductName: "iPhone 14",
							Price:       decimal.NewFromFloat(1200),
							Quantity:    1,
						},
						{
							ID:          2,
							ProductName: "Macbook Air M1",
							Price:       decimal.NewFromFloat(1200),
							Quantity:    1,
						},
					},
					CreatedAt: myCreatedTime,
				},
				{
					ID:         1,
					UserName:   "Quang Thuy",
					UserEmail:  "qthuy@gmail.com",
					TotalPrice: decimal.New(2400, 0),
					Status:     "Created",
					Items: []OrderItemOutput{
						{
							ID:          1,
							ProductName: "iPhone 14",
							Price:       decimal.NewFromFloat(1200),
							Quantity:    1,
						},
						{
							ID:          2,
							ProductName: "Macbook Air M1",
							Price:       decimal.NewFromFloat(1200),
							Quantity:    1,
						},
					},
					CreatedAt: myCreatedTime,
				},
			},
			mockOrderRepo: mockOrderRepo{
				filter: repositories.OrderFilterRepo{},
				output: []repositories.OrderOutputGraph{
					{
						ID:          2,
						UserName:    "Quang Thuy",
						UserEmail:   "qthuy@gmail.com",
						TotalPrice:  decimal.New(2400, 0),
						Status:      "Updated",
						ItemID:      "{1,2}",
						ProductName: "{iPhone 14,Macbook Air M1}",
						Quantity:    "{1,1}",
						ItemPrice:   "{1200,1200}",
						CreatedAt:   myCreatedTime,
					},
					{
						ID:          1,
						UserName:    "Quang Thuy",
						UserEmail:   "qthuy@gmail.com",
						TotalPrice:  decimal.New(2400, 0),
						Status:      "Created",
						ItemID:      "{1,2}",
						ProductName: "{iPhone 14,Macbook Air M1}",
						Quantity:    "{1,1}",
						ItemPrice:   "{1200,1200}",
						CreatedAt:   myCreatedTime,
					},
				},
			},
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.MockIRepository{}
			controller := NewController(&mockRepo)

			if tc.expCall {
				mockRepo.On("GetOrders", context.Background(), tc.mockOrderRepo.filter).Return(tc.mockOrderRepo.output, tc.mockOrderRepo.err)
			}

			orders, _, err := controller.GetOrders(context.Background(), tc.filter)
			if tc.expErr != nil {
				assert.EqualError(t, err, tc.expErr.Error())
			} else {
				assert.Equal(t, tc.orderOutput, orders)
			}
		})
	}
}
