package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims defines the structure of the JWT payload
type JWTClaims struct {
	AdopterID uint `json:"adopter_id"`
	jwt.RegisteredClaims
}

// Temporary ValidateToken function for testing
// It skips the actual validation and returns a dummy valid claim
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// Temporary bypass for testing (skip real token validation)
	// For now, return a dummy valid claim with AdopterID 1
	return &JWTClaims{
		AdopterID: 1, // Or any valid AdopterID for testing purposes
	}, nil
}

// Original GenerateToken function remains the same
func GenerateToken(adopterID uint) (string, error) {
	claims := JWTClaims{
		AdopterID: adopterID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token valid for 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secretKey))
}
