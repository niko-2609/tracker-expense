package logs

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var CustomLogger fiber.Handler

func init() {
	CustomLogger = logger.New(logger.Config{
		Format: "[${ip}]:${port} ${locals:requestid} ${status} - ${method} ${path}\n",
	})
}
