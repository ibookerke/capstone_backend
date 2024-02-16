package http

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/ibookerke/capstone_backend/internal/config"
	"github.com/ibookerke/capstone_backend/internal/service"
)

type UserHandler struct {
	logger      *slog.Logger
	userHandler *service.UserService
	cfg         config.Auth
}

func NewUserHandler(
	logger *slog.Logger,
	UserService *service.UserService,
	cfg config.Auth,
) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userHandler: UserService,
		cfg:         cfg,
	}
}

func (u *UserHandler) RegisterUserRoutes(r *gin.Engine) {
	authConfig := AuthConfig{JwtKey: u.cfg.JWTKey}

	r.GET("/user/current", AuthMiddleware(authConfig), u.getUser)

}

func (u *UserHandler) getUser(c *gin.Context) {
	userId, err := GetUserId(c)
	if err != nil {
		c.JSON(401, gin.H{"error": fmt.Errorf("get user id error: %w", err).Error()})
		return
	}

	user, err := u.userHandler.GetUserInfoById(c, userId)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Errorf("get user info error: %w", err).Error()})
	}

	c.JSON(200, gin.H{"user": user})
}
