package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/niko-2609/tracker-expense/constants"
	models "github.com/niko-2609/tracker-expense/models/auth"
	"golang.org/x/crypto/bcrypt"
)

func CreateJWTToken(userData models.UserCache) (string, error) {
	// Create a new token. Specify signing algorithm and an empty claims(payload)
	token := jwt.New(jwt.SigningMethodHS256)

	// Retrive the claims(payload) section of the JWT
	claims := token.Claims.(jwt.MapClaims)

	// Populate the claims
	claims["user_id"] = userData.ID
	claims["user_email"] = userData.Email

	// Sign the token with signing method defined above and our signing key
	jwtToken, err := token.SignedString([]byte(constants.SIGNING_KEY))
	if err != nil {
		return "", err
	}

	// Return new token
	return jwtToken, nil
}

// Create hash for he poasoed
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Compare password hashes
func CompareHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
