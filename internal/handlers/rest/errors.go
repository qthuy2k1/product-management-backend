package rest

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/qthuy2k1/product-management/internal/controllers"
)

type ErrorResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

var (
	ErrInvalidProductID        = &ErrorResponse{StatusCode: 400, Message: "invalid product ID"}
	ErrInvalidUserID           = &ErrorResponse{StatusCode: 400, Message: "invalid user ID"}
	ErrMissingName             = &ErrorResponse{StatusCode: 400, Message: "name cannot be blank"}
	ErrNameTooLong             = &ErrorResponse{StatusCode: 400, Message: "name too long"}
	ErrMissingDesc             = &ErrorResponse{StatusCode: 400, Message: "description cannot be blank"}
	ErrMissingCategoryName     = &ErrorResponse{StatusCode: 400, Message: "category cannot be blank"}
	ErrInvalidPrice            = &ErrorResponse{StatusCode: 400, Message: "price must be greater than 0 and less than 15 digits"}
	ErrInvalidQuantity         = &ErrorResponse{StatusCode: 400, Message: "quantity must be non-negative"}
	ErrInvalidEmail            = &ErrorResponse{StatusCode: 400, Message: "invalid email"}
	ErrInvalidAuthorID         = &ErrorResponse{StatusCode: 400, Message: "invalid author id"}
	ErrInvalidPassword         = &ErrorResponse{StatusCode: 400, Message: "password must between 6 and 72 characters"}
	ErrInvalidJson             = &ErrorResponse{StatusCode: 400, Message: "invalid json"}
	ErrMethodNotAllowed        = &ErrorResponse{StatusCode: 405, Message: "method not allowed"}
	ErrNotFound                = &ErrorResponse{StatusCode: 404, Message: "not found"}
	ErrProductNotFound         = &ErrorResponse{StatusCode: 404, Message: "product not found"}
	ErrProductCategoryNotFound = &ErrorResponse{StatusCode: 404, Message: "product category not found"}
	ErrUserNotFound            = &ErrorResponse{StatusCode: 404, Message: "user not found"}
	ErrDateBadRequest          = &ErrorResponse{StatusCode: 400, Message: "invalid date format, dates must follow the format yyyy-mm-dd"}
	ErrInvalidCSVFileType      = &ErrorResponse{StatusCode: 400, Message: "file must be a csv"}
	ErrNotEnoughColumns        = &ErrorResponse{StatusCode: 400, Message: "not enough columns"}
	ErrIncorrectColumnNames    = &ErrorResponse{StatusCode: 400, Message: "incorrect column names, the file must have columns named: Name, Description, Price, Quantity, AuthorID, Category"}
	ErrCSVFileFormat           = &ErrorResponse{StatusCode: 400, Message: "the file is not a valid CSV file"}
	ErrProductListEmpty        = &ErrorResponse{StatusCode: 200, Message: "the list of products is empty"}
)

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

// ServerErrorRenderer returns the error response with status code 500 and the "Internal server error" message
func ServerErrorRenderer() *ErrorResponse {
	return &ErrorResponse{
		StatusCode: 500,
		Message:    "internal server error",
	}
}

// ConvertCtrlError compares the error return with the error in controller and returns the corresponding ErrorResponse
func convertCtrlError(err error) *ErrorResponse {
	switch err {
	case controllers.ErrUserNotFound:
		return ErrUserNotFound
	case controllers.ErrProductCategoryNotFound:
		return ErrProductCategoryNotFound
	case controllers.ErrProductNotFound:
		return ErrProductNotFound
	case controllers.ErrNotEnoughColumns:
		return ErrNotEnoughColumns
	case controllers.ErrIncorrectColumnNames:
		return ErrIncorrectColumnNames
	case controllers.ErrCSVFileFormat:
		return ErrCSVFileFormat
	default:
		log.Println(err)
		return ServerErrorRenderer()
	}
}
