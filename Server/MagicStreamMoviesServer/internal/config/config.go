package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI              string
	DatabaseName          string
	SecretKey             string
	SecretRefreshKey      string
	OpenAPIKey            string
	BasePromptTemplate    string
	RecommendedMovieLimit int64
	AllowedOrigins        []string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	limitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")
	var limit int64 = 5
	if limitStr != "" {
		val, err := strconv.ParseInt(limitStr, 10, 64)
		if err == nil {
			limit = val
		}
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if allowedOrigins != "" {
		parts := strings.Split(allowedOrigins, ",")
		for _, p := range parts {
			origins = append(origins, strings.TrimSpace(p))
		}
	} else {
		origins = []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080"}
	}

	return &Config{
		MongoURI:              os.Getenv("MONGODB_URL"),
		DatabaseName:          os.Getenv("DATABASE_NAME"),
		SecretKey:             os.Getenv("SECRET_KEY"),
		SecretRefreshKey:      os.Getenv("SECRET_REFRESH_KEY"),
		OpenAPIKey:            os.Getenv("OPEN_API_KEY"),
		BasePromptTemplate:    os.Getenv("BASE_PROMPT_TEMPLATE"),
		RecommendedMovieLimit: limit,
		AllowedOrigins:        origins,
	}
}
