package route

import (
	"cutterproject/internal/delivery/http"
	"cutterproject/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware *middleware.AuthMiddleware
	UserController *http.UserController
}

func (c *RouteConfig) SetupRoute() {
	api := c.App.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	authGroup := api.Group("/auth")
	authGroup.Post("/register", c.UserController.Register)
	authGroup.Post("/login", c.UserController.Login)

	userGroup := api.Group("/users", c.AuthMiddleware.ProtectedRoute())
	userGroup.Get("/me", c.UserController.GetUserInfo)
	//userGroup.Get("/:userId", c.UserController.GetUserInfo)
	//userGroup.Delete("/:userId")
}
