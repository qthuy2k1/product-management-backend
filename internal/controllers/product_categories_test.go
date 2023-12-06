package controllers

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/qthuy2k1/product-management/internal/models"

	"github.com/qthuy2k1/product-management/internal/repositories"
	"github.com/stretchr/testify/assert"
)

// Test CreateProductCategory in Controller layer
func Test_ProductCategoryController_CreateProductCategory(t *testing.T) {
	type mockPCateRepo struct {
		expCall    bool
		pCateInput repositories.ProductCategory
		err        error
	}
	tests := map[string]struct {
		mockPCateRepo mockPCateRepo
		pCateInput    PCateInput
		err           error
	}{
		"success": {
			mockPCateRepo: mockPCateRepo{
				expCall: true,
				pCateInput: repositories.ProductCategory{
					Name:        "Cellphone",
					Description: "Cellphone",
				},
			},
			pCateInput: PCateInput{
				Name:        "Cellphone",
				Description: "Cellphone",
			},
		},
		"name too long": {
			mockPCateRepo: mockPCateRepo{
				expCall: true,
				pCateInput: repositories.ProductCategory{
					Name:        "cellphone" + strings.Repeat("a", 251),
					Description: "Cellphone",
				},
				err: errors.New("name too long"),
			},
			pCateInput: PCateInput{
				Name:        "cellphone" + strings.Repeat("a", 251),
				Description: "Cellphone",
			},
			err: errors.New("name too long"),
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.NewMockIRepository(t)
			controller := NewController(mockRepo)
			if tc.mockPCateRepo.expCall {
				mockRepo.On("CreateProductCategory", context.Background(), tc.mockPCateRepo.pCateInput).Return(tc.mockPCateRepo.err)
			}

			if err := controller.CreateProductCategory(context.Background(), tc.pCateInput); tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}
func Test_ProductCategoryController_GetProductCategoryByName(t *testing.T) {
	type mockPCateRepo struct {
		expCall   bool
		pCateName string
		output    models.ProductCategory
		err       error
	}
	tests := map[string]struct {
		mockPCateRepo mockPCateRepo
		pCateName     string
		output        PCateOutput
		err           error
	}{
		"success": {
			mockPCateRepo: mockPCateRepo{
				expCall:   true,
				pCateName: "Smartphone",
				output: models.ProductCategory{
					ID:          2,
					Name:        "Smartphone",
					Description: "A mobile device",
				},
			},
			pCateName: "Smartphone",
			output: PCateOutput{
				ID:          2,
				Name:        "Smartphone",
				Description: "A mobile device",
			},
		},
		"error when not valid user id": {
			mockPCateRepo: mockPCateRepo{
				expCall:   true,
				pCateName: "abc",
				err:       repositories.ErrProductCategoryNotFound,
			},
			pCateName: "abc",
			err:       ErrProductCategoryNotFound,
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.MockIRepository{}
			controller := NewController(&mockRepo)
			if tc.mockPCateRepo.expCall {
				mockRepo.On("GetProductCategoryByName", context.Background(), tc.mockPCateRepo.pCateName).Return(tc.mockPCateRepo.output, tc.mockPCateRepo.err)
			}
			user, err := controller.GetProductCategoryByName(context.Background(), tc.pCateName)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.Equal(t, user, tc.output)
			}
		})
	}
}
