package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/qthuy2k1/product-management/internal/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/qthuy2k1/product-management/internal/utils"
)

type userRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

// CreateUser gets the user data from body request, calls to CreateUser controller and returns the status
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userReq := userRequest{}
	ctx := r.Context()
	// Parse JSON request body into a User struct
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		render.Render(w, r, ErrInvalidJson)
		return
	}

	userInput, errResp := validateAndConvertUser(userReq)
	if errResp != nil {
		render.Render(w, r, errResp)
		return
	}

	// Set default value for new user when creating
	userInput.Role = "user"
	userInput.Status = "activated"

	// Create user
	if err := h.Controller.CreateUser(ctx, userInput); err != nil {
		log.Println(err)
		render.Render(w, r, convertCtrlError(err))
		return
	}

	// render json to notify success
	response := map[string]bool{"success": true}
	utils.RenderJson(w, response, http.StatusCreated)
}

type userResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUser receives the user id from url param, calls to Controller and retrieves the User repository
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil || id <= 0 {
		render.Render(w, r, ErrInvalidUserID)
		return
	}

	user, err := h.Controller.GetUser(ctx, id)
	if err != nil {
		render.Render(w, r, convertCtrlError(err))
		return
	}

	// assign value to user struct in handler layer
	userResp := userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	utils.RenderJson(w, userResp, http.StatusOK)
}

// validateAndConvertUser validates the user from body request and returns user struct in controller layer
func validateAndConvertUser(userReq userRequest) (controllers.UserInput, *ErrorResponse) {
	// Validate the name
	if len(strings.TrimSpace(userReq.Name)) == 0 {
		return controllers.UserInput{}, ErrMissingName
	}

	// Validate the email
	if !isValidEmail(userReq.Email) {
		return controllers.UserInput{}, ErrInvalidEmail
	}
	// Validate the password
	// The password must contain at least 6 characters
	if !isValidPassword(userReq.Password) {
		return controllers.UserInput{}, ErrInvalidPassword
	}

	userInput := controllers.UserInput{
		Name:     strings.TrimSpace(userReq.Name),
		Email:    strings.TrimSpace(userReq.Email),
		Password: userReq.Password,
	}

	return userInput, nil
}

// isValidEmail validates that an email address is in a valid format
func isValidEmail(email string) bool {
	if len(strings.TrimSpace(email)) == 0 {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

// isValidPassword validates that a password meets the minimum requirements
func isValidPassword(password string) bool {
	if password == "" {
		return false
	}
	// Check if the password must be between 6 and 72 characters
	return len(password) >= 6 && len(password) <= 72
}
