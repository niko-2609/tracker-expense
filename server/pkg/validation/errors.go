package validation

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CheckErrors(c *fiber.Ctx, errs []ValidationError, err error) string {
	var errorMsg string
	if errs != nil {
		if len(errs) == 1 {
			errorMsg = fmt.Sprintf("Invalid %s: %s", errs[0].Field, errs[0].Message)
		} else {
			errorMsg = "Invalid request: email or password is not valid."
		}
	} else {
		errorMsg = fmt.Sprintf("Invalid request: %s", err.Error())
	}
	return errorMsg
}
