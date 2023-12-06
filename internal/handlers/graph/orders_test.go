package graph

import (
	"context"
	"testing"
	"time"

	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/qthuy2k1/product-management/internal/handlers/graph/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_OrderHandler_CreateOrder(t *testing.T) {
	type mockOrderCtrl struct {
		expCall        bool
		orderInput     controllers.OrderInput
		orderItemInput []controllers.OrderItemInput
		err            error
	}

	testCases := map[string]struct {
		givenOrderInput model.OrderRequest // payload provided from end-user
		mockOrderCtrl   mockOrderCtrl
		expResp         bool
		expErr          error
	}{
		"create order successfully": {
			givenOrderInput: model.OrderRequest{
				UserID: 1,
				Status: "Created",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 1,
						Quantity:  1,
					},
					{
						ProductID: 2,
						Quantity:  1,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: true,
				orderInput: controllers.OrderInput{
					UserID: 1,
					Status: "Created",
				},
				orderItemInput: []controllers.OrderItemInput{
					{
						ProductID: 1,
						Quantity:  1,
					},
					{
						ProductID: 2,
						Quantity:  1,
					},
				},
			},
			expResp: true,
		},
		"invalid product id": {
			givenOrderInput: model.OrderRequest{
				UserID: 1,
				Status: "Created",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 0,
						Quantity:  1,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: false,
			},
			expResp: false,
			expErr:  ErrInvalidProductID,
		},
		"invalid quantity": {
			givenOrderInput: model.OrderRequest{
				UserID: 1,
				Status: "Created",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 1,
						Quantity:  0,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: false,
			},
			expResp: false,
			expErr:  ErrInvalidQuantity,
		},
		"invalid user id": {
			givenOrderInput: model.OrderRequest{
				UserID: 0,
				Status: "Created",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 1,
						Quantity:  1,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: false,
			},
			expResp: false,
			expErr:  ErrInvalidUserID,
		},
		"missing order status": {
			givenOrderInput: model.OrderRequest{
				UserID: 1,
				Status: "  ",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 1,
						Quantity:  1,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: false,
			},
			expResp: false,
			expErr:  ErrMissingOrderStatus,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			resolver := Resolver{Controller: mockController}

			if tc.mockOrderCtrl.expCall {
				mockController.On("CreateOrder", context.Background(), tc.mockOrderCtrl.orderInput, tc.mockOrderCtrl.orderItemInput).Return(tc.mockOrderCtrl.err)
			}

			result, err := resolver.Mutation().CreateOrder(context.Background(), tc.givenOrderInput)

			if err != nil {
				assert.EqualError(t, err, tc.expErr.Error())
			} else {
				assert.Equal(t, tc.expResp, result)
			}

			if tc.mockOrderCtrl.expCall {
				mockController.AssertCalled(t, "CreateOrder", context.Background(), tc.mockOrderCtrl.orderInput, tc.mockOrderCtrl.orderItemInput)
			}
		})
	}
}

