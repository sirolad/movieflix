package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/models"
	"github.com/sirolad/MagicStreamMovies/Server/MagicStreamMovies/Server/internal/service"
)

type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{
		service:  s,
		validate: validator.New(),
	}
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Register a new user with the provided details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User Registration Data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if err := h.validate.Struct(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	createdUser, err := h.service.RegisterUser(ctx, user)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully", "user": createdUser})
}

// LoginUser godoc
// @Summary      Login user
// @Description  Login user and return tokens
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.UserLogin  true  "User Login Credentials"
// @Success      200          {object}  models.UserResponse
// @Failure      400          {object}  map[string]interface{}
// @Failure      401          {object}  map[string]interface{}
// @Failure      500          {object}  map[string]interface{}
// @Router       /login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	var userLogin models.UserLogin
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if err := h.validate.Struct(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	resp, err := h.service.LoginUser(ctx, userLogin.Email, userLogin.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid email or password" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    resp.Token,
		Path:     "/",
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(http.StatusOK, resp)
}

// LogoutHandler godoc
// @Summary      Logout user
// @Description  Logout user and clear tokens
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      map[string]string  true  "User ID JSON {\"user_id\": \"id\"}"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /user/logout [post]
func (h *UserHandler) LogoutHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	var req struct {
		UserId string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.service.LogoutUser(ctx, req.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging out user"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

// RefreshTokenHandler godoc
// @Summary      Refresh user tokens
// @Description  Refresh access and refresh tokens
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /user/refresh-token [post]
func (h *UserHandler) RefreshTokenHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
	defer cancel()

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to retrieve refresh token from cookie"})
		return
	}

	token, newRefreshToken, err := h.service.RefreshToken(ctx, refreshToken)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid or expired refresh token" || err.Error() == "user not found" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		MaxAge:   604800, // 7 days
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
}
