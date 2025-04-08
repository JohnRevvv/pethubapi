package controllers

// import (
// 	"errors"
// 	"pethub_api/middleware"
// 	"pethub_api/models"
// 	"time"
// 	"pethub_api/models/response"

// 	"github.com/gofiber/fiber/v2"
// 	"golang.org/x/crypto/bcrypt"
// 	"gorm.io/gorm"
// )

// // CreateShelter creates an shelter account and info
// func RegisterShelter(c *fiber.Ctx) error {
// 	// Parse request body
// 	requestBody := struct {
// 		Username           string `json:"username"`
// 		Password           string `json:"password"`
// 		ShelterName        string `json:"shelter_name"`
// 		ShelterAddress     string `json:"shelter_address"`
// 		ShelterLandmark    string `json:"shelter_landmark"`
// 		ShelterContact     string `json:"shelter_contact"`
// 		ShelterEmail       string `json:"shelter_email"`
// 		ShelterOwner       string `json:"shelter_owner"`
// 		ShelterDescription string `json:"shelter_description"`
// 		ShelterSocial      string `json:"shelter_social"`
// 	}{}

// 	if err := c.BodyParser(&requestBody); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Invalid request body",
// 		})
// 	}

// 	// Check if username exists
// 	var existingUser models.ShelterAccount
// 	result := middleware.DBConn.Where("username = ?", requestBody.Username).First(&existingUser)
// 	if result.Error == nil {
// 		return c.JSON(response.ShelterResponseModel{
// 			RetCode: "400",
// 			Message: "Username already exists!",
// 			Data:    nil,
// 		})
// 	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return c.JSON(response.ShelterResponseModel{
// 			RetCode: "500",
// 			Message: "Database error",
// 			Data:    result.Error,
// 		})
// 	}

// 	// Hash the password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return c.JSON(response.ShelterResponseModel{
// 			RetCode: "500",
// 			Message: "Failed to hash password",
// 			Data:    err,
// 		})
// 	}

// 	// Create shelter account
// 	ShelterAccount := models.ShelterAccount{
// 		Username:  requestBody.Username,
// 		Password:  string(hashedPassword), // Store hashed password
// 		CreatedAt: time.Now(),
// 	}

// 	// Insert into shelteraccount and get the generated ShelterID
// 	if err := middleware.DBConn.Create(&ShelterAccount).Error; err != nil {
// 		return c.JSON(response.ShelterResponseModel{
// 			RetCode: "500",
// 			Message: "Failed to Register Shelter Account",
// 			Data:    err,
// 		})
// 	}

// 	// Create shelter info
// 	ShelterInfo := models.ShelterInfo{
// 		ShelterID:          ShelterAccount.ShelterID, // Link the ShelterInfo to ShelterAccount
// 		ShelterName:        requestBody.ShelterName,
// 		ShelterAddress:     requestBody.ShelterAddress,
// 		ShelterLandmark:    requestBody.ShelterLandmark,
// 		ShelterContact:     requestBody.ShelterContact,
// 		ShelterEmail:       requestBody.ShelterEmail,
// 		ShelterOwner:       requestBody.ShelterOwner,
// 		ShelterDescription: requestBody.ShelterDescription,
// 		ShelterSocial:      requestBody.ShelterSocial,
// 	}

// 	// Insert into Shelterinfo
// 	if err := middleware.DBConn.Create(&ShelterInfo).Error; err != nil {
// 		return c.JSON(response.ShelterResponseModel{
// 			RetCode: "500",
// 			Message: "Failed to Register Shelter Info",
// 			Data:    err,
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "Shelter registered successfully",
// 		"data": fiber.Map{
// 			"shelter": ShelterAccount,
// 			"info":    ShelterInfo,
// 		},
// 	})
// }
