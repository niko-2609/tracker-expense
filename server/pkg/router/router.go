package router

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/niko-2609/tracker-expense/pkg/handlers/auth"
	middleware "github.com/niko-2609/tracker-expense/pkg/middleware/auth"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/register", handlers.SignUp)

	//test
	test := api.Group("/test")
	test.Get("", middleware.Protected(), func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "middleware authentication is working",
		})
	})
}
