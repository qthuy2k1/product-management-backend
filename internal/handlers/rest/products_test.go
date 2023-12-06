package rest

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ProductHandler_CreateProduct(t *testing.T) {
	type mockProductCtrl struct {
		expCall      bool
		productInput controllers.ProductInput
		err          error
	}

	testCases := map[string]struct {
		givenInput      string // payload provided from end-user
		mockProductCtrl mockProductCtrl
		expResp         string
		expCode         int
	}{
		"create product successfully": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				productInput: controllers.ProductInput{
					Name:         "iPhone 14",
					Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:        decimal.New(1500, 0),
					Quantity:     20,
					AuthorID:     1,
					CategoryName: "Smartphone",
				},
			},
			givenInput: `{"name":"iPhone 14","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"author_id":1,"category":"Smartphone"}`,
			expResp:    `{"success":true}`,
			expCode:    http.StatusCreated,
		},
		"user not exists": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				productInput: controllers.ProductInput{
					Name:         "iPhone 14",
					Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:        decimal.New(1500, 0),
					Quantity:     20,
					AuthorID:     100,
					CategoryName: "Smartphone",
				},
				err: controllers.ErrUserNotFound,
			},
			givenInput: `{"name":"iPhone 14","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"author_id":100,"category":"Smartphone"}`,
			expResp:    `{"message":"user not found"}`,
			expCode:    http.StatusNotFound,
		},
		"product category not exists": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				productInput: controllers.ProductInput{
					Name:         "iPhone 14",
					Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:        decimal.New(1500, 0),
					Quantity:     20,
					AuthorID:     40,
					CategoryName: "smartwatchhhh",
				},
				err: controllers.ErrProductCategoryNotFound,
			},
			givenInput: `{"name":"iPhone 14","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"author_id":40,"category":"smartwatchhhh"}`,
			expResp:    `{"message":"product category not found"}`,
			expCode:    http.StatusNotFound,
		},
		"create product with missing name field": {
			givenInput: `{"description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":1}`,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"name cannot be blank"}`,
			expCode: http.StatusBadRequest,
		},
		"create product with name too long": {
			// 260 'a' characters
			givenInput: `{"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":1}`,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"name too long"}`,
			expCode: http.StatusBadRequest,
		},
		"create product with invalid JSON": {
			givenInput: `{"name":"iPhone 14","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":1,"author_id":1`,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid json"}`,
			expCode: http.StatusBadRequest,
		},
		"create product with wrong type field": {
			givenInput: `{"name":123123,"description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":1}`,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid json"}`,
			expCode: http.StatusBadRequest,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			handler := NewHandler(mockController)

			r := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(tc.givenInput))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tc.mockProductCtrl.expCall {
				mockController.On("CreateProduct", context.Background(), tc.mockProductCtrl.productInput).Return(tc.mockProductCtrl.err)
			}
			handler.CreateProduct(w, r)

			assert.Equal(t, tc.expResp, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expCode, w.Code)

			if tc.mockProductCtrl.expCall {
				mockController.AssertCalled(t, "CreateProduct", context.Background(), tc.mockProductCtrl.productInput)
			}
		})
	}
}

