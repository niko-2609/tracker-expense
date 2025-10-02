package validation

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CheckErrors(c *fiber.Ctx, errs []ValidationError, err error) string {
	var errorMsg string
	if errs != nil && len(errs) > 0 {
		messages := make([]string, 0, len(errs))
		for _, e := range errs {
			messages = append(messages, fmt.Sprintf("%s: %s", e.Field, e.Message))
		}
		errorMsg = "Invalid request - " + strings.Join(messages, "; ")
		return errorMsg
	}

	if err != nil {
		errorMsg = fmt.Sprintf("Invalid request: %s", err.Error())
	}
	return errorMsg
}
