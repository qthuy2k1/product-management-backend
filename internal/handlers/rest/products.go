package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/shopspring/decimal"

	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/qthuy2k1/product-management/internal/utils"
)

type productRequest struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Price        decimal.Decimal `json:"price"`
	Quantity     int             `json:"quantity"`
	AuthorID     int             `json:"author_id"`
	CategoryName string          `json:"category"`
}

// CreateProduct gets the product data from body request, calls to CreateProduct controller and returns the status
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	pReq := productRequest{}
	ctx := r.Context()
	// Parse JSON request body into a Product struct
	if err := json.NewDecoder(r.Body).Decode(&pReq); err != nil {
		render.Render(w, r, ErrInvalidJson)
		return
	}

	product, errResp := validateAndConvertProduct(pReq)
	if errResp != nil {
		render.Render(w, r, errResp)
		return
	}

	if err := h.Controller.CreateProduct(ctx, product); err != nil {
		log.Println(err)
		render.Render(w, r, convertCtrlError(err))
		return
	}

	response := map[string]bool{"success": true}
	utils.RenderJson(w, response, http.StatusCreated)
}

// UpdateProduct gets the product data from body request, calls to UpdateProduct controller and returns the status
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil || id <= 0 {
		render.Render(w, r, ErrInvalidProductID)
		return
	}

	pReq := productRequest{}
	if err := json.NewDecoder(r.Body).Decode(&pReq); err != nil {
		render.Render(w, r, ErrInvalidJson)
		return
	}

	product, errResp := validateAndConvertProduct(pReq)
	if errResp != nil {
		render.Render(w, r, errResp)
		return
	}

	product.ID = id

	if err = h.Controller.UpdateProduct(ctx, product); err != nil {
		render.Render(w, r, convertCtrlError(err))
		return
	}

	response := map[string]bool{"success": true}
	utils.RenderJson(w, response, http.StatusOK)
}

// DeleteProduct retrieves an id and call to DeleteProduct in Controller layer and returns the status
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "productID"))
	if err != nil || id <= 0 {
		render.Render(w, r, ErrInvalidProductID)
		return
	}

	if err = h.Controller.DeleteProduct(ctx, id); err != nil {
		render.Render(w, r, convertCtrlError(err))
		return
	}

	response := map[string]bool{"success": true}
	utils.RenderJson(w, response, http.StatusOK)
}

type ProductFilterReq struct {
	Name  string
	Date  string
	Email []string
}

type ProductResponse struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Price        decimal.Decimal `json:"price"`
	Quantity     int             `json:"quantity"`
	AuthorID     int             `json:"author_id"`
	CategoryName string          `json:"category"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// GetProducts retrieves all the products in db
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	pCtrlFilter, errResp := validateAndConvertProductFilter(query)
	if errResp != nil {
		render.Render(w, r, errResp)
		return
	}

	products, err := h.Controller.GetProducts(ctx, pCtrlFilter)
	if err != nil {
		render.Render(w, r, convertCtrlError(err))
		return
	}

	var pResp []ProductResponse
	for _, p := range products {
		pResp = append(pResp, ProductResponse{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			Quantity:     p.Quantity,
			AuthorID:     p.AuthorID,
			CategoryName: p.CategoryName,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		})
	}

	utils.RenderJson(w, pResp, http.StatusOK)
}

// ImportProductsFromCSV imports list of products data from a CSV file
func (h *Handler) ImportProductsFromCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get the file that was selected from form file input
	file, fileHeader, err := r.FormFile("products")
	if err != nil {
		render.Render(w, r, convertCtrlError(err))
		return
	}
	defer file.Close()

	// check if file is not a csv file
	if filepath.Ext(fileHeader.Filename) != ".csv" {
		render.Render(w, r, ErrInvalidCSVFileType)
		return
	}

	if err = h.Controller.ImportProductsFromCSV(ctx, file); err != nil {
		log.Println(err)
		render.Render(w, r, convertCtrlError(err))
		return
	}

	response := map[string]bool{"success": true}
	utils.RenderJson(w, response, http.StatusOK)
}

// ExportProductsToCSV exports all the products in db to a csv file and send an email to user
func (h *Handler) ExportProductsToCSV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pCtrlFilter, errResp := validateAndConvertProductFilter(r.URL.Query())
	if errResp != nil {
		render.Render(w, r, errResp)
		return
	}

	reader, err := h.Controller.ExportProductsToCSV(ctx, pCtrlFilter)
	if err != nil {
		log.Println(err)
		render.Render(w, r, convertCtrlError(err))
		return
	}

	if len(pCtrlFilter.Email) != 0 {
		if err := h.Controller.SendEmailProduct(pCtrlFilter.Email, reader); err != nil {
			log.Println(err)
			render.Render(w, r, convertCtrlError(err))
			return
		}
	} else {
		// Set HTTP headers
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="export_products.csv"`)

		// copy data from pipe reader to HTTP response writer
		if _, err := io.Copy(w, reader); err != nil {
			log.Println(err)
			render.Render(w, r, convertCtrlError(err))
			return
		}
	}

	response := map[string]bool{"success": true}
	utils.RenderJson(w, response, http.StatusOK)
}

