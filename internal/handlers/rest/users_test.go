package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test CreateUser in Handler layer
func Test_UserHandler_CreateUser(t *testing.T) {
	type mockUserCtrl struct {
		expCall   bool
		userInput controllers.UserInput
		err       error
	}
	// Set up a test case
	testCases := map[string]struct {
		givenInput   string
		mockUserCtrl mockUserCtrl
		expResp      string
		expCode      int
	}{
		"create user successfully": {
			givenInput: `{"name":"John Doe","email":"johndoe@example.com","password":"password123"}`,
			mockUserCtrl: mockUserCtrl{
				expCall: true,
				userInput: controllers.UserInput{
					Name:     "John Doe",
					Email:    "johndoe@example.com",
					Password: "password123",
					Role:     "user",
					Status:   "activated",
				},
			},
			expResp: `{"success":true}`,
			expCode: http.StatusCreated,
		},
		"missing email field": {
			givenInput: `{"name":"John Doe","password":"password123"}`,
			mockUserCtrl: mockUserCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid email"}`,
			expCode: http.StatusBadRequest,
		},
		"invalid email address": {
			givenInput: `{"name":"John Doe","email":"invalid-email","password":"password123"}`,
			mockUserCtrl: mockUserCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid email"}`,
			expCode: http.StatusBadRequest,
		},
		"password is less than 6 characters": {
			givenInput: `{"name":"John Doe","email":"johndoe@example.com","password":"123"}`,
			mockUserCtrl: mockUserCtrl{
				expCall: false,
			},
			expResp: `{"message":"password must between 6 and 72 characters"}`,
			expCode: http.StatusBadRequest,
		},
	}

	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			handler := NewHandler(mockController)

			r := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tc.givenInput))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tc.mockUserCtrl.expCall {
				mockController.On("CreateUser", context.Background(), tc.mockUserCtrl.userInput).Return(tc.mockUserCtrl.err)
			}
			handler.CreateUser(w, r)

			assert.Equal(t, tc.expResp, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expCode, w.Code)

			if tc.mockUserCtrl.expCall {
				mockController.AssertCalled(t, "CreateUser", context.Background(), tc.mockUserCtrl.userInput)
			}
		})
	}
}

// Test GetUser in Handler layer
func Test_UserHandler_GetUser(t *testing.T) {
	myCreatedTime, err := time.Parse("2006-01-02T15:04:05.999999Z", "2023-05-11T09:01:53.102071Z")
	assert.NoError(t, err)
	myUpdatedTime, err := time.Parse("2006-01-02T15:04:05.999999Z", "2023-05-11T09:01:53.102071Z")
	assert.NoError(t, err)

	type mockUserCtrl struct {
		expCall bool
		userID  int
		output  controllers.UserOutput
		err     error
	}
	// Set up the test cases
	testCases := map[string]struct {
		userID       int
		mockUserCtrl mockUserCtrl
		expResp      string
		expCode      int
	}{
		"get user successfully": {
			userID: 1,
			mockUserCtrl: mockUserCtrl{
				expCall: true,
				userID:  1,
				output: controllers.UserOutput{
					ID:        1,
					Name:      "Thuy Nguyen",
					Email:     "qthuy@gmail.com",
					Password:  "$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC",
					Role:      "user",
					Status:    "activated",
					CreatedAt: myCreatedTime,
					UpdatedAt: myUpdatedTime,
				},
			},
			expResp: `{"id":1,"name":"Thuy Nguyen","email":"qthuy@gmail.com","password":"$2a$14$1VRhclX4P/JkiqJ7nKXYvuC/4NdvKXreBPay9sDJv1CpWs4eFhXAC","role":"user","status":"activated","created_at":"2023-05-11T09:01:53.102071Z","updated_at":"2023-05-11T09:01:53.102071Z"}`,
			expCode: http.StatusOK,
		},
		"invalid user ID": {
			userID: 0,
			mockUserCtrl: mockUserCtrl{
				expCall: false,
			},
			expResp: `{"message":"invalid user ID"}`,
			expCode: http.StatusBadRequest,
		},
	}

	// Test each test case
	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			mockController := &controllers.MockIController{}
			handler := NewHandler(mockController)
			if tc.mockUserCtrl.expCall {
				mockController.On("GetUser", mock.Anything, tc.userID).Return(tc.mockUserCtrl.output, tc.mockUserCtrl.err)
			}
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", tc.userID), nil)
			w := httptest.NewRecorder()
			// New route
			rctx := chi.NewRouteContext()
			// Add user id to url params
			rctx.URLParams.Add("userID", strconv.Itoa(tc.userID))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler.GetUser(w, r)

			assert.Equal(t, tc.expResp, strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expCode, w.Code)

			if tc.mockUserCtrl.expCall {
				mockController.AssertCalled(t, "GetUser", mock.Anything, tc.userID)
			}
		})
	}
}
