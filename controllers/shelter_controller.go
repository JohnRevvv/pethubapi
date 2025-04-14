package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"pethub_api/middleware"
	"pethub_api/models"
	"pethub_api/models/response"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Username already exists!",
			Data:    nil,
		})
	} else if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    result.Error,
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to hash password",
			Data:    err,
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
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to Register Shelter Account",
			Data:    err,
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
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to Register Shelter Info",
			Data:    err,
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

func LoginShelter(c *fiber.Ctx) error {
	// Parse request body
	requestBody := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400", // Bad request
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// Check if the adopter exists
	var ShelterAccount models.ShelterAccount
	result := middleware.DBConn.Where("username = ?", requestBody.Username).First(&ShelterAccount)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400", // Bad request
			Message: "Invalid username or password",
			Data:    nil,
		})
	} else if result.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500", // Internal server error
			Message: "Database error",
			Data:    nil,
		})
	}

	// Check password using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(ShelterAccount.Password), []byte(requestBody.Password))
	if err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400", // Bad request
			Message: "Invalid username or password",
			Data:    nil,
		})
	}

	// Fetch shelter info
	var ShelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", ShelterAccount.ShelterID).First(&ShelterInfo)

	if infoResult.Error != nil && !errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500", // Internal server error
			Message: "Failed to fetch shelter info",
			Data:    nil,
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

func UpdateShelterDetails(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Fetch shelter info by ID
	var shelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Shelter info not found",
			Data:    nil,
		})
	} else if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter info",
			Data:    nil,
		})
	}

	// Parse JSON body for shelter info updates
	var updateRequest models.ShelterInfo
	if err := c.BodyParser(&updateRequest); err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400", // Bad request
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// Convert struct to a map for updating only non-empty fields
	updateData := map[string]interface{}{}
	if updateRequest.ShelterName != "" {
		updateData["shelter_name"] = updateRequest.ShelterName
	}
	if updateRequest.ShelterAddress != "" {
		updateData["shelter_address"] = updateRequest.ShelterAddress
	}
	if updateRequest.ShelterLandmark != "" {
		updateData["shelter_landmark"] = updateRequest.ShelterLandmark
	}
	if updateRequest.ShelterContact != "" {
		updateData["shelter_contact"] = updateRequest.ShelterContact
	}
	if updateRequest.ShelterEmail != "" {
		updateData["shelter_email"] = updateRequest.ShelterEmail
	}
	if updateRequest.ShelterOwner != "" {
		updateData["shelter_owner"] = updateRequest.ShelterOwner
	}
	if updateRequest.ShelterDescription != "" {
		updateData["shelter_description"] = updateRequest.ShelterDescription
	}
	if updateRequest.ShelterSocial != "" {
		updateData["shelter_social"] = updateRequest.ShelterSocial
	}

	// Debugging log
	fmt.Printf("Update Data: %+v\n", updateData)

	if len(updateData) == 0 {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "No fields to update",
			Data:    nil,
		})
	}

	// Update shelter info fields
	if err := middleware.DBConn.Model(&shelterInfo).Updates(updateData).Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400", // Bad request
			Message: "Failed to update shelter info",
			Data:    nil,
		})
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Shelter details updated successfully",
		Data:    shelterInfo,
	})
}

func UploadShelterMedia(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Parse shelter ID
	parsedShelterID, err := strconv.ParseUint(shelterID, 10, 32)
	if err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Invalid shelter ID",
			Data:    err,
		})
	}
	// Fetch existing shelter media or create a new one
	var shelterMedia models.ShelterMedia
	mediaResult := middleware.DBConn.Where("shelter_id = ?", parsedShelterID).First(&shelterMedia)

	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		// Create a new shelter media record if not found
		shelterMedia = models.ShelterMedia{ShelterID: uint(parsedShelterID)}
		middleware.DBConn.Create(&shelterMedia)
	} else if mediaResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter media",
			Data:    err,
		})
	}

	// Handle profile image upload
	profileFile, err := c.FormFile("shelter_profile")
	if err == nil {
		fileContent, err := profileFile.Open()
		if err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "400",
				Message: "Failed to open profile image",
				Data:    err,
			})
		}
		defer fileContent.Close()

		fileBytes, err := ioutil.ReadAll(fileContent)
		if err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "400",
				Message: "Failed to read profile image",
				Data:    err,
			})
		}
		shelterMedia.ShelterProfile = base64.StdEncoding.EncodeToString(fileBytes) // Replace old image
	}

	// Handle cover image upload
	coverFile, err := c.FormFile("shelter_cover")
	if err == nil {
		fileContent, err := coverFile.Open()
		if err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "500",
				Message: "Database error while fetching shelter media",
				Data:    err,
			})
		}
		defer fileContent.Close()

		fileBytes, err := ioutil.ReadAll(fileContent)
		if err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "400",
				Message: "Failed to read cover image",
				Data:    err,
			})
		}
		shelterMedia.ShelterCover = base64.StdEncoding.EncodeToString(fileBytes) // Replace old image
	}

	// Explicitly update fields with WHERE condition
	updateResult := middleware.DBConn.Model(&models.ShelterMedia{}).
		Where("shelter_id = ?", parsedShelterID).
		Updates(map[string]interface{}{
			"shelter_profile": shelterMedia.ShelterProfile,
			"shelter_cover":   shelterMedia.ShelterCover,
		})

	if updateResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Failed to update shelter media",
			Data:    updateResult.Error.Error(),
		})
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Shelter media uploaded/updated successfully",
		Data:  shelterMedia,
	})
}

