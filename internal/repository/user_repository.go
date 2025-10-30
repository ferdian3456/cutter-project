package repository

import (
	"context"
	"cutterproject/internal/constant"
	"cutterproject/internal/model"
	"errors"
	"fmt"
	"time"

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

// Postgresql - Nosql
func (repository *UserRepository) Register(ctx context.Context, tx pgx.Tx, user model.User) (int, error) {
	query := "INSERT INTO users (username,email,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id"

	var userId int
	err := tx.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&userId)
	if err != nil {
		return userId, err
	}

	return userId, nil
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

func (repository *UserRepository) GetUserAuth(ctx context.Context, email string) (int, string, error) {
	query := "SELECT id,password FROM users WHERE email=$1 LIMIT 1"

	var id int
	var passwordHash string

	err := repository.DB.QueryRow(ctx, query, email).Scan(&id, &passwordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return id, passwordHash, &model.ValidationError{
				Code:    constant.ERR_VALIDATION_CODE,
				Message: "Email is not found",
				Param:   "email",
			}
		}
		return id, passwordHash, err
	}

	return id, passwordHash, nil
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

// Redis - Cache
func (repository *UserRepository) SetAuthTokenInCache(ctx context.Context, accessToken string, refreeshToken string, userId int) error {
	accessTokenKey := fmt.Sprintf("auth:acccessToken:%d", userId)
	refreshTokenKey := fmt.Sprintf("auth:refreshToken:%d", userId)

	err := repository.DBCache.Set(ctx, accessTokenKey, accessToken, 15*time.Minute).Err()
	if err != nil {
		return err
	}

	err = repository.DBCache.Set(ctx, refreshTokenKey, refreeshToken, 15*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepository) GetAccessTokenInCache(ctx context.Context, userId int) (string, error) {
	accessTokenKey := fmt.Sprintf("auth:acccessToken:%d", userId)
	token, err := repository.DBCache.Get(ctx, accessTokenKey).Result()
	if err == redis.Nil {
		return token, &model.ValidationError{
			Code:    constant.ERR_NOT_FOUND_ERROR,
			Message: "Authorization token not found or expired",
			Param:   "accessToken",
		}
	} else if err != nil {
		return token, err
	}

	return token, nil
}
