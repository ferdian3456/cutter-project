package usecase

import (
	"cutterproject/internal/constant"
	"cutterproject/internal/model"
	"cutterproject/internal/repository"
	"cutterproject/internal/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepository *repository.UserRepository
	DB             *pgxpool.Pool
	Log            *zap.Logger
	Config         *koanf.Koanf
}

func NewUserUsecase(userRepository *repository.UserRepository, db *pgxpool.Pool, zap *zap.Logger, koanf *koanf.Koanf) *UserUsecase {
	return &UserUsecase{
		UserRepository: userRepository,
		DB:             db,
		Log:            zap,
		Config:         koanf,
	}
}

func (usecase *UserUsecase) Register(ctx *fiber.Ctx, payload model.UserCreateRequest) (model.TokenResponse, error) {
	ctxContext := ctx.Context()
	token := model.TokenResponse{}

	if payload.Username == "" {
		return token, &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "Username is required to not be empty",
			Param:   "username",
		}
	} else if len(payload.Username) < 4 {
		return token, &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "Username must be at least 4 characters",
			Param:   "username",
		}
	} else if len(payload.Username) > 22 {
		return token, &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "username must be at most 22 characters",
			Param:   "username",
		}
	}

	if payload.Email == "" {
		return token, &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "Email is required to not be empty",
			Param:   "email",
		}
	} else if len(payload.Email) < 16 {
		return token, &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "email must be at least 16 characters",
			Param:   "email",
		}
	} else if len(payload.Email) > 80 {
		return token, &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "Email must be at most 80 characters",
			Param:   "email",
		}
	}

	if payload.Password == "" {
		return token, &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "Password is required to not be empty",
			Param:   "email",
		}
	} else if len(payload.Password) < 5 {
		return token, &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "Password must be at least 5 characters",
			Param:   "email",
		}
	} else if len(payload.Password) > 20 {
		return token, &model.ValidationError{
			Code:    constant.ERR_VALIDATION_CODE,
			Message: "Password must be at most 20 characters",
			Param:   "email",
		}
	}

	err := usecase.UserRepository.CheckUsernameOrEmailUnique(ctxContext, payload.Username, payload.Email)
	if err != nil {
		return token, err
	}

	now := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return token, err
	}

	user := model.User{
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	// start transaction
	tx, err := usecase.DB.Begin(ctx.Context())
	if err != nil {
		return token, err
	}

	defer tx.Rollback(ctxContext)

	userId, err := usecase.UserRepository.Register(ctxContext, tx, user)
	if err != nil {
		return token, err
	}

	err = tx.Commit(ctxContext)
	if err != nil {
		return token, err
	}

	token, err = util.GenerateTokenPair(userId, usecase.Config.String("JWT_SECRET_KEY"))
	if err != nil {
		return token, err
	}

	return token, nil
}

func (usecase *UserUsecase) GetUserInfo(ctx *fiber.Ctx, id int) (model.UserResponse, error) {
	user, err := usecase.UserRepository.GetUserInfo(ctx.Context(), id)
	if err != nil {
		return user, err
	}

	return user, nil
}
