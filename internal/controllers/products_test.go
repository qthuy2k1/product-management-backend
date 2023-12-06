package controllers

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/qthuy2k1/product-management/internal/repositories"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_ProductControler_CreateProduct(t *testing.T) {
	type mockProductRepo struct {
		productInput repositories.Product
		err          error
	}

	type mockUserRepo struct {
		userID int
		output models.User
		err    error
	}

	type mockPCateRepo struct {
		pCateID int
		output  models.ProductCategory
		err     error
	}

	testCases := map[string]struct {
		expCall         bool
		productInput    ProductInput
		mockProductRepo mockProductRepo
		mockUserRepo    mockUserRepo
		mockPCateRepo   mockPCateRepo
		expErr          error
	}{
		"success": {
			expCall: true,
			productInput: ProductInput{
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        decimal.New(1500, 20),
				Quantity:     20,
				CategoryName: "Smartphone",
				AuthorID:     1,
			},
			mockProductRepo: mockProductRepo{
				productInput: repositories.Product{
					Name:        "iPhone 14",
					Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:       decimal.New(1500, 20),
					Quantity:    20,
					CategoryID:  1,
					AuthorID:    1,
				},
			},
			mockUserRepo: mockUserRepo{
				userID: 1,
				output: models.User{
					ID:       1,
					Name:     "Thuy Nguyen",
					Email:    "qthuy@gmail.com",
					Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
				},
			},
			mockPCateRepo: mockPCateRepo{
				pCateID: 1,
				output: models.ProductCategory{
					ID:          1,
					Name:        "Cellphone",
					Description: "Cellphone",
				},
			},
		},
		"user not found": {
			expCall: true,
			productInput: ProductInput{
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        decimal.New(1500, 20),
				Quantity:     20,
				CategoryName: "Smartphone",
				AuthorID:     100,
			},
			mockUserRepo: mockUserRepo{
				err: repositories.ErrUserNotFound,
			},
			expErr: ErrUserNotFound,
		},
		"product category not found": {
			expCall: true,
			productInput: ProductInput{
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        decimal.New(1500, 20),
				Quantity:     20,
				CategoryName: "smartwatchhhh",
				AuthorID:     1,
			},
			mockUserRepo: mockUserRepo{
				userID: 1,
				output: models.User{
					ID:       1,
					Name:     "Thuy Nguyen",
					Email:    "qthuy@gmail.com",
					Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
				},
			},
			mockPCateRepo: mockPCateRepo{
				pCateID: 100,
				err:     repositories.ErrProductCategoryNotFound,
			},
			expErr: ErrProductCategoryNotFound,
		},
		"create product with name longer than 255 characters": {
			expCall: true,
			productInput: ProductInput{
				Name:         "iPhone 14 Pro Maxxxx" + strings.Repeat("a", 251),
				Description:  "An Apple cellphone with A16 Bionic chip, 6GB RAM and 128GB storage",
				Price:        decimal.New(1500, 20),
				Quantity:     20,
				CategoryName: "Smartphone",
				AuthorID:     1,
			},
			mockProductRepo: mockProductRepo{
				productInput: repositories.Product{
					Name:        "iPhone 14 Pro Maxxxx" + strings.Repeat("a", 251),
					Description: "An Apple cellphone with A16 Bionic chip, 6GB RAM and 128GB storage",
					Price:       decimal.New(1500, 20),
					Quantity:    20,
					CategoryID:  1,
					AuthorID:    1,
				},
				err: errors.New("models: unable to insert into products: pq: value too long for type character varying(255)"),
			},
			mockUserRepo: mockUserRepo{
				userID: 1,
				output: models.User{
					ID:       1,
					Name:     "Thuy Nguyen",
					Email:    "qthuy@gmail.com",
					Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
				},
			},
			mockPCateRepo: mockPCateRepo{
				pCateID: 1,
				output: models.ProductCategory{
					ID:          1,
					Name:        "Cellphone",
					Description: "Cellphone",
				},
			},
			expErr: errors.New("models: unable to insert into products: pq: value too long for type character varying(255)"),
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockRepo := &repositories.MockIRepository{}
			controller := NewController(mockRepo)

			if tc.expCall {
				mockRepo.On("GetUser", context.Background(), tc.productInput.AuthorID).Return(tc.mockUserRepo.output, tc.mockUserRepo.err)
				mockRepo.On("GetProductCategoryByName", context.Background(), tc.productInput.CategoryName).Return(tc.mockPCateRepo.output, tc.mockPCateRepo.err)
				mockRepo.On("CreateProduct", context.Background(), tc.mockProductRepo.productInput).Return(tc.expErr)
			}

			if err := controller.CreateProduct(context.Background(), tc.productInput); tc.expErr != nil {
				assert.EqualError(t, err, tc.expErr.Error())
			}
		})
	}
}
func Test_ProductControler_UpdateProduct(t *testing.T) {
	type mockGetProductRepo struct {
		productID int
		output    models.Product
		err       error
	}

	type mockUpdateProductRepo struct {
		productInput models.Product
		err          error
	}

	type mockUserRepo struct {
		userID int
		output models.User
		err    error
	}

	type mockPCateRepo struct {
		pCateID int
		output  models.ProductCategory
		err     error
	}

	testCases := map[string]struct {
		expCall               bool
		productInput          ProductInput
		mockGetProductRepo    mockGetProductRepo
		mockUpdateProductRepo mockUpdateProductRepo
		mockUserRepo          mockUserRepo
		mockPCateRepo         mockPCateRepo
		expErr                error
	}{
		"success": {
			expCall: true,
			productInput: ProductInput{
				ID:           1,
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        decimal.New(1500, 20),
				Quantity:     20,
				CategoryName: "Smartphone",
				AuthorID:     1,
			},
			mockUpdateProductRepo: mockUpdateProductRepo{
				productInput: models.Product{
					ID:          1,
					Name:        "iPhone 14",
					Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:       decimal.New(1500, 20),
					Quantity:    20,
					CategoryID:  1,
					AuthorID:    1,
				},
			},
			mockGetProductRepo: mockGetProductRepo{
				productID: 1,
				output: models.Product{
					ID:          1,
					Name:        "iPhone 14",
					Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:       decimal.New(1500, 20),
					Quantity:    20,
					CategoryID:  1,
					AuthorID:    1,
				},
			},
			mockUserRepo: mockUserRepo{
				userID: 1,
				output: models.User{
					ID:       1,
					Name:     "Thuy Nguyen",
					Email:    "qthuy@gmail.com",
					Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
				},
			},
			mockPCateRepo: mockPCateRepo{
				pCateID: 1,
				output: models.ProductCategory{
					ID:          1,
					Name:        "Cellphone",
					Description: "Cellphone",
				},
			},
		},
		"user not found": {
			expCall: true,
			productInput: ProductInput{
				ID:           1,
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        decimal.New(1500, 20),
				Quantity:     20,
				CategoryName: "Smartphone",
				AuthorID:     100,
			},
			mockUpdateProductRepo: mockUpdateProductRepo{},
			mockGetProductRepo: mockGetProductRepo{
				productID: 1,
				output: models.Product{
					ID:          1,
					Name:        "iPhone 14",
					Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:       decimal.New(1500, 20),
					Quantity:    20,
					CategoryID:  1,
					AuthorID:    1,
				},
			},
			mockUserRepo: mockUserRepo{
				err: repositories.ErrUserNotFound,
			},
			mockPCateRepo: mockPCateRepo{},
			expErr:        ErrUserNotFound,
		},
		"product category not found": {
			expCall: true,
			productInput: ProductInput{
				ID:           1,
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        decimal.New(1500, 20),
				Quantity:     20,
				CategoryName: "smartwatchhhh",
				AuthorID:     1,
			},
			mockUpdateProductRepo: mockUpdateProductRepo{},
			mockGetProductRepo: mockGetProductRepo{
				productID: 1,
				output: models.Product{
					ID:          1,
					Name:        "iPhone 14",
					Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:       decimal.New(1500, 20),
					Quantity:    20,
					CategoryID:  1,
					AuthorID:    1,
				},
			},
			mockUserRepo: mockUserRepo{
				userID: 1,
				output: models.User{
					ID:       1,
					Name:     "Thuy Nguyen",
					Email:    "qthuy@gmail.com",
					Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
				},
			},
			mockPCateRepo: mockPCateRepo{
				pCateID: 100,
				err:     repositories.ErrProductCategoryNotFound,
			},
			expErr: ErrProductCategoryNotFound,
		},
		"create product with name longer than 255 characters": {
			expCall: true,
			productInput: ProductInput{
				ID:           1,
				Name:         "iPhone 14 Pro Maxxxx" + strings.Repeat("a", 251),
				Description:  "An Apple cellphone with A16 Bionic chip, 6GB RAM and 128GB storage",
				Price:        decimal.New(1500, 20),
				Quantity:     20,
				CategoryName: "Smartphone",
				AuthorID:     1,
			},
			mockUpdateProductRepo: mockUpdateProductRepo{
				productInput: models.Product{
					ID:          1,
					Name:        "iPhone 14 Pro Maxxxx" + strings.Repeat("a", 251),
					Description: "An Apple cellphone with A16 Bionic chip, 6GB RAM and 128GB storage",
					Price:       decimal.New(1500, 20),
					Quantity:    20,
					CategoryID:  1,
					AuthorID:    1,
				},
				err: errors.New("models: unable to insert into products: pq: value too long for type character varying(255)"),
			},
			mockGetProductRepo: mockGetProductRepo{
				productID: 1,
				output: models.Product{
					ID:          1,
					Name:        "iPhone 14",
					Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:       decimal.New(1500, 20),
					Quantity:    20,
					CategoryID:  1,
					AuthorID:    1,
				},
			},
			mockUserRepo: mockUserRepo{
				userID: 1,
				output: models.User{
					ID:       1,
					Name:     "Thuy Nguyen",
					Email:    "qthuy@gmail.com",
					Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
				},
			},
			mockPCateRepo: mockPCateRepo{
				pCateID: 1,
				output: models.ProductCategory{
					ID:          1,
					Name:        "Cellphone",
					Description: "Cellphone",
				},
			},
			expErr: errors.New("models: unable to insert into products: pq: value too long for type character varying(255)"),
		},
		"product not found": {
			expCall: true,
			productInput: ProductInput{
				ID:           1000,
				Name:         "iPhone 14",
				Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
				Price:        decimal.New(1500, 20),
				Quantity:     20,
				CategoryName: "Smartphone",
				AuthorID:     1,
			},
			mockUpdateProductRepo: mockUpdateProductRepo{
				productInput: models.Product{
					ID:          1,
					Name:        "iPhone 14",
					Description: "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:       decimal.New(1500, 20),
					Quantity:    20,
					CategoryID:  1,
					AuthorID:    1,
				},
			},
			mockGetProductRepo: mockGetProductRepo{
				productID: 1000,
				err:       errors.New("product not found"),
			},
			expErr: errors.New("product not found"),
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockRepo := &repositories.MockIRepository{}
			controller := NewController(mockRepo)
			if tc.expCall {
				mockRepo.On("GetProduct", context.Background(), tc.mockGetProductRepo.productID).Return(tc.mockGetProductRepo.output, tc.mockGetProductRepo.err)
				mockRepo.On("GetUser", context.Background(), tc.productInput.AuthorID).Return(tc.mockUserRepo.output, tc.mockUserRepo.err)
				mockRepo.On("GetProductCategoryByName", context.Background(), tc.productInput.CategoryName).Return(tc.mockPCateRepo.output, tc.mockPCateRepo.err)
				mockRepo.On("UpdateProduct", context.Background(), tc.mockUpdateProductRepo.productInput).Return(tc.expErr)
			}

			if err := controller.UpdateProduct(context.Background(), tc.productInput); tc.expErr != nil {
				assert.EqualError(t, err, tc.expErr.Error())
			}
		})
	}
}
func Test_ProductController_DeleteProduct(t *testing.T) {
	type mockProductRepo struct {
		expCall   bool
		productID int
		err       error
	}
	tests := map[string]struct {
		mockProductRepo mockProductRepo
		input           int
		err             error
	}{
		"success": {
			mockProductRepo: mockProductRepo{
				expCall:   true,
				productID: 2,
			},
			input: 2,
		},
		"product not found": {
			mockProductRepo: mockProductRepo{
				expCall:   true,
				productID: -1,
				err:       repositories.ErrProductNotFound,
			},
			input: -1,
			err:   ErrProductNotFound,
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.MockIRepository{}
			controller := NewController(&mockRepo)
			if tc.mockProductRepo.expCall {
				mockRepo.On("DeleteProduct", context.Background(), tc.mockProductRepo.productID).Return(tc.mockProductRepo.err)
			}
			err := controller.DeleteProduct(context.Background(), tc.input)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}

func Test_ProductController_GetProducts(t *testing.T) {
	myCreatedTime, err := time.Parse("2006-01-02 15:04:05.999999", "2023-06-02 09:08:36.046843")
	assert.NoError(t, err)
	myUpdatedTime, err := time.Parse("2006-01-02 15:04:05.999999", "2023-06-02 09:08:36.046843")
	assert.NoError(t, err)

	type mockProductRepo struct {
		input  repositories.ProductRepoFilter
		output []repositories.ProductOutput
		err    error
	}

	tests := map[string]struct {
		expCall         bool
		mockProductRepo mockProductRepo
		input           ProductCtrlFilter
		output          []ProductOutput
		err             error
	}{
		"get products successfully": {
			expCall: true,
			mockProductRepo: mockProductRepo{
				output: []repositories.ProductOutput{
					{
						ID:           194,
						Name:         "iPhone 14",
						Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:        decimal.New(1500, 20),
						Quantity:     50,
						CategoryName: "Smartphone",
						AuthorID:     136,
						CreatedAt:    myCreatedTime,
						UpdatedAt:    myUpdatedTime,
					},
					{
						ID:           195,
						Name:         "Macbook",
						Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:        decimal.New(1500, 20),
						Quantity:     50,
						CategoryName: "Smartphone",
						AuthorID:     136,
						CreatedAt:    myCreatedTime,
						UpdatedAt:    myUpdatedTime,
					},
				},
			},
			output: []ProductOutput{
				{
					ID:           194,
					Name:         "iPhone 14",
					Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
					Price:        decimal.New(1500, 20),
					Quantity:     50,
					CategoryName: "Smartphone",
					AuthorID:     136,
					CreatedAt:    myCreatedTime,
					UpdatedAt:    myUpdatedTime,
				},
				{
					ID:           195,
					Name:         "Macbook",
					Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
					Price:        decimal.New(1500, 20),
					Quantity:     50,
					CategoryName: "Smartphone",
					AuthorID:     136,
					CreatedAt:    myCreatedTime,
					UpdatedAt:    myUpdatedTime,
				},
			},
		},
		"get products successfully with filter": {
			expCall: true,
			mockProductRepo: mockProductRepo{
				input: repositories.ProductRepoFilter{
					Name: "iPhone",
				},
				output: []repositories.ProductOutput{
					{
						ID:           194,
						Name:         "iPhone 14",
						Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:        decimal.New(1500, 20),
						Quantity:     50,
						CategoryName: "Smartphone",
						AuthorID:     136,
						CreatedAt:    myCreatedTime,
						UpdatedAt:    myUpdatedTime,
					},
				},
			},
			input: ProductCtrlFilter{
				Name: "iPhone",
			},
			output: []ProductOutput{
				{
					ID:           194,
					Name:         "iPhone 14",
					Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
					Price:        decimal.New(1500, 20),
					Quantity:     50,
					CategoryName: "Smartphone",
					AuthorID:     136,
					CreatedAt:    myCreatedTime,
					UpdatedAt:    myUpdatedTime,
				},
			},
		},
		"products not found with filter": {
			expCall: true,
			mockProductRepo: mockProductRepo{
				input: repositories.ProductRepoFilter{
					Name: "imac",
				},
			},
			input: ProductCtrlFilter{
				Name: "imac",
			},
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.MockIRepository{}
			controller := NewController(&mockRepo)
			if tc.expCall {
				mockRepo.On("GetProducts", context.Background(), tc.mockProductRepo.input).Return(tc.mockProductRepo.output, tc.mockProductRepo.err)
			}
			products, err := controller.GetProducts(context.Background(), tc.input)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.Equal(t, products, tc.output)
			}
		})
	}
}

