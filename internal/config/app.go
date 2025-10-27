package config

import (
	http "cutterproject/internal/delivery/http"
	"cutterproject/internal/delivery/http/middleware"
	"cutterproject/internal/delivery/http/route"
	"cutterproject/internal/repository"
	"cutterproject/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Router  *fiber.App
	DB      *pgxpool.Pool
	DBCache *redis.Client
	Log     *zap.Logger
	Config  *koanf.Koanf
}

func Server(config *ServerConfig) {
	userRepository := repository.NewUserRepository(config.Log, config.DB, config.DBCache)
	userUsecase := usecase.NewUserUsecase(userRepository, config.DB, config.Log, config.Config)
	userController := http.NewUserController(userUsecase, config.Log, config.Config)

	authMiddleware := middleware.NewAuthMiddleware(config.Router, config.Log, config.Config, userUsecase)

	routeConfig := route.RouteConfig{
		App:            config.Router,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}

	routeConfig.SetupRoute()
}
