package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/niko-2609/tracker-expense/database"
	"github.com/niko-2609/tracker-expense/pkg/logs"
	"github.com/niko-2609/tracker-expense/pkg/router"
)

func main() {
	// New fiber app instance
	app := fiber.New()

	// Initialize requestid to track requests
	app.Use(requestid.New())

	// Log remote IP, port and request id.
	app.Use(logs.CustomLogger)

	// CORS settings
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Authorization, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Setup routing
	router.SetupRoutes(app)

	// Connect to database
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Exiting service, %s", err)
	}

	// Start server
	app.Listen(":3000")

}
