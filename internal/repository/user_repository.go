package repository

import (
	"context"
	"cutterproject/internal/constant"
	"cutterproject/internal/model"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserRepository struct {
	Log     *zap.Logger
	DB      *pgxpool.Pool
	DBCache *redis.Client
}

func NewUserRepository(zap *zap.Logger, db *pgxpool.Pool, dbCache *redis.Client) *UserRepository {
	return &UserRepository{
		Log:     zap,
		DB:      db,
		DBCache: dbCache,
	}
}

func (repository *UserRepository) Register(ctx context.Context, tx pgx.Tx, user model.User) error {
	query := "INSERT INTO users (username,email,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5)"
	_, err := tx.Exec(ctx, query, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepository) CheckUsernameOrEmailUnique(ctx context.Context, username string, email string) error {
	query := "SELECT username,email FROM users WHERE username=$1 OR email=$2 LIMIT 1"

	var existUsername string
	var existEmail string
	err := repository.DB.QueryRow(ctx, query, username, email).Scan(&existUsername, &existEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	if existUsername == username {
		return &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "Username is already exist",
			Param:   "username",
		}
	}

	if existEmail == email {
		return &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "Email is already exist",
			Param:   "email",
		}
	}

	return nil
}

func (repository *UserRepository) GetUserInfo(ctx context.Context, id int) (model.UserResponse, error) {
	query := "SELECT id,username,email, created_at,updated_at FROM users WHERE id=$1 LIMIT 1"

	user := model.UserResponse{}
	err := repository.DB.QueryRow(ctx, query, id).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, &model.ValidationError{
				Code:    constant.ERR_NOT_FOUND_ERROR,
				Message: "User not found",
				Param:   "userId",
			}
		}
		return user, err
	}

	return user, nil
}
