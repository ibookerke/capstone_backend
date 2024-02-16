package http

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/ibookerke/capstone_backend/internal/config"
	"github.com/ibookerke/capstone_backend/internal/service"
)

type NotificationSettingsHandler struct {
	logger                      *slog.Logger
	notificationSettingsService *service.NotificationSettingsService
	validator                   *validator.Validate
	cfg                         config.Auth
}

func NewNotificationSettingsHandler(
	logger *slog.Logger,
	notificationSettingsService *service.NotificationSettingsService,
	cfg config.Auth,
) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{
		logger:                      logger,
		notificationSettingsService: notificationSettingsService,
		cfg:                         cfg,
		validator:                   validator.New(),
	}
}

func (n *NotificationSettingsHandler) RegisterNotificationSettingsRoutes(r *gin.Engine) {
	authConfig := AuthConfig{JwtKey: n.cfg.JWTKey}
	r.Group("/notification-settings")
	{
		r.POST("/notification-settings/toggle", AuthMiddleware(authConfig), n.toggleNotificationSettingsOption)
	}
}

type ToggleNotificationSettingsRequest struct {
	Option  string `json:"option"`
	Enabled bool   `json:"enabled"`
}

func (n *NotificationSettingsHandler) toggleNotificationSettingsOption(c *gin.Context) {
	userID, err := GetUserId(c)
	if err != nil {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	var req ToggleNotificationSettingsRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": fmt.Errorf("Invalid input: %w", err).Error()})
		return
	}

	if err := n.validator.Struct(req); err != nil {
		c.JSON(400, gin.H{"message": fmt.Errorf("invalid input: %w", err).Error()})
		return
	}

	err = n.notificationSettingsService.ToggleNotificationSettingsOption(c, &service.ToggleNotificationSettingsInput{
		UserID:  userID,
		Option:  req.Option,
		Enabled: req.Enabled,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": fmt.Errorf("toggle notification settings error: %w", err).Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Success"})
}
