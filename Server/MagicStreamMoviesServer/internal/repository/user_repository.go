package repository

import (
	"context"
	"errors"
	"time"

	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (*mongo.InsertOneResult, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CountUsersByEmail(ctx context.Context, email string) (int64, error)
	UpdateTokens(ctx context.Context, userId string, token string, refreshToken string, updatedAt time.Time) error
	GetUserFavouriteGenres(ctx context.Context, userId string) ([]string, error)
}

type mongoUserRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepository{
		userCollection: db.Collection("users"),
	}
}

func (r *mongoUserRepository) CreateUser(ctx context.Context, user models.User) (*mongo.InsertOneResult, error) {
	return r.userCollection.InsertOne(ctx, user)
}

func (r *mongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) CountUsersByEmail(ctx context.Context, email string) (int64, error) {
	return r.userCollection.CountDocuments(ctx, bson.M{"email": email})
}

func (r *mongoUserRepository) UpdateTokens(ctx context.Context, userId string, token string, refreshToken string, updatedAt time.Time) error {
	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    updatedAt,
		},
	}
	_, err := r.userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	return err
}

func (r *mongoUserRepository) GetUserFavouriteGenres(ctx context.Context, userId string) ([]string, error) {
	filter := bson.D{{Key: "user_id", Value: userId}}
	projection := bson.M{"favourite_genres": 1, "_id": 0}
	var result bson.M

	opts := options.FindOne().SetProjection(projection)
	err := r.userCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return nil, err
	}

	favGenresArray, ok := result["favourite_genres"].(bson.A)
	if !ok {
		// If it's nil or not an array, maybe return empty? 
		// Or check if user exists first? 
		// Keeping logic similar to original.
		return []string{}, errors.New("unable to retrieve favourite genres for user")
	}

	var genreNames []string
	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if genreName, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, genreName)
					}
				}
			}
		}
	}
	return genreNames, nil
}
