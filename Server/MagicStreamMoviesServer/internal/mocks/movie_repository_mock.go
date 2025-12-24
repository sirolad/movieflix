package mocks

import (
	"context"

	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MockMovieRepository struct {
	mock.Mock
}

func (m *MockMovieRepository) GetMovies(ctx context.Context) ([]models.Movie, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Movie), args.Error(1)
}

func (m *MockMovieRepository) GetMovie(ctx context.Context, imdbID string) (*models.Movie, error) {
	args := m.Called(ctx, imdbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Movie), args.Error(1)
}

func (m *MockMovieRepository) CreateMovie(ctx context.Context, movie models.Movie) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, movie)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMovieRepository) UpdateMovieReview(ctx context.Context, imdbID string, review string, sentiment string, rankVal int) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, imdbID, review, sentiment, rankVal)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockMovieRepository) GetRankings(ctx context.Context) ([]models.Ranking, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Ranking), args.Error(1)
}

func (m *MockMovieRepository) GetRecommendedMovies(ctx context.Context, genres []string, limit int64) ([]models.Movie, error) {
	args := m.Called(ctx, genres, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Movie), args.Error(1)
}

func (m *MockMovieRepository) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Genre), args.Error(1)
}
