package service

import (
	"context"
	"errors"
	"time"

	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/config"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/repository"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService interface {
	RegisterUser(ctx context.Context, user models.User) (*models.User, error)
	LoginUser(ctx context.Context, email, password string) (*models.UserResponse, error)
	LogoutUser(ctx context.Context, userId string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

type userService struct {
	repo   repository.UserRepository
	config *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		repo:   repo,
		config: cfg,
	}
}

func (s *userService) RegisterUser(ctx context.Context, user models.User) (*models.User, error) {
	count, err := s.repo.CountUsersByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.ID = bson.NewObjectID()
	user.UserID = user.ID.Hex()

	_, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) LoginUser(ctx context.Context, email, password string) (*models.UserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	signedToken, signedRefreshToken, err := utils.GenerateAllTokens(
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
		user.UserID,
		s.config.SecretKey,
		s.config.SecretRefreshKey,
	)
	if err != nil {
		return nil, err
	}

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	err = s.repo.UpdateTokens(ctx, user.UserID, signedToken, signedRefreshToken, updateAt)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		UserID:         user.UserID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Role:           user.Role,
		Token:          signedToken,
		RefreshToken:   signedRefreshToken,
		FavoriteGenres: user.FavoriteGenres,
	}, nil
}

func (s *userService) LogoutUser(ctx context.Context, userId string) error {
	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	return s.repo.UpdateTokens(ctx, userId, "", "", updateAt)
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := utils.ValidateToken(refreshToken, s.config.SecretRefreshKey)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	user, err := s.repo.GetUserByUserID(ctx, claims.UserId)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	signedToken, signedRefreshToken, err := utils.GenerateAllTokens(
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
		user.UserID,
		s.config.SecretKey,
		s.config.SecretRefreshKey,
	)
	if err != nil {
		return "", "", err
	}

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	err = s.repo.UpdateTokens(ctx, user.UserID, signedToken, signedRefreshToken, updateAt)
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}
