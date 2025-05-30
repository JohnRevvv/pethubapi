package middleware

import (
	"fmt"
	"log"
	"os"
	"pethub_api/models/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load() // Load from .env file
	if err != nil {
		log.Println("Warning: No .env file found")
	}
}

// Secret key for signing tokens (should be stored in env variables)
var SecretKey = os.Getenv("SECRET_KEY")

// GenerateJWT generates a new JWT token
func GenerateJWT(ID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Expires in 72 hours

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: "Unauthorized: No token provided",
				Data:    nil,
			})
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: "Unauthorized: Invalid token",
				Data:    nil,
			})
		}

		fmt.Println("Token received:", tokenString)

		// Get "id" from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: "Unauthorized: Invalid token claims",
				Data:    nil,
			})
		}

		adopterID, ok := claims["id"].(float64)
		if !ok {
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: "Unauthorized: Missing adopter ID in token",
				Data:    nil,
			})
		}

		// Store adopter ID in context
		c.Locals("adopter_id", uint(adopterID))
		return c.Next()

	}
}

// GetAdopterIDFromJWT retrieves the adopter ID from the JWT token stored in the Fiber context
func GetAdopterIDFromJWT(c *fiber.Ctx) (uint, error) {
	adopterID, ok := c.Locals("adopter_id").(uint) // JWT claims are stored as float64 in Go
	if !ok {
		return 0, fmt.Errorf("adopter ID not found in context")
	}
	return uint(adopterID), nil
}
