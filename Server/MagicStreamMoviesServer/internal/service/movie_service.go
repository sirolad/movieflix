package service

import (
	"context"
	"errors"
	"strings"

	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/config"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/repository"
	"github.com/tmc/langchaingo/llms/openai"
)

type MovieService interface {
	GetMovies(ctx context.Context) ([]models.Movie, error)
	GetMovie(ctx context.Context, imdbID string) (*models.Movie, error)
	AddMovie(ctx context.Context, movie models.Movie) error
	UpdateAdminReview(ctx context.Context, imdbID string, review string) (string, string, error)
	GetRecommendedMovies(ctx context.Context, userId string) ([]models.Movie, error)
	GetAllGenres(ctx context.Context) ([]models.Genre, error)
}

type movieService struct {
	movieRepo repository.MovieRepository
	userRepo  repository.UserRepository
	config    *config.Config
}

func NewMovieService(movieRepo repository.MovieRepository, userRepo repository.UserRepository, cfg *config.Config) MovieService {
	return &movieService{
		movieRepo: movieRepo,
		userRepo:  userRepo,
		config:    cfg,
	}
}

func (s *movieService) GetMovies(ctx context.Context) ([]models.Movie, error) {
	return s.movieRepo.GetMovies(ctx)
}

func (s *movieService) GetMovie(ctx context.Context, imdbID string) (*models.Movie, error) {
	return s.movieRepo.GetMovie(ctx, imdbID)
}

func (s *movieService) AddMovie(ctx context.Context, movie models.Movie) error {
	_, err := s.movieRepo.CreateMovie(ctx, movie)
	return err
}

func (s *movieService) UpdateAdminReview(ctx context.Context, imdbID string, review string) (string, string, error) {
	sentiment, rankVal, err := s.getReviewRanking(ctx, review)
	if err != nil {
		return "", "", err
	}

	_, err = s.movieRepo.UpdateMovieReview(ctx, imdbID, review, sentiment, rankVal)
	if err != nil {
		return "", "", err
	}

	return sentiment, review, nil
}

func (s *movieService) GetRecommendedMovies(ctx context.Context, userId string) ([]models.Movie, error) {
	genres, err := s.userRepo.GetUserFavouriteGenres(ctx, userId)
	if err != nil {
		return nil, err
	}

	return s.movieRepo.GetRecommendedMovies(ctx, genres, s.config.RecommendedMovieLimit)
}

func (s *movieService) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	return s.movieRepo.GetAllGenres(ctx)
}

func (s *movieService) getReviewRanking(ctx context.Context, adminReview string) (string, int, error) {
	rankings, err := s.movieRepo.GetRankings(ctx)
	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""
	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited += ranking.RankingName + ","
		}
	}
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	if s.config.OpenAPIKey == "" {
		return "", 0, errors.New("could not read OPEN_API_KEY")
	}

	llm, err := openai.New(openai.WithToken(s.config.OpenAPIKey))
	if err != nil {
		return "", 0, err
	}

	base_prompt := strings.Replace(s.config.BasePromptTemplate, "{rankings}", sentimentDelimited, 1)
	response, err := llm.Call(ctx, base_prompt+"\n"+adminReview)
	if err != nil {
		return "", 0, err
	}

	// Clean response if needed (e.g. trim spaces)
	response = strings.TrimSpace(response)

	rankVal := 0
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil
}
