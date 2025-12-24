package mocks

import (
	"context"

	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(ctx context.Context, user models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) LoginUser(ctx context.Context, email, password string) (*models.UserResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) LogoutUser(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockUserService) RefreshToken(ctx context.Context, userId string) (string, string, error) {
	args := m.Called(ctx, userId)
	return args.String(0), args.String(1), args.Error(2)
}

type MockMovieService struct {
	mock.Mock
}

func (m *MockMovieService) GetMovies(ctx context.Context) ([]models.Movie, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Movie), args.Error(1)
}

func (m *MockMovieService) GetMovie(ctx context.Context, imdbID string) (*models.Movie, error) {
	args := m.Called(ctx, imdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Movie), args.Error(1)
}

func (m *MockMovieService) AddMovie(ctx context.Context, movie models.Movie) error {
	args := m.Called(ctx, movie)
	return args.Error(0)
}

func (m *MockMovieService) UpdateAdminReview(ctx context.Context, imdbID string, review string) (string, string, error) {
	args := m.Called(ctx, imdbID, review)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockMovieService) GetRecommendedMovies(ctx context.Context, userId string) ([]models.Movie, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Movie), args.Error(1)
}

func (m *MockMovieService) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Genre), args.Error(1)
}
