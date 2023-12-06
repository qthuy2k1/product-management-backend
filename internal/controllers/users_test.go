package controllers

import (
	"context"
	"errors"
	"testing"

	"github.com/qthuy2k1/product-management/internal/models"
	"github.com/qthuy2k1/product-management/internal/repositories"
	"github.com/stretchr/testify/assert"
)

// Test CreateUser in controller layer
func Test_UserController_CreateUser(t *testing.T) {
	type mockUserRepo struct {
		expCall bool
		input   repositories.User
		err     error
	}
	tests := map[string]struct {
		mockUserRepo mockUserRepo
		input        UserInput
		err          error
	}{
		"success": {
			mockUserRepo: mockUserRepo{
				expCall: true,
				input: repositories.User{
					Name:     "John Doe",
					Email:    "doe@gmail.com",
					Password: "password",
					Role:     "user",
					Status:   "activated",
				},
			},
			input: UserInput{
				Name:     "John Doe",
				Email:    "doe@gmail.com",
				Password: "password",
				Role:     "user",
				Status:   "activated",
			},
			err: nil,
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.MockIRepository{}
			controller := NewController(&mockRepo)
			// password := html.EscapeString(strings.TrimSpace(tc.input.Password))
			// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
			// assert.NoError(t, err)
			// tc.mockUserRepo.input.Password = string(hashedPassword)
			if tc.mockUserRepo.expCall {
				mockRepo.On("CreateUser", context.Background(), tc.mockUserRepo.input).Return(tc.mockUserRepo.err)
			}
			err := controller.CreateUser(context.Background(), tc.input)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			}
		})
	}
}

// Test GetUser in controller layer
func Test_UserController_GetUser(t *testing.T) {
	type mockUserRepo struct {
		expCall bool
		input   int
		output  models.User
		err     error
	}
	tests := map[string]struct {
		mockUserRepo mockUserRepo
		input        int
		output       UserOutput
		err          error
	}{
		"success": {
			mockUserRepo: mockUserRepo{
				expCall: true,
				input:   2,
				output: models.User{
					ID:       2,
					Name:     "John Doe",
					Email:    "doe@gmail.com",
					Password: "$2a$14$7WNpJLl4Oc8OysHpAO9G.e5LnRo.XYb1BMFXIUKQr2sX8s8NOuQGy",
					Role:     "user",
					Status:   "activated",
				},
			},
			input: 2,
			output: UserOutput{
				ID:       2,
				Name:     "John Doe",
				Email:    "doe@gmail.com",
				Password: "$2a$14$7WNpJLl4Oc8OysHpAO9G.e5LnRo.XYb1BMFXIUKQr2sX8s8NOuQGy",
				Role:     "user",
				Status:   "activated",
			},
		},
		"error when not valid user id": {
			mockUserRepo: mockUserRepo{
				expCall: true,
				input:   -1,
				err:     errors.New("invalid ID"),
			},
			input: -1,
			err:   errors.New("invalid ID"),
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			mockRepo := repositories.MockIRepository{}
			controller := NewController(&mockRepo)
			if tc.mockUserRepo.expCall {
				mockRepo.On("GetUser", context.Background(), tc.mockUserRepo.input).Return(tc.mockUserRepo.output, tc.mockUserRepo.err)
			}
			user, err := controller.GetUser(context.Background(), tc.input)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.Equal(t, user, tc.output)
			}
		})
	}
}
