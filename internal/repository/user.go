package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibookerke/capstone_backend/internal/domain"
	trmpgx "github.com/ibookerke/capstone_backend/internal/pkg/pgx"
	"github.com/ibookerke/capstone_backend/internal/pkg/trm"
)

type UserRepository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	trm    trm.Manager
}

func NewUserRepository(
	pool *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	trm trm.Manager,
) *UserRepository {
	return &UserRepository{
		pool:   pool,
		getter: getter,
		trm:    trm,
	}
}

const createUserSQL = `INSERT INTO har.users 
		(email, password, name, surname, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6)
	RETURNING id`

const updateUserSQL = `UPDATE har.users
	SET email = $1, name = $2, surname = $3, updated_at = $4
	WHERE id = $5`

const getUserByEmailSQL = `SELECT 
    	users.id, users.email, users.password, users.name, users.surname, users.created_at, users.updated_at
	FROM har.users
	where users.email = $1`

const getUserByIDSQL = `SELECT
		users.id, users.email, users.password, users.name, users.surname, users.created_at, users.updated_at
	FROM har.users
	where users.id = $1`

func (a *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := a.trm.Do(ctx, func(ctx context.Context) error {
		exec := a.getter.DefaultTrOrDB(ctx, a.pool)

		pwd, err := HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		err = exec.QueryRow(
			ctx,
			createUserSQL,
			user.Email,
			pwd,
			user.Name,
			user.Surname,
			time.Now(),
			time.Now(),
		).Scan(
			&user.ID,
		)
		if err != nil {
			return fmt.Errorf("save user: %w", wrapScanError(err))
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (a *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := a.trm.Do(ctx, func(ctx context.Context) error {
		exec := a.getter.DefaultTrOrDB(ctx, a.pool)

		_, err := exec.Exec(
			ctx,
			updateUserSQL,
			user.Email,
			user.Name,
			user.Surname,
			time.Now(),
			user.ID,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("update user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (a *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	err := a.trm.Do(ctx, func(ctx context.Context) error {
		exec := a.getter.DefaultTrOrDB(ctx, a.pool)

		err := exec.QueryRow(
			ctx,
			getUserByEmailSQL,
			email,
		).Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Surname,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return fmt.Errorf("get user by email: %w", wrapScanError(err))
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (a *UserRepository) GetUserByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var user domain.User

	err := a.trm.Do(ctx, func(ctx context.Context) error {
		exec := a.getter.DefaultTrOrDB(ctx, a.pool)

		err := exec.QueryRow(
			ctx,
			getUserByIDSQL,
			id,
		).Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Surname,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("get user by id: %w", wrapScanError(err))
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}
