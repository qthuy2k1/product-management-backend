package graph

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/qthuy2k1/product-management/internal/handlers/graph/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_ProductHandler_CreateProduct(t *testing.T) {
	type mockProductCtrl struct {
		expCall      bool
		productInput controllers.ProductInput
		err          error
	}

	testCases := map[string]struct {
		givenInput      model.ProductRequest // payload provided from end-user
		mockProductCtrl mockProductCtrl
		expResp         bool
		expErr          error
	}{
		"create product successfully": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				productInput: controllers.ProductInput{
					Name:         "iPhone 14",
					Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:        decimal.NewFromBigInt(big.NewInt(15), 2),
					Quantity:     20,
					AuthorID:     1,
					CategoryName: "Smartphone",
				},
			},
			givenInput: model.ProductRequest{
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        1500,
				Quantity:     20,
				CategoryName: "Smartphone",
				AuthorID:     1,
			},
			expResp: true,
		},
		"user not found": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				productInput: controllers.ProductInput{
					Name:         "iPhone 14",
					Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:        decimal.NewFromBigInt(big.NewInt(15), 2),
					Quantity:     20,
					AuthorID:     100,
					CategoryName: "Smartphone",
				},
				err: controllers.ErrUserNotFound,
			},
			givenInput: model.ProductRequest{
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        1500,
				Quantity:     20,
				AuthorID:     100,
				CategoryName: "Smartphone",
			},
			expErr: errors.New("user not found"),
		},
		"product category not found": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				productInput: controllers.ProductInput{
					Name:         "iPhone 14",
					Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:        decimal.NewFromBigInt(big.NewInt(15), 2),
					Quantity:     20,
					AuthorID:     40,
					CategoryName: "smartwatchhhh",
				},
				err: controllers.ErrProductCategoryNotFound,
			},
			givenInput: model.ProductRequest{
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        1500,
				Quantity:     20,
				AuthorID:     40,
				CategoryName: "smartwatchhhh",
			},
			expErr: errors.New("product category not found"),
		},
		"create product with missing name field": {
			givenInput: model.ProductRequest{
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        1500,
				Quantity:     20,
				AuthorID:     40,
				CategoryName: "Smartwatch",
			},
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expErr: errors.New("name cannot be blank"),
		},
		"create product with name too long": {
			// 260 'a' characters
			givenInput: model.ProductRequest{
				Name:         strings.Repeat("a", 260),
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        1500,
				Quantity:     20,
				AuthorID:     40,
				CategoryName: "Smartwatch",
			},
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expErr: errors.New("name too long"),
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			resolver := Resolver{Controller: mockController}

			if tc.mockProductCtrl.expCall {
				mockController.On("CreateProduct", context.Background(), tc.mockProductCtrl.productInput).Return(tc.mockProductCtrl.err)
			}

			result, err := resolver.Mutation().CreateProduct(context.Background(), tc.givenInput)

			if err != nil {
				assert.EqualError(t, err, tc.expErr.Error())
			} else {
				assert.Equal(t, tc.expResp, result)
			}

			if tc.mockProductCtrl.expCall {
				mockController.AssertCalled(t, "CreateProduct", context.Background(), tc.mockProductCtrl.productInput)
			}
		})
	}
}

