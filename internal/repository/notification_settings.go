package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibookerke/capstone_backend/internal/domain"
	trmpgx "github.com/ibookerke/capstone_backend/internal/pkg/pgx"
	"github.com/ibookerke/capstone_backend/internal/pkg/trm"
)

type NotificationSettingsRepository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	trm    trm.Manager
}

func NewNotificationSettingsRepository(
	pool *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	trm trm.Manager,
) *NotificationSettingsRepository {
	return &NotificationSettingsRepository{
		pool:   pool,
		getter: getter,
		trm:    trm,
	}
}

const CheckIfUserNotificationOptionExistsSQL = `SELECT EXISTS(
	SELECT 1
	FROM har.user_notification_settings
	WHERE user_id = $1
	AND notification_option = $2
)`

const CreateUserNotificationOptionSQL = `INSERT INTO har.user_notification_settings
	(user_id, notification_option, enabled)
		VALUES
	($1, $2, $3)`

const GetUserNotificationSettingsByUserIDSQL = `SELECT
	id, user_id, notification_option, enabled
	FROM har.user_notification_settings
	WHERE user_id = $1`

const ToggleNotificationSettingsSQL = `UPDATE har.user_notification_settings
	SET enabled = $1
	WHERE user_id = $2
	AND notification_option = $3`

func (n NotificationSettingsRepository) CheckIfUserNotificationOptionExists(ctx context.Context, userID domain.UserID, notificationOption string) (bool, error) {
	var exists bool
	err := n.pool.QueryRow(ctx, CheckIfUserNotificationOptionExistsSQL, userID, notificationOption).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (n NotificationSettingsRepository) CreateNotificationSettingsOption(ctx context.Context, userID domain.UserID, notificationOption string, enabled bool) error {
	fmt.Println("repa:", userID, notificationOption, enabled)
	_, err := n.pool.Exec(
		ctx,
		CreateUserNotificationOptionSQL,
		userID,
		notificationOption,
		enabled,
	)
	return err
}

func (n NotificationSettingsRepository) GetNotificationSettingsByUserID(ctx context.Context, userID domain.UserID) ([]domain.NotificationSettings, error) {
	rows, err := n.pool.Query(
		ctx,
		GetUserNotificationSettingsByUserIDSQL,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []domain.NotificationSettings
	for rows.Next() {
		var s domain.NotificationSettings
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.NotificationOption,
			&s.Enabled)
		if err != nil {
			return nil, err
		}
		settings = append(settings, s)
	}
	return settings, nil
}

func (n NotificationSettingsRepository) ToggleNotificationSettings(ctx context.Context, userID domain.UserID, notificationOption string, enabled bool) error {
	_, err := n.pool.Exec(ctx, ToggleNotificationSettingsSQL, enabled, userID, notificationOption)
	return err
}
