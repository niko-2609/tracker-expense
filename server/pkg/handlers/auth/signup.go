package handlers

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/niko-2609/tracker-expense/database"
	authmodel "github.com/niko-2609/tracker-expense/models/auth"
	apiModel "github.com/niko-2609/tracker-expense/models/common/api"
	"github.com/niko-2609/tracker-expense/pkg/validation"
	"github.com/niko-2609/tracker-expense/utils"
	"gorm.io/gorm"
)

func SignUp(c *fiber.Ctx) error {
	input := new(authmodel.Credentials)

	// Verify raw request
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		log.Error("Bad Request - Invalid Credentials Object")
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Invalid Request",
			Data:    nil,
		})
	}

	// Verify request against the defined model
	if err := validation.Validate.Struct(input); err != nil {
		log.Error("Bad Request - Request validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Failed request validation",
			Data:    nil,
		})
	}

	email := input.Email
	password := input.Password

	var userModel *authmodel.User

	// Check if user exists in DB
	userModel, err := utils.GetUserByEmail(email)

	if err != nil {
		// For all other errors that `RecordNotFound`, we return an error.
		// This is done because `RecordNotFound` is not an error in this case and is valid
		// Since we will proceed if a record for the user does not exist in DB.
		if userModel == nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// For all other errors, return 500
			log.Error("An unexpected error occurred")
			return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
				Status:  "error",
				Message: "Internal server error",
				Data:    nil,
			})
		}
	}

	// If no errors and a valid user is present, return 409
	if userModel != nil {
		log.Warn("User already exists, should be redirected to /login")
		return c.Status(fiber.StatusConflict).JSON(apiModel.Response{
			Status:  "error",
			Message: "User exists, try signing in",
			Data:    nil,
		})
	}

	if !utils.IsEmail(email) {
		log.Error("Email provided is not valid")
		return c.Status(fiber.StatusBadRequest).JSON(apiModel.Response{
			Status:  "error",
			Message: "Invalid email or password, please try again",
			Data:    nil,
		})
	}

	//  Encrypt password
	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		log.Error("Error encrypting password")
		return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
			Status:  "error",
			Message: "Unable to sign up user, please try again",
			Data:    nil,
		})
	}

	// // Get Username from Email
	// username, err := utils.ExtractUserName(email)
	// if err != nil {
	// 	log.Error(err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
	// 		Status:  "error",
	// 		Message: "Unable to sign up user, provide a valid email",
	// 		Data:    nil,
	// 	})
	// }

	userModel = &authmodel.User{
		Email:    email,
		Password: hashedPass,
		Username: utils.ExtractUserName(email),
	}

	// Save user to DB
	if err := database.DB.Create(userModel).Error; err != nil {
		log.Error("Unable to create user in DB")
		return c.Status(fiber.StatusInternalServerError).JSON(apiModel.Response{
			Status:  "error",
			Message: "Internal server error, please try again",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(apiModel.Response{
		Status:  "success",
		Message: "Sign up successfull",
		Data:    nil,
	})

}
