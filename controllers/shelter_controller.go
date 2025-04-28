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
		Username  string `json:"username"`
		Password  string `json:"password"`
		RegStatus string `json:"reg_status"`
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

	// Check if account is pending
	if ShelterAccount.RegStatus == "pending" {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "403", // Forbidden
			Message: "Your account is pending, please wait for admin approval",
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

	token, err := middleware.GenerateJWT(ShelterAccount.ShelterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Error generating token",
			Data:    err.Error(),
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
			"token":      token,
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
		Data:    shelterMedia,
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

func GetShelterDonationInfo(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	var shelterDonations models.ShelterDonations
	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterDonations)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "Shelter not found",
			Data:    nil,
		})
	} else if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter donation info",
			Data:    infoResult.Error,
		})
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Shelter donation info retrieved successfully",
		Data:    shelterDonations,
	})
}

func UpdateShelterDonations(c *fiber.Ctx) error {
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

	// Check if donation record exists
	var shelterDonations models.ShelterDonations
	result := middleware.DBConn.Where("shelter_id = ?", parsedShelterID).First(&shelterDonations)

	isNew := errors.Is(result.Error, gorm.ErrRecordNotFound)
	if isNew {
		shelterDonations = models.ShelterDonations{ShelterID: uint(parsedShelterID)}
	}

	// Parse text fields from form-data
	shelterDonations.AccountNumber = c.FormValue("account_number", shelterDonations.AccountNumber)
	shelterDonations.AccountName = c.FormValue("account_name", shelterDonations.AccountName)

	// Handle QR image file upload
	qrFile, err := c.FormFile("qr_image")
	if err == nil {
		fileContent, err := qrFile.Open()
		if err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "400",
				Message: "Failed to open QR image",
				Data:    err,
			})
		}
		defer fileContent.Close()

		fileBytes, err := ioutil.ReadAll(fileContent)
		if err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "400",
				Message: "Failed to read QR image",
				Data:    err,
			})
		}
		shelterDonations.QRImage = base64.StdEncoding.EncodeToString(fileBytes)
	}

	// Save to DB
	if isNew {
		if err := middleware.DBConn.Create(&shelterDonations).Error; err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "500",
				Message: "Failed to create shelter donation info",
				Data:    err,
			})
		}
	} else {
		if err := middleware.DBConn.Model(&models.ShelterDonations{}).
			Where("shelter_id = ?", parsedShelterID).
			Updates(map[string]interface{}{
				"account_number": shelterDonations.AccountNumber,
				"account_name":   shelterDonations.AccountName,
				"qr_image":       shelterDonations.QRImage,
			}).Error; err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "500",
				Message: "Failed to update shelter donation info",
				Data:    err,
			})
		}
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Shelter donation info saved successfully",
		Data:    shelterDonations,
	})
}

func ShelterChangePassword(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	var shelterAccount models.ShelterAccount
	result := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterAccount)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Shelter account not found",
			Data:    nil,
		})
	}
	if result.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter account",
			Data:    nil,
		})
	}

	// Parse request body with confirm password
	requestBody := struct {
		OldPassword     string `json:"old_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	// Check if new and confirm passwords match
	if requestBody.NewPassword != requestBody.ConfirmPassword {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "New password and confirm password do not match",
			Data:    nil,
		})
	}

	// Check if old password matches
	err := bcrypt.CompareHashAndPassword([]byte(shelterAccount.Password), []byte(requestBody.OldPassword))
	if err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Old password is incorrect",
			Data:    nil,
		})
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to hash new password",
			Data:    nil,
		})
	}

	// Update the password in the database
	shelterAccount.Password = string(hashedPassword)
	if err := middleware.DBConn.Save(&shelterAccount).Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to update password",
			Data:    nil,
		})
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Password updated successfully",
		Data:    nil,
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

func CountPetsByShelter(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Check if shelter exists
	var shelterInfo models.ShelterInfo
	if err := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Shelter not found",
		})
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

func GetSingleAdoptionApplication(c *fiber.Ctx) error {
	shelterID := c.Query("shelter_id") // Assuming shelter_id is passed as a query
	status := c.Query("status")

	if status == "" || shelterID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "pet_id and shelter_id are required",
		})
	}

	var applications []models.AdoptionSubmission
	if err := middleware.DBConn.Debug().Where("shelter_id = ? AND status = ?", shelterID, status).Preload("Pet.PetMedia").Preload("Adopter").Find(&applications).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Failed to fetch adoption applications",
			Data:    err.Error(),
		})
	}

	return c.JSON(applications)
}

func GetAdoptionApplications(c *fiber.Ctx) error {
	shelterID := c.Params("id")
	status := c.Query("status")

	if status == "" || shelterID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "shelter_id and status are required",
		})
	}

	// Create a custom struct just for the response
	type AdoptionApplicationResponse struct {
		ApplicationID  uint   `json:"application_id"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		AdopterProfile string `json:"adopter_profile"`
		PetName        string `json:"pet_name"`
		Status         string `json:"status"`
	}

	var applications []models.AdoptionSubmission
	var responses []AdoptionApplicationResponse

	// Fetch the adoption submissions with related data
	if err := middleware.DBConn.Debug().
		Where("shelter_id = ? AND status = ?", shelterID, status).
		Preload("Adopter").
		Preload("Adopter.AdopterMedia"). // Preload adopter media
		Preload("Pet").                  // Preload pet data
		Find(&applications).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Failed to fetch adoption applications",
			Data:    err.Error(),
		})
	}

	// Format the response to match the required fields
	for _, app := range applications {
		responses = append(responses, AdoptionApplicationResponse{
			ApplicationID:  app.ApplicationID,
			FirstName:      app.Adopter.FirstName,
			LastName:       app.Adopter.LastName,
			AdopterProfile: app.Adopter.AdopterMedia.AdopterProfile, // Assuming `ProfileURL` is the field for the profile image
			PetName:        app.Pet.PetName,                         // Assuming `PetName` is the pet's name field
			Status:         app.Status,
		})
	}

	return c.JSON(responses)
}


func GetApplicationByApplicationID(c *fiber.Ctx) error {
	applicationID := c.Params("application_id")

	var adoptionSubmission models.AdoptionSubmission
	infoResult := middleware.DBConn.Debug().Where("application_id = ?", applicationID).First(&adoptionSubmission)

	if infoResult.Error != nil {
		if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
			return c.JSON(response.AdopterResponseModel{
				RetCode: "404",
				Message: "Application not found",
				Data:    nil,
			})
		}
		return c.JSON(response.AdopterResponseModel{
			RetCode: "500",
			Message: "Something went wrong",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adoption details retrieved successfully",
		"data": fiber.Map{
			"info": adoptionSubmission,
		},
	})
}
