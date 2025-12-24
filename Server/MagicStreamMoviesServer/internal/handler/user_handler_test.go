package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/mocks"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockUserService)
		userHandler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		user := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  "password123",
			Role:      "USER",
		}

		createdUser := user
		createdUser.UserID = "new-id"

		mockService.On("RegisterUser", mock.Anything, mock.MatchedBy(func(u models.User) bool {
			return u.Email == user.Email
		})).Return(&createdUser, nil)

		jsonBytes, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBytes))
		c.Request = req

		userHandler.RegisterUser(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		mockService := new(mocks.MockUserService)
		userHandler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		user := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  "password123",
			Role:      "USER",
		}

		mockService.On("RegisterUser", mock.Anything, mock.Anything).Return(nil, errors.New("user already exists"))

		jsonBytes, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBytes))
		c.Request = req

		userHandler.RegisterUser(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		mockService := new(mocks.MockUserService)
		userHandler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		user := models.User{
			FirstName: "", // Invalid
			Email:     "invalid-email",
		}

		jsonBytes, _ := json.Marshal(user)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBytes))
		c.Request = req

		userHandler.RegisterUser(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "RegisterUser")
	})
}
