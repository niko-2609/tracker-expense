package utils

import (
	"net/mail"

	"github.com/niko-2609/tracker-expense/database"
	models "github.com/niko-2609/tracker-expense/models/auth"
	"gorm.io/gorm"
)

// Check if string is a valid email
func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Get user from DB by email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where(&models.User{Email: email}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Get user by ID
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Model(&models.User{}).Where(&models.User{Model: gorm.Model{
		ID: id,
	}}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
