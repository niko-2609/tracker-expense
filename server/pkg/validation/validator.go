package validation

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationError struct {
	Field   string
	Message string
}

func ValidateRequest(c *fiber.Ctx, input any) ([]ValidationError, error) {
	// Check raw incoming request
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return nil, err
	}

	// Validate request using validator
	if err := Validate.Struct(input); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			var ve []ValidationError
			for _, e := range errs {
				ve = append(ve, ValidationError{
					Field:   e.Field(),
					Message: messageForTag(e),
				})
			}
			return ve, fmt.Errorf("Validation failed")
		}
		return nil, fmt.Errorf("Unable to validate request body")
	}
	return nil, nil
}

func messageForTag(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "format is not valid"
	case "min":
		return fmt.Sprintf("minimum length must be %s", fieldErr.Param())
	case "max":
		return fmt.Sprintf("maximum length must be  %s", fieldErr.Param())
	default:
		return "Invalid value"
	}

}