func Test_ProductHandler_UpdateProduct(t *testing.T) {
	type mockProductCtrl struct {
		expCall      bool
		productInput controllers.ProductInput
		err          error
	}

	testCases := map[string]struct {
		givenInput      string // payload provided from end-user
		productID       int
		mockProductCtrl mockProductCtrl
		expResp         string
		expCode         int
	}{
		"update product successfully": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				productInput: controllers.ProductInput{
					ID:           1,
					Name:         "iPhone 14",
					Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:        decimal.New(1500, 0),
					Quantity:     20,
					AuthorID:     1,
					CategoryName: "Smartphone",
				},
			},
			givenInput: `{"name":"iPhone 14","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":1}`,
			productID:  1,
			expResp:    `{"success":true}`,
			expCode:    http.StatusOK,
		},
		"user not exists": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				productInput: controllers.ProductInput{
					ID:           1,
					Name:         "iPhone 14",
					Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:        decimal.New(1500, 0),
					Quantity:     20,
					AuthorID:     100,
					CategoryName: "Smartphone",
				},
				err: controllers.ErrUserNotFound,
			},
			givenInput: `{"name":"iPhone 14","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":100}`,
			productID:  1,
			expResp:    `{"message":"user not found"}`,
			expCode:    http.StatusNotFound,
		},
		"product category not exists": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				productInput: controllers.ProductInput{
					ID:           1,
					Name:         "iPhone 14",
					Description:  "An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage",
					Price:        decimal.New(1500, 0),
					Quantity:     20,
					AuthorID:     40,
					CategoryName: "smartwatchhhh",
				},
				err: controllers.ErrProductCategoryNotFound,
			},
			givenInput: `{"name":"iPhone 14","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"smartwatchhhh","author_id":40}`,
			productID:  1,
			expResp:    `{"message":"product category not found"}`,
			expCode:    http.StatusNotFound,
		},
		"update product with missing name field": {
			givenInput: `{"description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":1}`,
			productID:  1,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"name cannot be blank"}`,
			expCode: http.StatusBadRequest,
		},
		"update product with name too long": {
			// 260 'a' characters
			givenInput: `{"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":1}`,
			productID:  1,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"name too long"}`,
			expCode: http.StatusBadRequest,
		},
		"update product with invalid JSON": {
			givenInput: `{"name":"iPhone 14","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":1`,
			productID:  1,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid json"}`,
			expCode: http.StatusBadRequest,
		},
		"update product with wrong type field": {
			givenInput: `{"name":123123,"description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":1}`,
			productID:  1,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid json"}`,
			expCode: http.StatusBadRequest,
		},
		"invalid product ID": {
			givenInput: `{"name":"iPhone 14","description":"An Apple cellphone with A15 Bionic chip, 6GB RAM and 128GB storage","price":1500,"quantity":20,"category":"Smartphone","author_id":1}`,
			productID:  0,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid product ID"}`,
			expCode: http.StatusBadRequest,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			handler := NewHandler(mockController)

			r := httptest.NewRequest(http.MethodPut, "/products", strings.NewReader(tc.givenInput))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// New route
			rctx := chi.NewRouteContext()
			// Add user id to url params
			rctx.URLParams.Add("productID", strconv.Itoa(tc.productID))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			if tc.mockProductCtrl.expCall {
				mockController.On("UpdateProduct", r.Context(), tc.mockProductCtrl.productInput).Return(tc.mockProductCtrl.err)
			}
			handler.UpdateProduct(w, r)

			assert.Equal(t, tc.expResp, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expCode, w.Code)

			if tc.mockProductCtrl.expCall {
				mockController.AssertCalled(t, "UpdateProduct", mock.Anything, tc.mockProductCtrl.productInput)
			}
		})
	}
}

func Test_ProductHandler_DeleteProduct(t *testing.T) {
	type mockProductCtrl struct {
		expCall bool
		err     error
	}
	// Set up the test cases
	testCases := map[string]struct {
		productID       int
		mockProductCtrl mockProductCtrl
		expResp         string
		expCode         int
	}{
		"delete product successfully": {
			productID: 1,
			mockProductCtrl: mockProductCtrl{
				expCall: true,
			},
			expResp: `{"success":true}`,
			expCode: http.StatusOK,
		},
		"invalid product ID": {
			productID: 0,
			mockProductCtrl: mockProductCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid product ID"}`,
			expCode: http.StatusBadRequest,
		},
		"product not found": {
			productID: 100,
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				err:     controllers.ErrProductNotFound,
			},
			expResp: `{"message":"product not found"}`,
			expCode: http.StatusNotFound,
		},
	}

	// Test each test case
	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			handler := NewHandler(mockController)

			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/products/%d", tc.productID), nil)
			w := httptest.NewRecorder()

			// New route
			rctx := chi.NewRouteContext()
			// Add user id to url params
			rctx.URLParams.Add("productID", strconv.Itoa(tc.productID))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			if tc.mockProductCtrl.expCall {
				mockController.On("DeleteProduct", mock.Anything, tc.productID).Return(tc.mockProductCtrl.err)
			}
			handler.DeleteProduct(w, r)

			assert.Equal(t, tc.expResp, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expCode, w.Code)

			if tc.mockProductCtrl.expCall {
				mockController.AssertCalled(t, "DeleteProduct", mock.Anything, tc.productID)
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
		expCall bool
		query   controllers.ProductCtrlFilter
		output  []controllers.ProductOutput
		err     error
	}

	testCases := map[string]struct {
		givenFilter     string // payload provided from end-user
		mockProductCtrl mockProductCtrl
		expResp         string
		expCode         int
	}{
		"get all products successfully": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				output: []controllers.ProductOutput{
					{
						ID:           194,
						Name:         "iPhone 14",
						Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:        decimal.New(1500, 0),
						Quantity:     50,
						AuthorID:     136,
						CategoryName: "Smartphone",
						CreatedAt:    myCreatedTime,
						UpdatedAt:    myUpdatedTime,
					},
					{
						ID:           195,
						Name:         "Macbook",
						Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:        decimal.New(1500, 0),
						Quantity:     50,
						AuthorID:     136,
						CategoryName: "Smartphone",
						CreatedAt:    myCreatedTime,
						UpdatedAt:    myUpdatedTime,
					},
				},
			},
			expResp: `[{"id":194,"name":"iPhone 14","description":"An Apple smartphone with A15 chip, 6GB RAM and 512GB storage","price":"1500","quantity":50,"author_id":136,"category":"Smartphone","created_at":"2023-06-02T09:08:36.046843Z","updated_at":"2023-06-02T09:08:36.046843Z"},{"id":195,"name":"Macbook","description":"An Apple smartphone with A15 chip, 6GB RAM and 512GB storage","price":"1500","quantity":50,"author_id":136,"category":"Smartphone","created_at":"2023-06-02T09:08:36.046843Z","updated_at":"2023-06-02T09:08:36.046843Z"}]`,
			expCode: http.StatusOK,
		},
		"get all products with filter request successfully": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				query: controllers.ProductCtrlFilter{
					Name: "iPhone",
				},
				output: []controllers.ProductOutput{
					{
						ID:           194,
						Name:         "iPhone 14",
						Description:  "An Apple smartphone with A15 chip, 6GB RAM and 512GB storage",
						Price:        decimal.New(1500, 0),
						Quantity:     50,
						AuthorID:     136,
						CategoryName: "Smartphone",
						CreatedAt:    myCreatedTime,
						UpdatedAt:    myUpdatedTime,
					},
				},
			},
			givenFilter: "queryName=iPhone",
			expResp:     `[{"id":194,"name":"iPhone 14","description":"An Apple smartphone with A15 chip, 6GB RAM and 512GB storage","price":"1500","quantity":50,"author_id":136,"category":"Smartphone","created_at":"2023-06-02T09:08:36.046843Z","updated_at":"2023-06-02T09:08:36.046843Z"}]`,
			expCode:     http.StatusOK,
		},
		"products not found with filter": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				query: controllers.ProductCtrlFilter{
					Name: "imac",
				},
			},
			givenFilter: "queryName=imac",
			expResp:     `null`,
			expCode:     http.StatusOK,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			handler := NewHandler(mockController)

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products?%s", tc.givenFilter), nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tc.mockProductCtrl.expCall {
				mockController.On("GetProducts", context.Background(), tc.mockProductCtrl.query).Return(tc.mockProductCtrl.output, tc.mockProductCtrl.err)
			}
			handler.GetProducts(w, r)

			assert.Equal(t, tc.expResp, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expCode, w.Code)

			if tc.mockProductCtrl.expCall {
				mockController.AssertCalled(t, "GetProducts", context.Background(), tc.mockProductCtrl.query)
			}
		})
	}
}

