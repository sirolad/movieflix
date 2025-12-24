package models

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestUserValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "Valid User",
			user: User{
				UserID:    "user123",
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Role:      "USER",
			},
			wantErr: false,
		},
		{
			name: "Invalid Email",
			user: User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "invalid-email",
				Password:  "password123",
				Role:      "USER",
			},
			wantErr: true,
		},
		{
			name: "Short Password",
			user: User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "123",
				Role:      "USER",
			},
			wantErr: true,
		},
		{
			name: "Invalid Role",
			user: User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Role:      "INVALID",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMovieValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		movie   Movie
		wantErr bool
	}{
		{
			name: "Valid Movie",
			movie: Movie{
				ImdbID:     "tt1234567",
				Title:      "Test Movie",
				PosterPath: "http://example.com/poster.jpg",
				YouTubeID:  "dQw4w9WgXcQ",
				Genre: []Genre{
					{GenreID: 1, GenreName: "Action"},
				},
				Ranking: Ranking{
					RankingValue: 8,
					RankingName:  "Good",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing Title",
			movie: Movie{
				ImdbID:     "tt1234567",
				PosterPath: "http://example.com/poster.jpg",
				YouTubeID:  "dQw4w9WgXcQ",
			},
			wantErr: true,
		},
		{
			name: "Invalid Poster URL",
			movie: Movie{
				ImdbID:     "tt1234567",
				Title:      "Test Movie",
				PosterPath: "not-a-url",
				YouTubeID:  "dQw4w9WgXcQ",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.movie)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
