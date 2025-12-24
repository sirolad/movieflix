package service_test

import (
	"context"
	"testing"

	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/config"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/mocks"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}
	svc := service.NewUserService(mockRepo, cfg)

	user := models.User{
		Email:    "test@example.com",
		Password: "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "USER",
	}

	// Expect CountUsersByEmail to return 0 (user does not exist)
	mockRepo.On("CountUsersByEmail", mock.Anything, user.Email).Return(int64(0), nil)

	// Expect CreateUser to be called
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u models.User) bool {
		return u.Email == user.Email && u.FirstName == user.FirstName
	})).Return(&mongo.InsertOneResult{}, nil)

	createdUser, err := svc.RegisterUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.NotEmpty(t, createdUser.Password) // Password should be hashed
	assert.NotEqual(t, "password123", createdUser.Password)
	assert.NotEmpty(t, createdUser.UserID)

	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UserExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := &config.Config{}
	svc := service.NewUserService(mockRepo, cfg)

	user := models.User{
		Email:    "existing@example.com",
		Password: "password123",
	}

	// Expect CountUsersByEmail to return 1 (user exists)
	mockRepo.On("CountUsersByEmail", mock.Anything, user.Email).Return(int64(1), nil)

	createdUser, err := svc.RegisterUser(context.Background(), user)

	assert.Error(t, err)
	assert.Nil(t, createdUser)
	assert.Equal(t, "user already exists", err.Error())

	// CreateUser should NOT be called
	mockRepo.AssertNotCalled(t, "CreateUser")

	mockRepo.AssertExpectations(t)
}
