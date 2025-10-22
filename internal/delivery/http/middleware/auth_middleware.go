package middleware

import (
	"net/http"

	"cutterproject/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	Handler     *fiber.Handler
	Log         *zap.Logger
	Config      *koanf.Koanf
	UserUsecase *usecase.UserUsecase
}

func NewAuthMiddleware(handler *fiber.Handler, zap *zap.Logger, koanf *koanf.Koanf, userUsecase *usecase.UserUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		Handler:     handler,
		Log:         zap,
		Config:      koanf,
		UserUsecase: userUsecase,
	}
}

func (middleware *AuthMiddleware) ProtectedRoute(userUsecase *usecase.UserUsecase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Next()
	}
}