func GetShelterInfoByID(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Fetch shelter info by ID
	var shelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Shelter info not found",
			Data:    nil,
		})
	} else if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter info",
			Data:    nil,
		})
	}

	// Fetch shelter media by ID
	var shelterMedia models.ShelterMedia
	mediaResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterMedia)

	// If no shelter media is found, set it to null
	var mediaResponse interface{}
	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		mediaResponse = nil
	} else if mediaResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter media",
			Data:    nil,
		})
	} else {
		// Decode Base64-encoded images
		decodedProfile, err := base64.StdEncoding.DecodeString(shelterMedia.ShelterProfile)
		if err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "500",
				Message: "Failed to decode profile image",
				Data:    err,
			})
		}

		decodedCover, err := base64.StdEncoding.DecodeString(shelterMedia.ShelterCover)
		if err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "500",
				Message: "Failed to decode cover image",
				Data:    err,
			})
		}

		// Include decoded images in the response
		mediaResponse = fiber.Map{
			"shelter_profile": decodedProfile,
			"shelter_cover":   decodedCover,
		}
	}

	// Combine shelter info and media into a single response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter info retrieved successfully",
		"data": fiber.Map{
			"info":  shelterInfo,
			"media": mediaResponse,
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
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter info",
			Data:    nil,
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
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Database error while fetching shelter media",
			Data:    nil,
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

func GetAllShelters(c *fiber.Ctx) error {
	// Fetch all shelter accounts
	var shelterAccounts []models.ShelterAccount
	accountResult := middleware.DBConn.Find(&shelterAccounts)

	if accountResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Failed to fetch shelter accounts",
			Data:    nil,
		})
	}

	// Fetch all shelter info
	var shelterInfos []models.ShelterInfo
	infoResult := middleware.DBConn.Find(&shelterInfos)

	if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Failed to fetch shelter info",
			Data:    nil,
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

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "All shelters retrieved successfully",
		Data:    shelters,
	})
}

func GetShelterByName(c *fiber.Ctx) error {
	shelterName := c.Query("shelter_name")
	if shelterName == "" {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Shelter name is required",
			Data:    nil,
		})
	}

	// Fetch shelter info by name
	var shelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_name = ?", shelterName).First(&shelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Shelter info not found",
			Data:    nil,
		})
	} else if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter info",
			Data:    nil,
		})
	}

	// Fetch shelter account associated with the shelter info
	var shelterAccount models.ShelterAccount
	accountResult := middleware.DBConn.Where("id = ?", shelterInfo.ShelterID).First(&shelterAccount)

	if accountResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Failed to fetch shelter account",
			Data:    nil,
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


//For Data Analytics
func CountPetsByShelter(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Check if shelter exists
	var shelterInfo models.ShelterInfo
	if err := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Shelter not found"})
	}

	var available, pending, adopted int64

	middleware.DBConn.Debug().Model(&models.PetInfo{}).Where("shelter_id = ? AND status = ?", shelterID, "available").Count(&available)
	middleware.DBConn.Model(&models.PetInfo{}).Where("shelter_id = ? AND status = ?", shelterID, "pending").Count(&pending)
	middleware.DBConn.Model(&models.PetInfo{}).Where("shelter_id = ? AND status = ?", shelterID, "adopted").Count(&adopted)

	return c.Status(200).JSON(fiber.Map{
		"message": "Counts fetched successfully",
		"data": fiber.Map{
			"shelter_id": shelterID,
			"available":  available,
			"pending":    pending,
			"adopted":    adopted,
		},
	})
}
