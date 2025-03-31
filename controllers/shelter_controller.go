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
		ShelterLandmark    string `json:"shelter_landmark"`
		ShelterContact     string `json:"shelter_contact"`
		ShelterEmail       string `json:"shelter_email"`
		ShelterOwner       string `json:"shelter_owner"`
		ShelterDescription string `json:"shelter_description"`
		ShelterSocial      string `json:"shelter_social"`
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
		ShelterID:          ShelterAccount.ShelterID, // Link the ShelterInfo to ShelterAccount
		ShelterName:        requestBody.ShelterName,
		ShelterAddress:     requestBody.ShelterAddress,
		ShelterLandmark:    requestBody.ShelterLandmark,
		ShelterContact:     requestBody.ShelterContact,
		ShelterEmail:       requestBody.ShelterEmail,
		ShelterOwner:       requestBody.ShelterOwner,
		ShelterDescription: requestBody.ShelterDescription,
		ShelterSocial:      requestBody.ShelterSocial,
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

	// Check password using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(ShelterAccount.Password), []byte(requestBody.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	// Fetch adopter info
	var ShelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", ShelterAccount.ShelterID).First(&ShelterInfo)

	if infoResult.Error != nil && !errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch adopter info",
		})
	}

	// Login successful, return shelter account, info, and ID
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data": fiber.Map{
			"shelter_id": ShelterAccount.ShelterID, // Include shelter ID in the response
			"Shelter":    ShelterAccount,
			"Info":       ShelterInfo,
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
			if account.ShelterID == info.ShelterID {
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

func GetShelterInfoByID(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Fetch shelter info by ID
	var shelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Fetch shelter account associated with the shelter info
	var shelterAccount models.ShelterAccount
	accountResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterAccount)

	if accountResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter account",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter info retrieved successfully",
		"data": fiber.Map{
			"info": shelterInfo,
		},
	})
}

func GetShelterDetailsByID(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Fetch shelter info by ID
	var shelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo)

	// If no shelter info is found, set it to null
	var infoResponse interface{}
	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		infoResponse = nil
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching shelter info",
		})
	} else {
		infoResponse = shelterInfo
	}

	// Fetch shelter media by ID
	var shelterMedia models.ShelterMedia
	mediaResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterMedia)

	// If no shelter media is found, set it to null
	var mediaResponse interface{}
	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		mediaResponse = nil
	} else if mediaResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching shelter media",
		})
	} else {
		mediaResponse = shelterMedia
	}

	// Combine shelter info and media into a single response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter details retrieved successfully",
		"data": fiber.Map{
			"info":  infoResponse,
			"media": mediaResponse,
		},
	})
}

func UpdateShelterDetails(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Fetch shelter info by ID
	var shelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching shelter info",
		})
	}

	// Parse request body for updated details
	updateData := struct {
		ShelterName        string `json:"shelter_name"`
		ShelterAddress     string `json:"shelter_address"`
		ShelterLandmark    string `json:"shelter_landmark"`
		ShelterContact     string `json:"shelter_contact"`
		ShelterEmail       string `json:"shelter_email"`
		ShelterOwner       string `json:"shelter_owner"`
		ShelterDescription string `json:"shelter_description"`
		ShelterSocial      string `json:"shelter_social"`
		ShelterProfile     string `json:"shelter_profile"` // Binary data for profile
		ShelterCover       string `json:"shelter_cover"`   // Binary data for cover
	}{}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Update shelter info fields
	shelterInfo.ShelterName = updateData.ShelterName
	shelterInfo.ShelterAddress = updateData.ShelterAddress
	shelterInfo.ShelterLandmark = updateData.ShelterLandmark
	shelterInfo.ShelterContact = updateData.ShelterContact
	shelterInfo.ShelterEmail = updateData.ShelterEmail
	shelterInfo.ShelterOwner = updateData.ShelterOwner
	shelterInfo.ShelterDescription = updateData.ShelterDescription
	shelterInfo.ShelterSocial = updateData.ShelterSocial

	// Save updated shelter info to the database
	if err := middleware.DBConn.Save(&shelterInfo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update shelter info",
		})
	}

	// Fetch shelter media by ID
	var shelterMedia models.ShelterMedia
	mediaResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterMedia)

	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		// If no media exists, create a new entry
		shelterMedia = models.ShelterMedia{
			ShelterID:      shelterInfo.ShelterID,
			ShelterProfile: updateData.ShelterProfile,
			ShelterCover:   updateData.ShelterCover,
		}
		if err := middleware.DBConn.Create(&shelterMedia).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create shelter media",
			})
		}
	} else if mediaResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching shelter media",
		})
	} else {
		// Update existing media
		shelterMedia.ShelterProfile = updateData.ShelterProfile
		shelterMedia.ShelterCover = updateData.ShelterCover

		if err := middleware.DBConn.Save(&shelterMedia).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update shelter media",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter details and media updated successfully",
		"data": fiber.Map{
			"info":  shelterInfo,
			"media": shelterMedia,
		},
	})
}