func Test_ProductController_ImportProductsFromCSV(t *testing.T) {
	type mockUserRepo struct {
		output models.User
		err    error
	}

	userDefaultRepo := models.User{
		ID:       1,
		Name:     "user default",
		Email:    "user@gmail.com",
		Password: "userpassword",
		Role:     "user",
		Status:   "activated",
	}

	type mockPCateRepo struct {
		output models.ProductCategory
		err    error
	}

	pCateDefaultRepo := models.ProductCategory{
		ID:          1,
		Name:        "Uncategorized",
		Description: "Product is uncategorized",
	}

	type mockProductRepo struct {
		input repositories.Product
	}

	type mockInput struct {
		mockProductRepo mockProductRepo
		mockPCateRepo   mockPCateRepo
		mockUserRepo    mockUserRepo
	}

	testCases := map[string]struct {
		expCall  bool
		input    []ProductInput
		mockRepo []mockInput
		csvData  []string
		err      error
	}{
		"import products successfully": {
			expCall: true,
			input: []ProductInput{
				{
					Name:         "Samsung S23",
					Description:  "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
					Price:        decimal.New(1500, 20),
					Quantity:     10,
					AuthorID:     2,
					CategoryName: "Smartphone",
				},
				{
					Name:         "Samsung Galaxy Watch 5",
					Description:  "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
					Price:        decimal.New(1500, 20),
					Quantity:     20,
					AuthorID:     2,
					CategoryName: "Smartwatch",
				},
			},
			mockRepo: []mockInput{
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung S23",
							Description: "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
							Price:       decimal.New(1500, 20),
							Quantity:    10,
							AuthorID:    2,
							CategoryID:  80,
						},
					},

					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          80,
							Name:        "Smartphone",
							Description: "A mobile device",
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
				},
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung Galaxy Watch 5",
							Description: "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
							Price:       decimal.New(1500, 20),
							Quantity:    20,
							AuthorID:    2,
							CategoryID:  90,
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A mobile device",
						},
					},
				},
			},
			csvData: []string{`Name,Description,Price,Quantity,AuthorID,Category`,
				`Samsung S23,"A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1100,10,2,Smartphone`,
				`Samsung Galaxy Watch 5,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",1500,20,2,Smartwatch`},
		},
		"change the sorting order of the column": {
			expCall: true,
			input: []ProductInput{
				{
					Name:         "Samsung S23",
					Description:  "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
					Price:        decimal.New(1500, 20),
					Quantity:     10,
					AuthorID:     2,
					CategoryName: "Smartphone",
				},
				{
					Name:         "Samsung Galaxy Watch 5",
					Description:  "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
					Price:        decimal.New(1500, 20),
					Quantity:     20,
					AuthorID:     2,
					CategoryName: "Smartwatch",
				},
			},
			mockRepo: []mockInput{
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung S23",
							Description: "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
							Price:       decimal.New(1500, 20),
							Quantity:    10,
							AuthorID:    2,
							CategoryID:  80,
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          80,
							Name:        "Smartphone",
							Description: "A mobile device",
						},
					},
				},
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung Galaxy Watch 5",
							Description: "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
							Price:       decimal.New(1500, 20),
							Quantity:    20,
							AuthorID:    2,
							CategoryID:  90,
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A mobile device",
						},
					},
				},
			},
			csvData: []string{`Name,Description,Price,Quantity,Category,AuthorID`,
				`Samsung S23,"A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1500,10,Smartphone,2`,
				`Samsung Galaxy Watch 5,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",1500,20,Smartwatch,2`},
		},
		"import products with product category not found": {
			expCall: true,
			input: []ProductInput{
				{
					Name:         "Samsung S23",
					Description:  "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
					Price:        decimal.New(1500, 20),
					Quantity:     10,
					AuthorID:     2,
					CategoryName: "Smartphone",
				},
				{
					Name:         "Samsung Galaxy Watch 5",
					Description:  "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
					Price:        decimal.New(1500, 20),
					Quantity:     20,
					AuthorID:     2,
					CategoryName: "Smartwatchhhhh",
				},
			},
			mockRepo: []mockInput{
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung S23",
							Description: "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
							Price:       decimal.New(1500, 20),
							Quantity:    10,
							AuthorID:    2,
							CategoryID:  80,
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          80,
							Name:        "Smartphone",
							Description: "A mobile device",
						},
					},
				},
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung Galaxy Watch 5",
							Description: "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
							Price:       decimal.New(1500, 20),
							Quantity:    20,
							AuthorID:    2,
							CategoryID:  1,
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
					mockPCateRepo: mockPCateRepo{
						output: pCateDefaultRepo,
					},
				},
			},
			csvData: []string{`Name,Description,Price,Quantity,AuthorID,Category`,
				`Samsung S23,"A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1500,10,2,Smartphone`,
				`Samsung Galaxy Watch 5,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",1500,20,2,Smartwatchhhhh`},
		},
		"import products with user not found": {
			expCall: true,
			input: []ProductInput{
				{
					Name:         "Samsung S23",
					Description:  "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
					Price:        decimal.New(1500, 20),
					Quantity:     10,
					AuthorID:     1000,
					CategoryName: "Smartphone",
				},
				{
					Name:         "Samsung Galaxy Watch 5",
					Description:  "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
					Price:        decimal.New(1500, 20),
					Quantity:     20,
					AuthorID:     2,
					CategoryName: "Smartwatch",
				},
			},
			mockRepo: []mockInput{
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung S23",
							Description: "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
							Price:       decimal.New(1500, 20),
							Quantity:    10,
							AuthorID:    1000,
							CategoryID:  80,
						},
					},
					mockUserRepo: mockUserRepo{
						output: userDefaultRepo,
					},
					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          80,
							Name:        "Smartphone",
							Description: "A mobile device",
						},
					},
				},
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung Galaxy Watch 5",
							Description: "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
							Price:       decimal.New(1500, 20),
							Quantity:    20,
							AuthorID:    2,
							CategoryID:  90,
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A wearable computer",
						},
					},
				},
			},
			csvData: []string{`Name,Description,Price,Quantity,AuthorID,Category`,
				`Samsung S23,"A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1100,10,1500,Smartphone`,
				`Samsung Galaxy Watch 5,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",1500,20,2,Smartwatch`},
		},
		"data has caret and double quotes": {
			expCall: true,
			input: []ProductInput{
				{
					Name:         `Samsung S23""`,
					Description:  "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
					Price:        decimal.New(1500, 20),
					Quantity:     10,
					AuthorID:     2,
					CategoryName: "Smartphone",
				},
				{
					Name:         "Samsung Galaxy Watch 5<>^",
					Description:  "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
					Price:        decimal.New(1500, 20),
					Quantity:     20,
					AuthorID:     2,
					CategoryName: "Smartwatch",
				},
			},
			mockRepo: []mockInput{
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung S23",
							Description: "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
							Price:       decimal.New(1500, 20),
							Quantity:    10,
							AuthorID:    2,
							CategoryID:  80,
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          80,
							Name:        "Smartphone",
							Description: "A mobile device",
						},
					},
				},
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung Galaxy Watch 5",
							Description: "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
							Price:       decimal.New(1500, 20),
							Quantity:    20,
							AuthorID:    2,
							CategoryID:  90,
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A wearable computer",
						},
					},
				},
			},
			csvData: []string{`Name,Description,Price,Quantity,AuthorID,Category`,
				`"Samsung S23""""","A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1500,10,2,Smartphone`,
				`Samsung Galaxy Watch 5<>^,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",1500,20,2,Smartwatch`},
		},
		"invalid csv format": {
			expCall: false,
			csvData: []string{`Name,Description,Price,Quantity,AuthorID`,
				`"Samsung S23""""","A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1100,10,1,1`,
				`Samsung Galaxy Watch 5<>^,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",250,20,1,1`},
			err: ErrCSVFileFormat,
		},
		"data not enough columns": {
			expCall: false,
			csvData: []string{`Name,Description,Price,Quantity,AuthorID`,
				`"Samsung S23""""","A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1100,10,1`,
				`Samsung Galaxy Watch 5<>^,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",250,20,1`},
			err: ErrNotEnoughColumns,
		},
		"missing price column data": {
			expCall: true,
			input: []ProductInput{
				{
					Name:         "Samsung Galaxy Watch 5<>^",
					Description:  "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
					Price:        decimal.New(1500, 20),
					Quantity:     20,
					AuthorID:     2,
					CategoryName: "Smartwatch",
				},
			},
			mockRepo: []mockInput{
				{
					mockProductRepo: mockProductRepo{
						input: repositories.Product{
							Name:        "Samsung Galaxy Watch 5",
							Description: "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
							Price:       decimal.New(1500, 20),
							Quantity:    20,
							AuthorID:    2,
							CategoryID:  90,
						},
					},
					mockUserRepo: mockUserRepo{
						output: models.User{
							ID:       2,
							Name:     "Thuy Nguyen",
							Email:    "qthuy@gmail.com",
							Password: "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
						},
					},
					mockPCateRepo: mockPCateRepo{
						output: models.ProductCategory{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A wearable computer",
						},
					},
				},
			},
			csvData: []string{`Name,Description,Price,Quantity,AuthorID,Category`,
				`Samsung S23,"A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",,10,1500,Smartphone`,
				`Samsung Galaxy Watch 5,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",1500,20,2,Smartwatch`},
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.MockIRepository{}
			controller := NewController(&mockRepo)

			var pRepoInput []repositories.Product

			f, err := os.CreateTemp("", desc+"products.csv")
			assert.NoError(t, err)

			// clean up
			defer os.Remove(f.Name())

			for _, s := range tc.csvData {
				_, err := f.WriteString(s + "\n")
				assert.NoError(t, err)
			}

			f.Close()

			csvData, err := os.Open(f.Name())
			assert.NoError(t, err)

			if tc.expCall {
				mockRepo.On("GetUserDefault", context.Background()).Return(userDefaultRepo, nil)
				mockRepo.On("GetProductCategoryDefault", context.Background()).Return(pCateDefaultRepo, nil)

				for i := range tc.mockRepo {
					mockRepo.On("GetUser", context.Background(), tc.input[i].AuthorID).Return(tc.mockRepo[i].mockUserRepo.output, tc.mockRepo[i].mockUserRepo.err)
					mockRepo.On("GetProductCategoryByName", context.Background(), tc.input[i].CategoryName).Return(tc.mockRepo[i].mockPCateRepo.output, tc.mockRepo[i].mockPCateRepo.err)

					pRepoInput = append(pRepoInput, tc.mockRepo[i].mockProductRepo.input)
				}
				mockRepo.On("ImportProducts", context.Background(), pRepoInput).Return(tc.err)
			}

			if err = controller.ImportProductsFromCSV(context.Background(), csvData); tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}

