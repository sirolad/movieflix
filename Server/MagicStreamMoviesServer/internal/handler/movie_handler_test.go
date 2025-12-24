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
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestGetMovies(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockMovieService)
		movieHandler := NewMovieHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/movies", nil)
		c.Request = req

		movies := []models.Movie{
			{Title: "Movie 1", ImdbID: "tt1"},
			{Title: "Movie 2", ImdbID: "tt2"},
		}

		mockService.On("GetMovies", mock.Anything).Return(movies, nil)

		movieHandler.GetMovies(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockService := new(mocks.MockMovieService)
		movieHandler := NewMovieHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/movies", nil)
		c.Request = req

		mockService.On("GetMovies", mock.Anything).Return(nil, errors.New("db error"))

		movieHandler.GetMovies(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetMovie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockMovieService)
		movieHandler := NewMovieHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		// Setup param
		c.Params = []gin.Param{{Key: "imdb_id", Value: "tt123"}}
		req := httptest.NewRequest("GET", "/movie/tt123", nil)
		c.Request = req

		movie := models.Movie{Title: "Movie 1", ImdbID: "tt123"}

		mockService.On("GetMovie", mock.Anything, "tt123").Return(&movie, nil)

		movieHandler.GetMovie(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService := new(mocks.MockMovieService)
		movieHandler := NewMovieHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		c.Params = []gin.Param{{Key: "imdb_id", Value: "tt123"}}
		req := httptest.NewRequest("GET", "/movie/tt123", nil)
		c.Request = req

		mockService.On("GetMovie", mock.Anything, "tt123").Return(nil, mongo.ErrNoDocuments)

		movieHandler.GetMovie(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestAddMovie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mocks.MockMovieService)
		movieHandler := NewMovieHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		movie := models.Movie{
			Title:      "Test Movie",
			ImdbID:     "tt1234567",
			PosterPath: "http://example.com/poster.jpg",
			YouTubeID:  "dQw4w9WgXcQ",
			Genre: []models.Genre{
				{GenreID: 1, GenreName: "Action"},
			},
			Ranking: models.Ranking{
				RankingValue: 8,
				RankingName:  "Good",
			},
		}

		mockService.On("AddMovie", mock.Anything, mock.MatchedBy(func(m models.Movie) bool {
			return m.ImdbID == movie.ImdbID
		})).Return(nil)

		jsonBytes, _ := json.Marshal(movie)
		req := httptest.NewRequest("POST", "/movie", bytes.NewBuffer(jsonBytes))
		c.Request = req

		movieHandler.AddMovie(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		mockService := new(mocks.MockMovieService)
		movieHandler := NewMovieHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		movie := models.Movie{
			Title: "Incomplete Movie",
		}

		jsonBytes, _ := json.Marshal(movie)
		req := httptest.NewRequest("POST", "/movie", bytes.NewBuffer(jsonBytes))
		c.Request = req

		movieHandler.AddMovie(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "AddMovie")
	})
}
