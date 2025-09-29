package validation

import "github.com/go-playground/validator/v10"

// Global Validate object
var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}