func Test_ProductHandler_ImportProductsFromCSV(t *testing.T) {
	type mockProductCtrl struct {
		expCall          bool
		productFileInput multipart.File
		err              error
	}

	testCases := map[string]struct {
		csvData         string // payload provided from end-user
		mockProductCtrl mockProductCtrl
		expResp         string
		expCode         int
	}{
		"import products successfully": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
			},
			csvData: `Name,Description,Price,Quantity,AuthorID,Category
		Samsung S23,"A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1100,10,179,Smartphone
		Samsung Galaxy Watch 5,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",250,20,179,Smartwatch`,
			expResp: `{"success":true}`,
			expCode: http.StatusOK,
		},
		"import products with product category not found": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				err:     controllers.ErrProductCategoryNotFound,
			},
			csvData: `Name,Description,Price,Quantity,AuthorID,Category
Samsung S23,"A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1100,10,179,Smartphone
Samsung Galaxy Watch 5,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",250,20,179,smartwatchhhhhhhh`,
			expResp: `{"message":"product category not found"}`,
			expCode: http.StatusNotFound,
		},
		"import products with user not found": {
			mockProductCtrl: mockProductCtrl{
				expCall: true,
				err:     controllers.ErrUserNotFound,
			},
			csvData: `Name,Description,Price,Quantity,AuthorID,Category
Samsung S23,"A Samsung device with Snapdragon 8 Gen 2 chipset, 8GB RAM, 128GB storage.",1100,10,1799,Smartphone
Samsung Galaxy Watch 5,"A Samsung smartwatch with a Super Amoled 1.2 inch screen, 40mm clock face.",250,20,1799,Smartwatch`,
			expResp: `{"message":"user not found"}`,
			expCode: http.StatusNotFound,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			handler := NewHandler(mockController)

			// Create a new multipart writer and add the CSV data as a form file
			body := new(bytes.Buffer)
			multipartWriter := multipart.NewWriter(body)
			formFile, err := multipartWriter.CreateFormFile("products", "products.csv")
			assert.NoError(t, err)
			if _, err := formFile.Write([]byte(tc.csvData)); err != nil {
				assert.NoError(t, err)
			}

			// Set the content type header and close the multipart writer
			contentType := multipartWriter.FormDataContentType()
			if err := multipartWriter.Close(); err != nil {
				assert.NoError(t, err)
			}

			r := httptest.NewRequest(http.MethodPost, "/products/import-csv", body)
			r.Header.Add("Content-Type", contentType)
			w := httptest.NewRecorder()

			// get the file from form file
			file, _, err := r.FormFile("products")
			if err != nil {
				assert.EqualError(t, err, tc.mockProductCtrl.err.Error())
			}

			tc.mockProductCtrl.productFileInput = file

			if tc.mockProductCtrl.expCall {
				mockController.On("ImportProductsFromCSV", context.Background(), tc.mockProductCtrl.productFileInput).Return(tc.mockProductCtrl.err)
			}

			handler.ImportProductsFromCSV(w, r)

			assert.Equal(t, tc.expResp, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expCode, w.Code)

			if tc.mockProductCtrl.expCall {
				mockController.AssertCalled(t, "ImportProductsFromCSV", context.Background(), tc.mockProductCtrl.productFileInput)
			}
		})
	}
}