func Test_OrderHandler_UpdateOrder(t *testing.T) {
	type mockOrderCtrl struct {
		expCall    bool
		orderID    int
		orderInput controllers.OrderInput
		err        error
	}

	testCases := map[string]struct {
		givenOrderIDInput int                // payload provided from end-user
		givenOrderInput   model.OrderRequest // payload provided from end-user
		mockOrderCtrl     mockOrderCtrl
		expResp           bool
		expErr            error
	}{
		"update order successfully": {
			givenOrderIDInput: 1,
			givenOrderInput: model.OrderRequest{
				UserID: 1,
				Status: "Updated",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 1,
						Quantity:  1,
					},
					{
						ProductID: 2,
						Quantity:  1,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: true,
				orderID: 1,
				orderInput: controllers.OrderInput{
					UserID: 1,
					Status: "Updated",
					OrderItem: []controllers.OrderItemInput{
						{
							ProductID: 1,
							Quantity:  1,
						},
						{
							ProductID: 2,
							Quantity:  1,
						},
					},
				},
			},
			expResp: true,
		},
		"invalid order id": {
			givenOrderInput: model.OrderRequest{
				UserID: 1,
				Status: "Updated",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 0,
						Quantity:  1,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: false,
			},
			expResp: false,
			expErr:  ErrInvalidOrderID,
		},
		"invalid product id": {
			givenOrderIDInput: 1,
			givenOrderInput: model.OrderRequest{
				UserID: 1,
				Status: "Updated",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 0,
						Quantity:  1,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: false,
			},
			expResp: false,
			expErr:  ErrInvalidProductID,
		},
		"invalid quantity": {
			givenOrderIDInput: 1,
			givenOrderInput: model.OrderRequest{
				UserID: 1,
				Status: "Created",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 1,
						Quantity:  0,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: false,
			},
			expResp: false,
			expErr:  ErrInvalidQuantity,
		},
		"invalid user id": {
			givenOrderIDInput: 1,
			givenOrderInput: model.OrderRequest{
				UserID: 0,
				Status: "Created",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 1,
						Quantity:  1,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: false,
			},
			expResp: false,
			expErr:  ErrInvalidUserID,
		},
		"missing order status": {
			givenOrderIDInput: 1,
			givenOrderInput: model.OrderRequest{
				UserID: 1,
				Status: "  ",
				Items: []*model.OrderItemRequest{
					{
						ProductID: 1,
						Quantity:  1,
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: false,
			},
			expResp: false,
			expErr:  ErrMissingOrderStatus,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			resolver := Resolver{Controller: mockController}

			if tc.mockOrderCtrl.expCall {
				mockController.On("UpdateOrder", context.Background(), tc.mockOrderCtrl.orderID, tc.mockOrderCtrl.orderInput).Return(tc.mockOrderCtrl.err)
			}

			result, err := resolver.Mutation().UpdateOrder(context.Background(), tc.givenOrderIDInput, tc.givenOrderInput)

			if err != nil {
				assert.EqualError(t, err, tc.expErr.Error())
			} else {
				assert.Equal(t, tc.expResp, result)
			}

			if tc.mockOrderCtrl.expCall {
				mockController.AssertCalled(t, "UpdateOrder", context.Background(), tc.mockOrderCtrl.orderID, tc.mockOrderCtrl.orderInput)
			}
		})
	}
}

func Test_OrderHandler_GetOrders(t *testing.T) {
	myCreatedTime, err := time.Parse("2006-01-02 15:04:05", "2023-06-02 00:00:00")
	assert.NoError(t, err)
	myCreatedTimeStr := "02-06-2023 00:00:00"
	type mockOrderCtrl struct {
		expCall bool
		filter  controllers.OrderFilterCtrl
		output  []controllers.OrderOutputGraph
		err     error
	}

	testCases := map[string]struct {
		dateStart      string
		dateEnd        string
		sortDate       bool
		sortDateDesc   bool
		sortStatus     bool
		sortStatusDesc bool
		pageSize       int
		pageNumber     int
		mockOrderCtrl  mockOrderCtrl
		expResp        []*model.Order
		expErr         error
	}{
		"get all orders successfully": {
			expResp: []*model.Order{
				{
					ID:     1,
					Status: "Updated",
					User: &model.User{
						Name:  "Quang Thuy",
						Email: "qthuy@gmail.com",
					},
					CreatedAt: &myCreatedTimeStr,
					Total:     &[]float64{2400}[0],
					Items: []*model.OrderItem{
						{
							ID: 1,
							Product: &model.Product{
								Name: "iPhone 14",
							},
							Quantity: 1,
							Price:    1200,
						},
						{
							ID: 2,
							Product: &model.Product{
								Name: "Macbook Air M1",
							},
							Quantity: 1,
							Price:    1200,
						},
					},
				},
				{
					ID:     2,
					Status: "Updated",
					User: &model.User{
						Name:  "Quang Thuy",
						Email: "qthuy@gmail.com",
					},
					CreatedAt: &myCreatedTimeStr,
					Total:     &[]float64{2400}[0],
					Items: []*model.OrderItem{
						{
							ID: 1,
							Product: &model.Product{
								Name: "iPhone 14",
							},
							Quantity: 1,
							Price:    1200,
						},
						{
							ID: 2,
							Product: &model.Product{
								Name: "Macbook Air M1",
							},
							Quantity: 1,
							Price:    1200,
						},
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: true,
				output: []controllers.OrderOutputGraph{
					{
						ID:         1,
						UserName:   "Quang Thuy",
						UserEmail:  "qthuy@gmail.com",
						Status:     "Updated",
						TotalPrice: decimal.New(2400, 0),
						CreatedAt:  myCreatedTime,
						Items: []controllers.OrderItemOutput{
							{
								ID:          1,
								ProductName: "iPhone 14",
								Price:       decimal.New(1200, 0),
								Quantity:    1,
							},
							{
								ID:          2,
								ProductName: "Macbook Air M1",
								Price:       decimal.New(1200, 0),
								Quantity:    1,
							},
						},
					},
					{
						ID:         2,
						UserName:   "Quang Thuy",
						UserEmail:  "qthuy@gmail.com",
						Status:     "Updated",
						TotalPrice: decimal.New(2400, 0),
						CreatedAt:  myCreatedTime,
						Items: []controllers.OrderItemOutput{
							{
								ID:          1,
								ProductName: "iPhone 14",
								Price:       decimal.New(1200, 0),
								Quantity:    1,
							},
							{
								ID:          2,
								ProductName: "Macbook Air M1",
								Price:       decimal.New(1200, 0),
								Quantity:    1,
							},
						},
					},
				},
			},
		},
		"get all orders successfully with filter": {
			sortStatus:     true,
			sortStatusDesc: true,
			expResp: []*model.Order{
				{
					ID:     2,
					Status: "Updated",
					User: &model.User{
						Name:  "Quang Thuy",
						Email: "qthuy@gmail.com",
					},
					CreatedAt: &myCreatedTimeStr,
					Total:     &[]float64{2400}[0],
					Items: []*model.OrderItem{
						{
							ID: 1,
							Product: &model.Product{
								Name: "iPhone 14",
							},
							Quantity: 1,
							Price:    1200,
						},
						{
							ID: 2,
							Product: &model.Product{
								Name: "Macbook Air M1",
							},
							Quantity: 1,
							Price:    1200,
						},
					},
				},
				{
					ID:     1,
					Status: "Created",
					User: &model.User{
						Name:  "Quang Thuy",
						Email: "qthuy@gmail.com",
					},
					CreatedAt: &myCreatedTimeStr,
					Total:     &[]float64{2400}[0],
					Items: []*model.OrderItem{
						{
							ID: 1,
							Product: &model.Product{
								Name: "iPhone 14",
							},
							Quantity: 1,
							Price:    1200,
						},
						{
							ID: 2,
							Product: &model.Product{
								Name: "Macbook Air M1",
							},
							Quantity: 1,
							Price:    1200,
						},
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				expCall: true,
				filter:  controllers.OrderFilterCtrl{},
				output: []controllers.OrderOutputGraph{
					{
						ID:         2,
						UserName:   "Quang Thuy",
						UserEmail:  "qthuy@gmail.com",
						Status:     "Updated",
						TotalPrice: decimal.New(2400, 0),
						CreatedAt:  myCreatedTime,
						Items: []controllers.OrderItemOutput{
							{
								ID:          1,
								ProductName: "iPhone 14",
								Price:       decimal.New(1200, 0),
								Quantity:    1,
							},
							{
								ID:          2,
								ProductName: "Macbook Air M1",
								Price:       decimal.New(1200, 0),
								Quantity:    1,
							},
						},
					},
					{
						ID:         1,
						UserName:   "Quang Thuy",
						UserEmail:  "qthuy@gmail.com",
						Status:     "Created",
						TotalPrice: decimal.New(2400, 0),
						CreatedAt:  myCreatedTime,
						Items: []controllers.OrderItemOutput{
							{
								ID:          1,
								ProductName: "iPhone 14",
								Price:       decimal.New(1200, 0),
								Quantity:    1,
							},
							{
								ID:          2,
								ProductName: "Macbook Air M1",
								Price:       decimal.New(1200, 0),
								Quantity:    1,
							},
						},
					},
				},
			},
		},
		"filter date bad request": {
			dateStart: "2023-07-19", // correct way: 19-07-2023
			expResp: []*model.Order{
				{
					ID:     2,
					Status: "Updated",
					User: &model.User{
						Name:  "Quang Thuy",
						Email: "qthuy@gmail.com",
					},
					CreatedAt: &myCreatedTimeStr,
					Total:     &[]float64{2400}[0],
					Items: []*model.OrderItem{
						{
							ID: 1,
							Product: &model.Product{
								Name: "iPhone 14",
							},
							Quantity: 1,
							Price:    1200,
						},
						{
							ID: 2,
							Product: &model.Product{
								Name: "Macbook Air M1",
							},
							Quantity: 1,
							Price:    1200,
						},
					},
				},
				{
					ID:     1,
					Status: "Updated",
					User: &model.User{
						Name:  "Quang Thuy",
						Email: "qthuy@gmail.com",
					},
					CreatedAt: &myCreatedTimeStr,
					Total:     &[]float64{2400}[0],
					Items: []*model.OrderItem{
						{
							ID: 1,
							Product: &model.Product{
								Name: "iPhone 14",
							},
							Quantity: 1,
							Price:    1200,
						},
						{
							ID: 2,
							Product: &model.Product{
								Name: "Macbook Air M1",
							},
							Quantity: 1,
							Price:    1200,
						},
					},
				},
			},
			expErr: ErrDateBadRequest,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			resolver := Resolver{Controller: mockController}

			if tc.mockOrderCtrl.expCall {
				mockController.On("GetOrders", context.Background(), tc.mockOrderCtrl.filter).Return(tc.mockOrderCtrl.output, tc.mockOrderCtrl.err)
			}

			result, err := resolver.Query().GetOrders(context.Background(), nil, nil, model.PaginationInput{})

			if err != nil {
				assert.EqualError(t, err, tc.expErr.Error())
			} else {
				assert.Equal(t, tc.expResp, result)
			}

			if tc.mockOrderCtrl.expCall {
				mockController.On("GetOrders", context.Background(), tc.mockOrderCtrl.filter).Return(tc.mockOrderCtrl.output, tc.mockOrderCtrl.err)
			}
		})
	}
}
