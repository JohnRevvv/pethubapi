package controllers

import (
	"errors"
	"pethub_api/middleware"
	"pethub_api/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	// Create shelter account
	ShelterAccount := models.ShelterAccount{
		Username:  requestBody.Username,
		Password:  string(hashedPassword), // Store hashed password
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
		ShelterID:          ShelterAccount.ID, // Link the ShelterInfo to ShelterAccount
		ShelterName:        requestBody.ShelterName,
		ShelterAddress:     requestBody.ShelterAddress,
		ShelterContact:     requestBody.ShelterContact,
		ShelterOwner:       requestBody.ShelterOwner,
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

// Login
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

func GetAllShelters(c *fiber.Ctx) error {
	// Fetch all shelter accounts
	var shelterAccounts []models.ShelterAccount
	accountResult := middleware.DBConn.Find(&shelterAccounts)

	if accountResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter accounts",
		})
	}

	// Fetch all shelter info
	var shelterInfos []models.ShelterInfo
	infoResult := middleware.DBConn.Find(&shelterInfos)

	if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter info",
		})
	}

	// Combine accounts and info into a single response
	shelters := []fiber.Map{}
	for _, account := range shelterAccounts {
		for _, info := range shelterInfos {
			if account.ID == info.ShelterID {
				shelters = append(shelters, fiber.Map{
					"shelter": account,
					"info":    info,
				})
				break
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelters retrieved successfully",
		"data":    shelters,
	})
}

func GetShelterByName(c *fiber.Ctx) error {
	shelterName := c.Query("shelter_name")
	if shelterName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Shelter name query parameter is missing",
		})
	}

	// Fetch shelter info by name
	var shelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_name = ?", shelterName).First(&shelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Fetch shelter account associated with the shelter info
	var shelterAccount models.ShelterAccount
	accountResult := middleware.DBConn.Where("id = ?", shelterInfo.ShelterID).First(&shelterAccount)

	if accountResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter account",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter retrieved successfully",
		"data": fiber.Map{
			"shelter": shelterAccount,
			"info":    shelterInfo,
		},
	})
}