func Test_ProductHandler_GetProducts(t *testing.T) {
	myCreatedTime, err := time.Parse("2006-01-02 15:04:05.999999", "2023-06-02 09:08:36.046843")
	assert.NoError(t, err)
	myUpdatedTime, err := time.Parse("2006-01-02 15:04:05.999999", "2023-06-02 09:08:36.046843")
	assert.NoError(t, err)

	type mockProductCtrl struct {
		query  controllers.ProductCtrlFilter
		output []controllers.ProductOutputGraph
		err    error
	}

	testCases := map[string]struct {
		givenFilterName string // payload provided from end-user
		givenFilterDate string // payload provided from end-user
		mockProductCtrl mockProductCtrl
		expErr          string
		expOutput       []*model.Product
		expCall         bool
	}{
		"get all products successfully": {
			expCall: true,
			mockProductCtrl: mockProductCtrl{
				output: []controllers.ProductOutputGraph{
					{
						ID:          194,
						Name:        "iPhone 14",
						Description: "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:       decimal.New(1500, 0),
						Quantity:    50,
						Author: controllers.UserOutput{
							ID:        2,
							Name:      "Thuy Nguyen",
							Email:     "qthuy@gmail.com",
							CreatedAt: myCreatedTime,
							UpdatedAt: myUpdatedTime,
						},
						Category: controllers.PCateOutput{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A wearable computer",
							CreatedAt:   myCreatedTime,
							UpdatedAt:   myUpdatedTime,
						},
						CreatedAt: myCreatedTime,
						UpdatedAt: myUpdatedTime,
					},
					{
						ID:          195,
						Name:        "Macbook",
						Description: "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:       decimal.New(1500, 0),
						Quantity:    50,
						Author: controllers.UserOutput{
							ID:        2,
							Name:      "Thuy Nguyen",
							Email:     "qthuy@gmail.com",
							CreatedAt: myCreatedTime,
							UpdatedAt: myUpdatedTime,
						},
						Category: controllers.PCateOutput{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A wearable computer",
							CreatedAt:   myCreatedTime,
							UpdatedAt:   myUpdatedTime,
						},
						CreatedAt: myCreatedTime,
						UpdatedAt: myUpdatedTime,
					},
				},
			},
			expOutput: []*model.Product{
				{
					ID:          194,
					Name:        "iPhone 14",
					Description: "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
					Price:       1500,
					Quantity:    50,
					Author: &model.User{
						ID:        2,
						Name:      "Thuy Nguyen",
						Email:     "qthuy@gmail.com",
						CreatedAt: myCreatedTime.String(),
						UpdatedAt: myUpdatedTime.String(),
					},
					Category: &model.ProductCategory{
						ID:          90,
						Name:        "Smartwatch",
						Description: "A wearable computer",
						CreatedAt:   myCreatedTime.String(),
						UpdatedAt:   myUpdatedTime.String(),
					},
					CreatedAt: myCreatedTime.String(),
					UpdatedAt: myUpdatedTime.String(),
				},
				{
					ID:          195,
					Name:        "Macbook",
					Description: "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
					Price:       1500,
					Quantity:    50,
					Author: &model.User{
						ID:        2,
						Name:      "Thuy Nguyen",
						Email:     "qthuy@gmail.com",
						CreatedAt: myCreatedTime.String(),
						UpdatedAt: myUpdatedTime.String(),
					},
					Category: &model.ProductCategory{
						ID:          90,
						Name:        "Smartwatch",
						Description: "A wearable computer",
						CreatedAt:   myCreatedTime.String(),
						UpdatedAt:   myUpdatedTime.String(),
					},
					CreatedAt: myCreatedTime.String(),
					UpdatedAt: myUpdatedTime.String(),
				},
			},
		},
		"get all products with filter request successfully": {
			expCall:         true,
			givenFilterName: "iphone",
			mockProductCtrl: mockProductCtrl{
				query: controllers.ProductCtrlFilter{
					Name: "iphone",
				},
				output: []controllers.ProductOutputGraph{
					{
						ID:          194,
						Name:        "iPhone 14",
						Description: "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:       decimal.New(1500, 0),
						Quantity:    50,
						Author: controllers.UserOutput{
							ID:        2,
							Name:      "Thuy Nguyen",
							Email:     "qthuy@gmail.com",
							CreatedAt: myCreatedTime,
							UpdatedAt: myUpdatedTime,
						},
						Category: controllers.PCateOutput{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A wearable computer",
							CreatedAt:   myCreatedTime,
							UpdatedAt:   myUpdatedTime,
						},
						CreatedAt: myCreatedTime,
						UpdatedAt: myUpdatedTime,
					},
				},
			},
			expOutput: []*model.Product{
				{
					ID:          194,
					Name:        "iPhone 14",
					Description: "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
					Price:       1500,
					Quantity:    50,
					Author: &model.User{
						ID:        2,
						Name:      "Thuy Nguyen",
						Email:     "qthuy@gmail.com",
						CreatedAt: myCreatedTime.String(),
						UpdatedAt: myUpdatedTime.String(),
					},
					Category: &model.ProductCategory{
						ID:          90,
						Name:        "Smartwatch",
						Description: "A wearable computer",
						CreatedAt:   myCreatedTime.String(),
						UpdatedAt:   myUpdatedTime.String(),
					},
					CreatedAt: myCreatedTime.String(),
					UpdatedAt: myUpdatedTime.String(),
				},
			},
		},
		"products not found with filter": {
			expCall: true,
			mockProductCtrl: mockProductCtrl{
				query: controllers.ProductCtrlFilter{
					Name: "imac",
				},
			},
			givenFilterName: "imac",
		},
		"bad date format filter": {
			expCall:         false,
			givenFilterDate: "2023-06-999",
			expErr:          ErrDateBadRequest.Error(),
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			resolver := Resolver{Controller: mockController}

			if tc.expCall {
				mockController.On("GetProductsGraph", context.Background(), tc.mockProductCtrl.query).Return(tc.mockProductCtrl.output, tc.mockProductCtrl.err)
			}
			result, err := resolver.Query().GetProducts(context.Background(), tc.givenFilterName, tc.givenFilterDate)

			if err != nil {
				assert.EqualError(t, err, tc.expErr)
			} else {
				assert.Equal(t, tc.expOutput, result)
			}

			if tc.expCall {
				mockController.AssertCalled(t, "GetProductsGraph", context.Background(), tc.mockProductCtrl.query)
			}
		})
	}
}