// validateAndConvertProduct validates the product from body request and return product struct in controller layer
func validateAndConvertProduct(pReq productRequest) (controllers.ProductInput, *ErrorResponse) {
	if len(strings.TrimSpace(pReq.Name)) == 0 {
		return controllers.ProductInput{}, ErrMissingName
	}

	if len(strings.TrimSpace(pReq.Name)) > 255 {
		return controllers.ProductInput{}, ErrNameTooLong
	}

	if len(strings.TrimSpace(pReq.Description)) == 0 {
		return controllers.ProductInput{}, ErrMissingDesc
	}

	if pReq.Price.IsZero() && pReq.Price.IsNegative() && len(pReq.Price.String()) > 15 {
		return controllers.ProductInput{}, ErrInvalidPrice
	}

	if pReq.Quantity < 0 {
		return controllers.ProductInput{}, ErrInvalidQuantity
	}

	if pReq.AuthorID <= 0 {
		return controllers.ProductInput{}, ErrInvalidAuthorID
	}

	if len(strings.TrimSpace(pReq.CategoryName)) == 0 {
		return controllers.ProductInput{}, ErrMissingCategoryName
	}

	return controllers.ProductInput{
		Name:         strings.TrimSpace(pReq.Name),
		Description:  strings.TrimSpace(pReq.Description),
		Price:        pReq.Price,
		Quantity:     pReq.Quantity,
		AuthorID:     pReq.AuthorID,
		CategoryName: strings.TrimSpace(pReq.CategoryName),
	}, nil
}

func validateAndConvertProductFilter(query url.Values) (controllers.ProductCtrlFilter, *ErrorResponse) {
	var filterReq ProductFilterReq
	nameFilter := strings.TrimSpace(query.Get("queryName"))
	if nameFilter != "" {
		filterReq.Name = nameFilter
	}

	dateFilter := strings.TrimSpace(query.Get("date"))
	if dateFilter != "" {
		if _, err := time.Parse("2006-01-02", dateFilter); err != nil {
			return controllers.ProductCtrlFilter{}, ErrDateBadRequest
		}
		filterReq.Date = dateFilter
	}

	emails := strings.TrimSpace(query.Get("emails"))
	if emails != "" {
		emailToList := strings.Split(emails, ",")
		for i := range emailToList {
			emailToList[i] = strings.TrimSpace(emailToList[i])
			// remove the email if not valid
			if !isValidEmail(emailToList[i]) {
				log.Println("email invalid:", emailToList[i])
				emailToList = append(emailToList[:i], emailToList[i+1:]...)
			}
		}

		if len(emailToList) == 0 {
			return controllers.ProductCtrlFilter{}, ErrInvalidEmail
		}
		filterReq.Email = emailToList
	}

	pCtrlFilter := controllers.ProductCtrlFilter{
		Name:  filterReq.Name,
		Date:  filterReq.Date,
		Email: filterReq.Email,
	}
	return pCtrlFilter, nil
}
