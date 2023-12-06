package rest

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/qthuy2k1/product-management/internal/utils"
	"net/http"
	"strings"
	"time"
)

type pCateRequest struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductCategory gets the product category data from form data, calls to CreateProductCategory controller and returns the status
func (h *Handler) CreateProductCategory(w http.ResponseWriter, r *http.Request) {
	pCateReq := pCateRequest{}
	ctx := r.Context()
	// Parse JSON request body into a Product struct
	if err := json.NewDecoder(r.Body).Decode(&pCateReq); err != nil {
		render.Render(w, r, ErrInvalidJson)
		return
	}

	// validate product category
	pCateInput, errResp := validateAndConvertProductCategory(pCateReq)
	if errResp != nil {
		render.Render(w, r, errResp)
		return
	}

	if err := h.Controller.CreateProductCategory(ctx, pCateInput); err != nil {
		render.Render(w, r, convertCtrlError(err))
		return
	}

	response := map[string]bool{"success": true}
	utils.RenderJson(w, response, http.StatusCreated)
}

// validateAndConvertProductCategory validates the product category from body request and returns product category struct in controller layer
func validateAndConvertProductCategory(pCateReq pCateRequest) (controllers.PCateInput, *ErrorResponse) {
	if len(strings.TrimSpace(pCateReq.Name)) == 0 {
		return controllers.PCateInput{}, ErrMissingName
	}

	if len(strings.TrimSpace(pCateReq.Name)) > 255 {
		return controllers.PCateInput{}, ErrNameTooLong
	}

	if len(strings.TrimSpace(pCateReq.Description)) == 0 {
		return controllers.PCateInput{}, ErrMissingDesc
	}

	// assign value to product struct in repository layer
	pCate := controllers.PCateInput{
		Name:        strings.TrimSpace(pCateReq.Name),
		Description: strings.TrimSpace(pCateReq.Description),
	}

	return pCate, nil
}

type OrderItemRequest struct {
	OrderID   int
	ProductID int
	Quantity  int
}
