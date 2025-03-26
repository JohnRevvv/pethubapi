package controllers

import (
	"errors"
	"time"

	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateAdopter creates an adopter account and info
func RegisterAdopter(c *fiber.Ctx) error {
	// Parse request body
	requestBody := struct {
		Username      string `json:"username"`
		Password      string `json:"password"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Age           int    `json:"age"`
		Sex           string `json:"sex"`
		Address       string `json:"address"`
		ContactNumber int `json:"contact_number"`
		Email         string `json:"email"`
		Occupation    string `json:"occupation"`
		CivilStatus   string `json:"civil_status"`
		SocialMedia   string `json:"social_media"`
	}{}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Check if username exists
	var existingUser models.AdopterAccount
	result := middleware.DBConn.Where("username = ?", requestBody.Username).First(&existingUser)
	if result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already exists",
		})
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Create adopter account
	adopterAccount := models.AdopterAccount{
		Username: requestBody.Username,
		Password: requestBody.Password, // Hash the password before storing it
		CreatedAt: time.Now(),
	}

	// Insert into adopteraccount and get the generated AdopterID
	if err := middleware.DBConn.Create(&adopterAccount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register adopter",
		})
	}

	// Create adopter info
	adopterInfo := models.AdopterInfo{
		AdopterID:     adopterAccount.AdopterID, // Link the AdopterInfo to AdopterAccount
		FirstName:     requestBody.FirstName,
		LastName:      requestBody.LastName,
		Age:           requestBody.Age,
		Sex:           requestBody.Sex,
		Address:       requestBody.Address,
		ContactNumber: requestBody.ContactNumber,
		Email:         requestBody.Email,
		Occupation:    requestBody.Occupation,
		CivilStatus:   requestBody.CivilStatus,
		SocialMedia:   requestBody.SocialMedia,
	}

	// Insert into adopterinfo
	if err := middleware.DBConn.Create(&adopterInfo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register adopter info",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Adopter registered successfully",
		"data": fiber.Map{
			"adopter": adopterAccount,
			"info":    adopterInfo,
		},
	})
}


//Login
// LoginAdopter authenticates an adopter and retrieves their info
func LoginAdopter(c *fiber.Ctx) error {
	// Parse request body
	requestBody := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Check if the adopter exists
	var adopterAccount models.AdopterAccount
	result := middleware.DBConn.Where("username = ?", requestBody.Username).First(&adopterAccount)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	} else if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Check password (assuming password is stored as plain text for now)
	if adopterAccount.Password != requestBody.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	// Fetch adopter info
	var adopterInfo models.AdopterInfo
	infoResult := middleware.DBConn.Where("adopter_id = ?", adopterAccount.AdopterID).First(&adopterInfo)

	if infoResult.Error != nil && !errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch adopter info",
		})
	}

	// Login successful, return adopter account and info
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data": fiber.Map{
			"adopter": adopterAccount,
			"info":    adopterInfo,
		},
	})
}


