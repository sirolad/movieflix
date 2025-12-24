package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/middleware"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/service"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MovieHandler struct {
	service  service.MovieService
	validate *validator.Validate
}

func NewMovieHandler(s service.MovieService) *MovieHandler {
	return &MovieHandler{
		service:  s,
		validate: validator.New(),
	}
}

// GetMovies godoc
// @Summary      Get all movies
// @Description  Get a list of all movies
// @Tags         movies
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Movie
// @Failure      500  {object}  map[string]interface{}
// @Router       /movies [get]
func (h *MovieHandler) GetMovies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	movies, err := h.service.GetMovies(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetAllGenres godoc
// @Summary      Get all genres
// @Description  Get a list of all unique genres
// @Tags         movies
// @Produce      json
// @Success      200  {array}   models.Genre
// @Failure      500  {object}  map[string]interface{}
// @Router       /genres [get]
func (h *MovieHandler) GetAllGenres(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	genres, err := h.service.GetAllGenres(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching genres"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// GetMovie godoc
// @Summary      Get a movie by ID
// @Description  Get details of a specific movie
// @Tags         movies
// @Produce      json
// @Security     BearerAuth
// @Param        imdb_id  path      string  true  "IMDB ID"
// @Success      200      {object}  models.Movie
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /movie/{imdb_id} [get]
func (h *MovieHandler) GetMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	movieID := c.Param("imdb_id")
	if movieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
		return
	}

	movie, err := h.service.GetMovie(ctx, movieID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movie"})
		}
		return
	}

	c.JSON(http.StatusOK, movie)
}

// AddMovie godoc
// @Summary      Add a new movie
// @Description  Add a new movie to the database
// @Tags         movies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        movie  body      models.Movie  true  "Movie Data"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /movie [post]
func (h *MovieHandler) AddMovie(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.validate.Struct(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	err := h.service.AddMovie(ctx, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding movie"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Movie added successfully"})
}

// UpdateAdminReview godoc
// @Summary      Update admin review (Admin only)
// @Description  Update the admin review for a movie (requires ADMIN role)
// @Tags         movies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        imdb_id  path      string  true  "IMDB ID"
// @Param        review   body      map[string]string  true  "Admin Review JSON {\"admin_review\": \"review\"}"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /movie/{imdb_id}/review [patch]
func (h *MovieHandler) UpdateAdminReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	role, err := middleware.GetRoleFromContext(c)
	if err != nil || role != "ADMIN" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admin access required"})
		return
	}

	movieID := c.Param("imdb_id")
	if movieID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
		return
	}

	var req struct {
		AdminReview string `json:"admin_review" validate:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	sentiment, _, err := h.service.UpdateAdminReview(ctx, movieID, req.AdminReview)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ranking_name": sentiment,
		"admin_review": req.AdminReview,
	})
}

// GetRecommendedMovies godoc
// @Summary      Get recommended movies
// @Description  Get recommended movies based on user's favorite genres
// @Tags         movies
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Movie
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /recommendedMovies [get]
func (h *MovieHandler) GetRecommendedMovies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	userId, err := middleware.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: " + err.Error()})
		return
	}

	movies, err := h.service.GetRecommendedMovies(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}
