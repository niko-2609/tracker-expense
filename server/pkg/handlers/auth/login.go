package handlers

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	authModel "github.com/niko-2609/tracker-expense/models/auth"
	models "github.com/niko-2609/tracker-expense/models/auth"
	apiModel "github.com/niko-2609/tracker-expense/models/common/api"
	"github.com/niko-2609/tracker-expense/pkg/validation"
	"github.com/niko-2609/tracker-expense/utils"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {
	input := new(authModel.Credentials)
	var usercache authModel.UserCache

	// Check raw incoming request
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		log.Error("Bad Request - Invalid Request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Request",
			"data":    nil,
		})
	}

	// Validate request using validator
	if err := validation.Validate.Struct(input); err != nil {
		log.Errorf("Bad Request - Request validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Failed request validation",
			Data:    nil,
		})
	}

	// Process request

	// Retrive email and password from request
	email := input.Email
	password := input.Password

	// Object to handle data from DB in DB's format
	userModel, err := new(authModel.User), *new(error)

	if utils.IsEmail(email) {
		userModel, err = utils.GetUserByEmail(email)
	}

	// If we have an error
	if err != nil {
		// Check if the record for the requested user exists. If not, return `unauthorized` and return 401
		if userModel == nil && errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("User not found in database")
			return c.Status(fiber.StatusUnauthorized).JSON(apiModel.Response{
				Status:  "error",
				Message: "Unauthorized user, please try again after signing up",
				Data:    nil,
			})
		}

		// For all other errors, return 500
		log.Error("An unexpected error occurred")
		return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
			Status:  "error",
			Message: "Internal server error",
			Data:    nil,
		})
	}

	// Checks passed, user can now be processed
	usercache = models.UserCache{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Username, // gives encoded string of password
		Password: userModel.Password,
	}

	// Check password validity
	if !utils.CompareHash(password, usercache.Password) {
		log.Error("Verification failed for password hashes")
		return c.Status(fiber.StatusUnauthorized).JSON(apiModel.Response{
			Status:  "error",
			Message: "User not authorized, please try again ",
			Data:    nil,
		})
	}

	token, err := utils.CreateJWTToken(usercache)
	if err != nil {
		log.Error("Error creating JWT token")
		return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
			Status:  "error",
			Message: "Login Failed, please try again",
			Data:    nil,
		})
	}

	// Construct new response payload
	resPayload := authModel.AccessPayload{
		Token: token,
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(apiModel.Response{
		Status:  "success",
		Message: "Login successfull",
		Data:    resPayload,
	})

}
