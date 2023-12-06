package controllers

import (
	"context"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/qthuy2k1/product-management/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Name     string
	Email    string
	Password string
	Role     string
	Status   string
}

// CreateUser adds an user to database
func (c *Controller) CreateUser(ctx context.Context, user UserInput) error {
	// Sanitize and hash password
	password := html.EscapeString(strings.TrimSpace(user.Password))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	userInput := repositories.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPassword),
	}

	return c.Repository.CreateUser(ctx, userInput)
}

type UserOutput struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Role      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetUser gets an user from database by ID
func (c *Controller) GetUser(ctx context.Context, id int) (UserOutput, error) {
	user, err := c.Repository.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return UserOutput{}, ErrUserNotFound
		}
		return UserOutput{}, err
	}
	userResponse := UserOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return userResponse, nil
}
