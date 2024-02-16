package service

import (
	"context"
	"fmt"

	"github.com/ibookerke/capstone_backend/internal/config"
	"github.com/ibookerke/capstone_backend/internal/domain"
	"github.com/ibookerke/capstone_backend/internal/pkg/trm"
)

type UserService struct {
	trm                            trm.Manager
	userRepository                 domain.UserRepository
	notificationSettingsRepository domain.NotificationSettingsRepository
	cfg                            config.Auth
}

func NewUserService(
	trm trm.Manager,
	userRepository domain.UserRepository,
	notificationSettingsRepository domain.NotificationSettingsRepository,
	cfg config.Auth,
) *UserService {
	return &UserService{
		trm:                            trm,
		userRepository:                 userRepository,
		notificationSettingsRepository: notificationSettingsRepository,
		cfg:                            cfg,
	}
}

func (u *UserService) GetUserInfoById(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	user, err := u.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id error: %w", err)
	}

	nSettings, err := u.notificationSettingsRepository.GetNotificationSettingsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get notification settings by user id error: %w", err)
	}

	user.NotificationSettings = nSettings

	return user, nil
}
