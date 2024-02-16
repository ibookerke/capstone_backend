package service

import (
	"context"
	"fmt"

	"github.com/ibookerke/capstone_backend/internal/config"
	"github.com/ibookerke/capstone_backend/internal/domain"
	"github.com/ibookerke/capstone_backend/internal/pkg/trm"
)

type NotificationSettingsService struct {
	trm                            trm.Manager
	notificationSettingsRepository domain.NotificationSettingsRepository
	cfg                            config.Auth
}

func NewNotificationSettingsService(
	trm trm.Manager,
	notificationSettingsRepository domain.NotificationSettingsRepository,
	cfg config.Auth,
) *NotificationSettingsService {
	return &NotificationSettingsService{
		trm:                            trm,
		notificationSettingsRepository: notificationSettingsRepository,
		cfg:                            cfg,
	}
}

type ToggleNotificationSettingsInput struct {
	UserID  domain.UserID `json:"user_id" validate:"required"`
	Option  string        `json:"option" validate:"required"`
	Enabled bool          `json:"enabled" validate:"required"`
}

func (n *NotificationSettingsService) ToggleNotificationSettingsOption(
	ctx context.Context,
	input *ToggleNotificationSettingsInput,
) error {
	exists, err := n.notificationSettingsRepository.CheckIfUserNotificationOptionExists(
		ctx,
		input.UserID,
		input.Option,
	)
	if err != nil {
		return fmt.Errorf("error checking if user notification option exists: %w", err)
	}

	fmt.Println(input)

	if exists {
		err := n.notificationSettingsRepository.ToggleNotificationSettings(
			ctx,
			input.UserID,
			input.Option,
			input.Enabled,
		)
		if err != nil {
			return fmt.Errorf("error updating user notification option: %w", err)
		}
	} else {
		err := n.notificationSettingsRepository.CreateNotificationSettingsOption(
			ctx,
			input.UserID,
			input.Option,
			input.Enabled,
		)
		if err != nil {
			return fmt.Errorf("error creating user notification option: %w", err)
		}
	}
	return nil
}
