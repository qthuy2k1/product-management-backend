package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/stretchr/testify/assert"
)

// Test CreateProductCategory in Handler layer
func Test_ProductCategoryHandler_CreateProductCategory(t *testing.T) {
	type mockPCateCtrl struct {
		expCall    bool
		pCateInput controllers.PCateInput
		err        error
	}
	testCases := map[string]struct {
		givenInput    string
		mockPCateCtrl mockPCateCtrl
		expResp       string
		expCode       int
	}{
		"create product category successfully": {
			givenInput: `{"name":"Cellphone","description":"Cellphone"}`,
			mockPCateCtrl: mockPCateCtrl{
				expCall: true,
				pCateInput: controllers.PCateInput{
					Name:        "Cellphone",
					Description: "Cellphone",
				},
			},
			expResp: `{"success":true}`,
			expCode: http.StatusCreated,
		},
		"create product category with missing name field": {
			givenInput: `{"description":"Cellphone"}`,
			mockPCateCtrl: mockPCateCtrl{
				expCall: false,
			},
			expResp: `{"message":"name cannot be blank"}`,
			expCode: http.StatusBadRequest,
		},
		"create product category with missing description field": {
			givenInput: `{"name":"Cellphone"}`,
			mockPCateCtrl: mockPCateCtrl{
				expCall: false,
			},
			expResp: `{"message":"description cannot be blank"}`,
			expCode: http.StatusBadRequest,
		},
		"create product category with invalid JSON": {
			givenInput: `{"name":"Cellphone","description":"Cellphone"`,
			mockPCateCtrl: mockPCateCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid json"}`,
			expCode: http.StatusBadRequest,
		},
		"create product category with wrong type field": {
			givenInput: `{"name":123,"description":"Cellphone"}`,
			mockPCateCtrl: mockPCateCtrl{
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

			r := httptest.NewRequest(http.MethodPost, "/product-categories", strings.NewReader(tc.givenInput))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tc.mockPCateCtrl.expCall {
				mockController.On("CreateProductCategory", context.Background(), tc.mockPCateCtrl.pCateInput).Return(tc.mockPCateCtrl.err)
			}

			handler.CreateProductCategory(w, r)

			assert.Equal(t, tc.expResp, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expCode, w.Code)

			if tc.mockPCateCtrl.expCall {
				mockController.AssertCalled(t, "CreateProductCategory", context.Background(), tc.mockPCateCtrl.pCateInput)
			}
		})
	}
}
