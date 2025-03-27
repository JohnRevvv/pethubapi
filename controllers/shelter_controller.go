package controllers

import (
	"errors"
	"pethub_api/middleware"
	"pethub_api/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateShelter creates an shelter account and info
func RegisterShelter(c *fiber.Ctx) error {
	// Parse request body
	requestBody := struct {
		Username           string `json:"username"`
		Password           string `json:"password"`
		ShelterName        string `json:"shelter_name"`
		ShelterAddress     string `json:"shelter_address"`
		ShelterContact     int    `json:"shelter_contact"`
		ShelterOwner       string `json:"shelter_owner"`
		ShelterDescription string `json:"shelter_description"`
	}{}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Check if username exists
	var existingUser models.ShelterAccount
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

	// Create shelter account
	ShelterAccount := models.ShelterAccount{
		Username:  requestBody.Username,
		Password:  requestBody.Password, // Hash the password before storing it
		CreatedAt: time.Now(),
	}

	// Insert into shelteraccount and get the generated ShelterID
	if err := middleware.DBConn.Create(&ShelterAccount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register Shelter",
		})
	}

	// Create shelter info
	ShelterInfo := models.ShelterInfo{
		ShelterID:     ShelterAccount.ID, // Link the ShelterInfo to ShelterAccount
		ShelterName:     requestBody.ShelterName,
		ShelterAddress: requestBody.ShelterAddress,
		ShelterContact: requestBody.ShelterContact,
		ShelterOwner: requestBody.ShelterOwner,
		ShelterDescription: requestBody.ShelterDescription,
	}

	// Insert into Shelterinfo
	if err := middleware.DBConn.Create(&ShelterInfo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register shelter info",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Shelter registered successfully",
		"data": fiber.Map{
			"shelter": ShelterAccount,
			"info":    ShelterInfo,
		},
	})
}

//Login
func LoginShelter(c *fiber.Ctx) error {
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
	var ShelterAccount models.ShelterAccount
	result := middleware.DBConn.Where("username = ?", requestBody.Username).First(&ShelterAccount)

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
	if ShelterAccount.Password != requestBody.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	// Fetch adopter info
	var ShelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", ShelterAccount.ID).First(&ShelterInfo)

	if infoResult.Error != nil && !errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch adopter info",
		})
	}

	// Login successful, return adopter account and info
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data": fiber.Map{
			"Shelter": ShelterAccount,
			"Info":    ShelterInfo,
		},
	})
}