func Test_ProductController_GetProductsGraph(t *testing.T) {
	myCreatedTime, err := time.Parse("2006-01-02 15:04:05.999999", "2023-06-02 09:08:36.046843")
	assert.NoError(t, err)
	myUpdatedTime, err := time.Parse("2006-01-02 15:04:05.999999", "2023-06-02 09:08:36.046843")
	assert.NoError(t, err)

	type mockProductRepo struct {
		input  repositories.ProductRepoFilter
		output []repositories.ProductOutput
		err    error
	}

	type mockUserRepo struct {
		userID int
		output repositories.User
		err    error
	}

	type mockPCateRepo struct {
		pCateID int
		output  repositories.ProductCategory
		err     error
	}

	type mockRepo struct {
		mockUserRepo  mockUserRepo
		mockPCateRepo mockPCateRepo
	}

	tests := map[string]struct {
		expCall         bool
		mockProductRepo mockProductRepo
		mockRepo        []mockRepo
		input           ProductCtrlFilter
		output          []ProductOutput
		err             error
	}{
		"get products successfully": {
			expCall: true,
			mockProductRepo: mockProductRepo{
				output: []repositories.ProductOutput{
					{
						ID:           194,
						Name:         "iPhone 14",
						Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:        decimal.New(1500, 20),
						Quantity:     50,
						CategoryName: "Smartphone",
						AuthorID:     136,
						CreatedAt:    myCreatedTime,
						UpdatedAt:    myUpdatedTime,
					},
					{
						ID:           195,
						Name:         "Macbook",
						Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:        decimal.New(1500, 20),
						Quantity:     50,
						CategoryName: "Smartphone",
						AuthorID:     136,
						CreatedAt:    myCreatedTime,
						UpdatedAt:    myUpdatedTime,
					},
				},
			},
			mockRepo: []mockRepo{
				{
					mockUserRepo: mockUserRepo{
						userID: 100,
						output: repositories.User{
							ID:        2,
							Name:      "Thuy Nguyen",
							Email:     "qthuy@gmail.com",
							CreatedAt: myCreatedTime,
							UpdatedAt: myUpdatedTime,
						},
					},
					mockPCateRepo: mockPCateRepo{
						pCateID: 100,
						output: repositories.ProductCategory{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A wearable computer",
							CreatedAt:   myCreatedTime,
							UpdatedAt:   myUpdatedTime,
						},
					},
				},
				{
					mockUserRepo: mockUserRepo{
						userID: 100,
						output: repositories.User{
							ID:        2,
							Name:      "Thuy Nguyen",
							Email:     "qthuy@gmail.com",
							CreatedAt: myCreatedTime,
							UpdatedAt: myUpdatedTime,
						},
					},
					mockPCateRepo: mockPCateRepo{
						pCateID: 100,
						output: repositories.ProductCategory{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A wearable computer",
							CreatedAt:   myCreatedTime,
							UpdatedAt:   myUpdatedTime,
						},
					},
				},
			},
			output: []ProductOutput{
				{
					ID:           194,
					Name:         "iPhone 14",
					Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
					Price:        decimal.New(1500, 20),
					Quantity:     50,
					CategoryName: "Smartphone",
					AuthorID:     136,
					CreatedAt:    myCreatedTime,
					UpdatedAt:    myUpdatedTime,
				},
				{
					ID:           195,
					Name:         "Macbook",
					Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
					Price:        decimal.New(1500, 20),
					Quantity:     50,
					CategoryName: "Smartphone",
					AuthorID:     136,
					CreatedAt:    myCreatedTime,
					UpdatedAt:    myUpdatedTime,
				},
			},
		},
		"get products successfully with filter": {
			expCall: true,
			mockProductRepo: mockProductRepo{
				input: repositories.ProductRepoFilter{
					Name: "iPhone",
				},
				output: []repositories.ProductOutput{
					{
						ID:           194,
						Name:         "iPhone 14",
						Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:        decimal.New(1500, 20),
						Quantity:     50,
						CategoryName: "Smartphone",
						AuthorID:     136,
						CreatedAt:    myCreatedTime,
						UpdatedAt:    myUpdatedTime,
					},
				},
			},
			mockRepo: []mockRepo{
				{
					mockUserRepo: mockUserRepo{
						userID: 100,
						output: repositories.User{
							ID:        2,
							Name:      "Thuy Nguyen",
							Email:     "qthuy@gmail.com",
							CreatedAt: myCreatedTime,
							UpdatedAt: myUpdatedTime,
						},
					},
					mockPCateRepo: mockPCateRepo{
						pCateID: 100,
						output: repositories.ProductCategory{
							ID:          90,
							Name:        "Smartwatch",
							Description: "A wearable computer",
							CreatedAt:   myCreatedTime,
							UpdatedAt:   myUpdatedTime,
						},
					},
				},
			},
			input: ProductCtrlFilter{
				Name: "iPhone",
			},
			output: []ProductOutput{
				{
					ID:           194,
					Name:         "iPhone 14",
					Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
					Price:        decimal.New(1500, 20),
					Quantity:     50,
					CategoryName: "Smartphone",
					AuthorID:     136,
					CreatedAt:    myCreatedTime,
					UpdatedAt:    myUpdatedTime,
				},
			},
		},
		"products not found with filter": {
			expCall: true,
			mockProductRepo: mockProductRepo{
				input: repositories.ProductRepoFilter{
					Name: "imac",
				},
			},
			input: ProductCtrlFilter{
				Name: "imac",
			},
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.MockIRepository{}
			controller := NewController(&mockRepo)
			if tc.expCall {
				mockRepo.On("GetProducts", context.Background(), tc.mockProductRepo.input).Return(tc.mockProductRepo.output, tc.mockProductRepo.err)
				for i := range tc.mockProductRepo.output {
					mockRepo.On("GetUser", context.Background(), tc.mockProductRepo.output[i].AuthorID).Return(tc.mockRepo[i].mockUserRepo.output, tc.mockRepo[i].mockUserRepo.err)
					mockRepo.On("GetProductCategoryByName", context.Background(), tc.mockProductRepo.output[i].CategoryName).Return(tc.mockRepo[i].mockPCateRepo.output, tc.mockRepo[i].mockPCateRepo.err)
				}
			}
			products, err := controller.GetProducts(context.Background(), tc.input)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.Equal(t, products, tc.output)
			}
		})
	}
}

