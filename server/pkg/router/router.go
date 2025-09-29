package router

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/niko-2609/tracker-expense/pkg/handlers/auth"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
}
