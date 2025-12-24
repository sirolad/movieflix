package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/config"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		SecretKey: "secret",
	}

	// Generate a valid token
	token, _, err := utils.GenerateAllTokens("test@example.com", "First", "Last", "USER", "user123", "secret", "refreshsecret")
	assert.NoError(t, err)

	tests := []struct {
		name           string
		token          string
		headerPrefix   string
		expectedStatus int
		expectedKeys   map[string]interface{}
	}{
		{
			name:           "Valid Token with Bearer",
			token:          token,
			headerPrefix:   "Bearer ",
			expectedStatus: http.StatusOK,
			expectedKeys: map[string]interface{}{
				"email":      "test@example.com",
				"first_name": "First",
				"last_name":  "Last",
				"role":       "USER",
				"user_id":    "user123",
			},
		},
		{
			name:           "Valid Token without Bearer",
			token:          token,
			headerPrefix:   "",
			expectedStatus: http.StatusOK,
			expectedKeys: map[string]interface{}{
				"email": "test@example.com",
			},
		},
		{
			name:           "Missing Token",
			token:          "",
			headerPrefix:   "",
			expectedStatus: http.StatusUnauthorized,
			expectedKeys:   nil,
		},
		{
			name:           "Invalid Token",
			token:          "invalid-token",
			headerPrefix:   "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedKeys:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest("GET", "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.headerPrefix+tt.token)
			}
			c.Request = req

			// Apply middleware
			handler := NewAuthMiddleware(cfg)
			
			// Dummy handler to check if next() was called and context set
			executed := false
			wrappedHandler := func(c *gin.Context) {
				handler(c)
				if !c.IsAborted() {
					executed = true
					// Check keys
					for k, v := range tt.expectedKeys {
						val, exists := c.Get(k)
						assert.True(t, exists, "Key %s should exist", k)
						assert.Equal(t, v, val)
					}
				}
			}

			wrappedHandler(c)

			if tt.expectedStatus == http.StatusOK {
				assert.True(t, executed, "Next handler should be executed")
			} else {
				assert.False(t, executed, "Next handler should NOT be executed")
				assert.Equal(t, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetRoleFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Test missing role
	_, err := GetRoleFromContext(c)
	assert.Error(t, err)
	assert.Equal(t, "role not found in context", err.Error())

	// Test invalid type
	c.Set(ContextKeyRole, 123)
	_, err = GetRoleFromContext(c)
	assert.Error(t, err)
	assert.Equal(t, "role is not a string", err.Error())

	// Test valid role
	c.Set(ContextKeyRole, "ADMIN")
	role, err := GetRoleFromContext(c)
	assert.NoError(t, err)
	assert.Equal(t, "ADMIN", role)
}