func Test_ProductController_ExportProductsToCSV(t *testing.T) {
	type mockProductRepo struct {
		input  repositories.ProductRepoFilter
		output []repositories.ProductOutput
		err    error
	}

	testCases := map[string]struct {
		expCall         bool
		mockProductRepo mockProductRepo
		input           ProductCtrlFilter
		expData         [][]string
		err             error
	}{
		"export products successfully": {
			expCall: true,
			mockProductRepo: mockProductRepo{
				output: []repositories.ProductOutput{
					{
						ID:           1,
						Name:         "Samsung S23",
						Description:  "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
						Price:        decimal.New(1500, 20),
						Quantity:     10,
						AuthorID:     2,
						CategoryName: "Smartphone",
					},
					{
						ID:           2,
						Name:         "Samsung Galaxy Watch 5",
						Description:  "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",
						Price:        decimal.New(1500, 20),
						Quantity:     20,
						AuthorID:     2,
						CategoryName: "Smartwatch",
					},
				},
			},
			expData: [][]string{
				{
					"ID", "Name", "Description", "Price", "Quantity", "AuthorID", "Category", "CreatedAt", "UpdatedAt",
				},
				{
					"1", "Samsung S23", "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.", "1100", "10", "2", "Smartphone", "0001-01-01 00:00:00 +0000 UTC", "0001-01-01 00:00:00 +0000 UTC",
				},
				{
					"2", "Samsung Galaxy Watch 5", "A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.", "250", "20", "2", "Smartwatch", "0001-01-01 00:00:00 +0000 UTC", "0001-01-01 00:00:00 +0000 UTC",
				},
			},
		},
		"export products successfully with filter": {
			expCall: true,
			mockProductRepo: mockProductRepo{
				input: repositories.ProductRepoFilter{
					Name: "S23",
				},
				output: []repositories.ProductOutput{
					{
						ID:           1,
						Name:         "Samsung S23",
						Description:  "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",
						Price:        decimal.New(1500, 20),
						Quantity:     10,
						AuthorID:     2,
						CategoryName: "Smartphone",
					},
				},
			},
			expData: [][]string{
				{
					"ID", "Name", "Description", "Price", "Quantity", "AuthorID", "Category", "CreatedAt", "UpdatedAt",
				},
				{
					"1", "Samsung S23", "A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.", "1100", "10", "2", "Smartphone", "0001-01-01 00:00:00 +0000 UTC", "0001-01-01 00:00:00 +0000 UTC",
				},
			},
			input: ProductCtrlFilter{
				Name: "S23",
			},
		},
		"export no products with filter": {
			expCall: true,
			mockProductRepo: mockProductRepo{
				input: repositories.ProductRepoFilter{
					Name: "aaa",
				},
			},
			expData: [][]string{
				{
					"ID", "Name", "Description", "Price", "Quantity", "AuthorID", "Category", "CreatedAt", "UpdatedAt",
				},
			},
			input: ProductCtrlFilter{
				Name: "aaa",
			},
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.MockIRepository{}
			controller := NewController(&mockRepo)

			if tc.expCall {
				mockRepo.On("GetProducts", context.Background(), tc.mockProductRepo.input).Return(tc.mockProductRepo.output, tc.mockProductRepo.err)
			}

			if _, err := controller.ExportProductsToCSV(context.Background(), tc.input); err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
