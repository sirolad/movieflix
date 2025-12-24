package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Save original env vars
	originalMongoURI := os.Getenv("MONGODB_URL")
	originalLimit := os.Getenv("RECOMMENDED_MOVIE_LIMIT")
	defer func() {
		os.Setenv("MONGODB_URL", originalMongoURI)
		os.Setenv("RECOMMENDED_MOVIE_LIMIT", originalLimit)
	}()

	// Set test env vars
	os.Setenv("MONGODB_URL", "mongodb://localhost:27017/test")
	os.Setenv("RECOMMENDED_MOVIE_LIMIT", "10")
	os.Setenv("SECRET_KEY", "testsecret")

	cfg := LoadConfig()

	assert.Equal(t, "mongodb://localhost:27017/test", cfg.MongoURI)
	assert.Equal(t, int64(10), cfg.RecommendedMovieLimit)
	assert.Equal(t, "testsecret", cfg.SecretKey)
	assert.Contains(t, cfg.AllowedOrigins, "http://localhost:3000")
}

func TestLoadConfigDefaults(t *testing.T) {
	// Clear env vars to test defaults
	os.Unsetenv("RECOMMENDED_MOVIE_LIMIT")
	os.Unsetenv("ALLOWED_ORIGINS")

	cfg := LoadConfig()

	assert.Equal(t, int64(5), cfg.RecommendedMovieLimit)
	assert.Contains(t, cfg.AllowedOrigins, "http://localhost:3000")
}
