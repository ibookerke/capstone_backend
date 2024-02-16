package domain

import "context"

type NotificationSettingsID int64

type NotificationSettings struct {
	ID                 NotificationSettingsID
	UserID             UserID
	NotificationOption string
	Enabled            bool
}

type NotificationSettingsRepository interface {
	GetNotificationSettingsByUserID(ctx context.Context, userID UserID) ([]NotificationSettings, error)
	CheckIfUserNotificationOptionExists(ctx context.Context, userID UserID, notificationOption string) (bool, error)
	CreateNotificationSettingsOption(ctx context.Context, userID UserID, notificationOption string, enabled bool) error
	ToggleNotificationSettings(ctx context.Context, userID UserID, notificationOption string, enabled bool) error
}
