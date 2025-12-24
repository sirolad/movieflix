package repository

import (
	"context"

	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MovieRepository interface {
	GetMovies(ctx context.Context) ([]models.Movie, error)
	GetMovie(ctx context.Context, imdbID string) (*models.Movie, error)
	CreateMovie(ctx context.Context, movie models.Movie) (*mongo.InsertOneResult, error)
	UpdateMovieReview(ctx context.Context, imdbID string, review string, sentiment string, rankVal int) (*mongo.UpdateResult, error)
	GetRankings(ctx context.Context) ([]models.Ranking, error)
	GetRecommendedMovies(ctx context.Context, genres []string, limit int64) ([]models.Movie, error)
	GetAllGenres(ctx context.Context) ([]models.Genre, error)
}

type mongoMovieRepository struct {
	movieCollection   *mongo.Collection
	rankingCollection *mongo.Collection
}

func NewMovieRepository(db *mongo.Database) MovieRepository {
	return &mongoMovieRepository{
		movieCollection:   db.Collection("movies"),
		rankingCollection: db.Collection("rankings"),
	}
}

func (r *mongoMovieRepository) GetMovies(ctx context.Context) ([]models.Movie, error) {
	var movies []models.Movie
	cursor, err := r.movieCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *mongoMovieRepository) GetMovie(ctx context.Context, imdbID string) (*models.Movie, error) {
	var movie models.Movie
	err := r.movieCollection.FindOne(ctx, bson.M{"imdb_id": imdbID}).Decode(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *mongoMovieRepository) CreateMovie(ctx context.Context, movie models.Movie) (*mongo.InsertOneResult, error) {
	return r.movieCollection.InsertOne(ctx, movie)
}

func (r *mongoMovieRepository) UpdateMovieReview(ctx context.Context, imdbID string, review string, sentiment string, rankVal int) (*mongo.UpdateResult, error) {
	filter := bson.M{"imdb_id": imdbID}
	update := bson.M{
		"$set": bson.M{
			"admin_review": review,
			"ranking": bson.M{
				"ranking_name":  sentiment,
				"ranking_value": rankVal,
			},
		},
	}
	return r.movieCollection.UpdateOne(ctx, filter, update)
}

func (r *mongoMovieRepository) GetRankings(ctx context.Context) ([]models.Ranking, error) {
	var rankings []models.Ranking
	cursor, err := r.rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}

func (r *mongoMovieRepository) GetRecommendedMovies(ctx context.Context, genres []string, limit int64) ([]models.Movie, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
	findOptions.SetLimit(limit)

	filter := bson.M{"genre.genre_name": bson.M{"$in": genres}}

	cursor, err := r.movieCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var recommendedMovies []models.Movie
	if err = cursor.All(ctx, &recommendedMovies); err != nil {
		return nil, err
	}
	return recommendedMovies, nil
}

func (r *mongoMovieRepository) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	// We use aggregation to unwind the genres array and then group by genre_id to get unique genres
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$genre"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$genre.genre_id"},
			{Key: "genre_name", Value: bson.D{{Key: "$first", Value: "$genre.genre_name"}}},
			{Key: "genre_id", Value: bson.D{{Key: "$first", Value: "$genre.genre_id"}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "genre_id", Value: 1},
			{Key: "genre_name", Value: 1},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "genre_name", Value: 1}}}},
	}

	cursor, err := r.movieCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var genres []models.Genre
	if err = cursor.All(ctx, &genres); err != nil {
		return nil, err
	}
	return genres, nil
}
