package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/ibookerke/capstone_backend/internal/config"
	"github.com/ibookerke/capstone_backend/internal/domain"
	"github.com/ibookerke/capstone_backend/internal/service"
)

type AuthHandler struct {
	logger      *slog.Logger
	authService *service.AuthService
	validator   *validator.Validate
	cfg         config.Auth
}

func NewAuthHandler(
	logger *slog.Logger,
	authService *service.AuthService,
	cfg config.Auth,
) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
		cfg:         cfg,
		validator:   validator.New(),
	}
}

func (a *AuthHandler) RegisterAuthRoutes(r *gin.Engine) {
	r.POST("/auth/login", a.login)
	r.POST("/auth/register", a.register)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func (a *AuthHandler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request: %w", err).Error()})
		return
	}

	if err := a.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request: %w", err).Error()})
		return
	}

	user, err := a.authService.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("login failed: %w", err).Error()})
		return
	}

	claims := &Claims{
		User: *user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(a.cfg.JWTKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("signing token error: %w", err).Error()})
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: signedToken,
		User:  *user,
	})
	return
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
}

func (a *AuthHandler) register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request: %w", err).Error()})
		return
	}

	if err := a.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request: %w", err).Error()})
		return
	}

	user, err := a.authService.Register(c, req.Email, req.Password, req.Name, req.Surname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("register failed: %w", err).Error()})
		return
	}

	c.JSON(http.StatusOK, user)
	return
}