func Test_ProductHandler_ExportProductsFromCSV(t *testing.T) {
	type mockExportProductsCtrl struct {
		input controllers.ProductCtrlFilter
		err   error
	}

	type mockWriteEmailContentCtrl struct {
		emailToList []string
		err         error
	}

	testCases := map[string]struct {
		expCall                   bool
		givenFilter               string
		mockExportProductsCtrl    mockExportProductsCtrl
		mockWriteEmailContentCtrl mockWriteEmailContentCtrl
		csvData                   [][]string
		expResp                   string
		expCode                   int
		expErr                    error
	}{
		"export products successfully": {
			expCall: true,
			mockWriteEmailContentCtrl: mockWriteEmailContentCtrl{
				emailToList: []string{"qthuy2609@gmail.com"},
			},
			csvData: [][]string{
				{"ID", "Name", "Description", "Price", "Quantity", "AuthorID", "Category", "CreatedAt", "UpdatedAt"}, {"1051", "Samsung Galaxy Watch 5", "A Samsung smartwatch with a Super Amoled 1.2 inch screen", " 40mm clock face.", "250", "18", "1", "Smartwatch", "2023-06-14T08:38:13Z", "2023-06-26T02:37:17Z"},
				{"1050", "Samsung S23", "A Samsung device with Snapdragon 8 Gen 2 chipset", " 8GB RAM", "128GB storage.", "1100", "100", "1", "Smartphone", "2023-06-14T08:38:13Z", "2023-06-28T02:34:17Z"},
			},
			expResp: `{"success":true}`,
			expCode: http.StatusOK,
		},
		"export products successfully with filter": {
			expCall: true,
			mockExportProductsCtrl: mockExportProductsCtrl{
				input: controllers.ProductCtrlFilter{
					Name: "S23",
				},
			},
			mockWriteEmailContentCtrl: mockWriteEmailContentCtrl{
				emailToList: []string{"qthuy2609@gmail.com"},
			},
			givenFilter: "queryName=S23",
			expResp:     `{"success":true}`,
			expCode:     http.StatusOK,
		},
		"export products with bad date request filter": {
			expCall:     false,
			givenFilter: "date=2023-06-222",
			expResp:     `{"message":"invalid date format, dates must follow the format yyyy-mm-dd"}`,
			expCode:     http.StatusBadRequest,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			handler := NewHandler(mockController)

			emailUsername := os.Getenv("EMAIL_USERNAME")
			pipeReader, pipeWriter := io.Pipe()

			// write the products to the pipe
			go func() {
				csvWriter := csv.NewWriter(pipeWriter)
				for _, p := range tc.csvData {
					csvWriter.Write(p)
				}
				csvWriter.Flush()
				pipeWriter.Close()
			}()

			if tc.expCall {
				mockController.On("ExportProductsToCSV", mock.Anything, tc.mockExportProductsCtrl.input).Return(pipeReader, tc.mockExportProductsCtrl.err)
				mockController.On("WriteEmailContent", emailUsername, pipeReader, tc.mockWriteEmailContentCtrl.emailToList).Return(mock.Anything, tc.mockWriteEmailContentCtrl.err)
			}

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products/export-products?%s&emails=%s", tc.givenFilter, strings.Join(tc.mockWriteEmailContentCtrl.emailToList, ",")), nil)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			pipeReader.Close()
			handler.ExportProductsToCSV(w, r)

			assert.Equal(t, tc.expResp, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expCode, w.Code)

			if tc.expCall {
				mockController.AssertCalled(t, "ExportProductsToCSV", mock.Anything, tc.mockExportProductsCtrl.input)
			}
		})
	}
}
