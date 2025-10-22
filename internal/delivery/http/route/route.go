package route

import (
	"cutterproject/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
	//UserController *http.UserController
}

func (c *RouteConfig) SetupRoute() {
	api := c.App.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	userGroup := api.Group("/user", c.App.Use(c.AuthMiddleware))
}
