package controllers

// import (
// 	"errors"
// 	"pethub_api/middleware"
// 	"pethub_api/models"

// 	"github.com/gofiber/fiber/v2"
// 	"golang.org/x/crypto/bcrypt"
// 	"gorm.io/gorm"
// )


// // Login
// func LoginShelter(c *fiber.Ctx) error {
// 	// Parse request body
// 	requestBody := struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}{}

// 	if err := c.BodyParser(&requestBody); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Invalid request body",
// 		})
// 	}

// 	// Check if the adopter exists
// 	var ShelterAccount models.ShelterAccount
// 	result := middleware.DBConn.Where("username = ?", requestBody.Username).First(&ShelterAccount)

// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"message": "Invalid username or password",
// 		})
// 	} else if result.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Database error",
// 		})
// 	}

// 	// Check password using bcrypt
// 	err := bcrypt.CompareHashAndPassword([]byte(ShelterAccount.Password), []byte(requestBody.Password))
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"message": "Invalid username or password",
// 		})
// 	}

// 	// Fetch shelter info
// 	var ShelterInfo models.ShelterInfo
// 	infoResult := middleware.DBConn.Where("shelter_id = ?", ShelterAccount.ShelterID).First(&ShelterInfo)

// 	if infoResult.Error != nil && !errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to fetch adopter info",
// 		})
// 	}

// 	// Login successful, return shelter account, info, and ID
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Login successful",
// 		"data": fiber.Map{
// 			"shelter_id": ShelterAccount.ShelterID, // Include shelter ID in the response
// 			"Shelter":    ShelterAccount,
// 			"Info":       ShelterInfo,
// 		},
// 	})
// }
