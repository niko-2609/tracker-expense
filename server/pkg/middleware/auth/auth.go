package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	apimodel "github.com/niko-2609/tracker-expense/models/common/api"
)

func Protected() fiber.Handler {
	return jwtware.New(
		jwtware.Config{
			SigningKey: jwtware.SigningKey{
				Key: []byte(os.Getenv("KEY")),
			},
			ErrorHandler: jwtError,
		},
	)
}

// This not just any handler, it takes an `err`
// along with `fiber.Ctx`. Its an `ErrorHandler`.
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(apimodel.Response{
			Status:  "error",
			Message: "Unable to verify user",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(apimodel.Response{
		Status:  "error",
		Message: err.Error(),
		Data:    nil,
	})
}
