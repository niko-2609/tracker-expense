package utils

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix()

	// Sign the token with signing method defined above and our signing key
	jwtToken, err := token.SignedString([]byte(os.Getenv("KEY")))
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

// Extract user name from password
func ExtractUserName(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 1 {
		return beautifyName(parts[0])
	}
	return ""
}

func beautifyName(rawName string) string {
	words := strings.Split(rawName, ".")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}
