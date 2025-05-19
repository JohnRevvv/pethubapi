package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"pethub_api/middleware"
	"pethub_api/models"
	"pethub_api/models/response"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetShelterInfo(c *fiber.Ctx) error {
	ShelterID := c.Params("shelter_id")

	var ShelterInfo models.ShelterInfo
	Result := middleware.DBConn.Debug().Preload("ShelterMedia").Where("shelter_id = ?", ShelterID).First(&ShelterInfo)

	if errors.Is(Result.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "Shelter Not Found",
			Data:    nil,
		})
	} else if Result.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter info",
			Data:    Result.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    ShelterInfo,
	})
}

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
			RetCode: "404",
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

func GetShelterDetailsRefined(c *fiber.Ctx) error {
	ShelterID := c.Params("shelter_id")

	var ShelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Debug().Preload("ShelterMedia").Where("shelter_id = ?", ShelterID).First(&ShelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "Shelter Not Found",
			Data:    nil,
		})
	} else if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching shelter info",
			Data:    infoResult.Error,
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    ShelterInfo,
	})
}

// not use
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

//old
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
		CreatedAt      string `json:"created_at"`


	}

	var applications []models.AdoptionSubmission
	var responses []AdoptionApplicationResponse

	// Fetch the adoption submissions with related data
	if err := middleware.DBConn.Debug().
		Where("shelter_id = ? AND status = ?", shelterID, status).
		Preload("Adopter").
		Preload("Adopter.AdopterMedia"). // Preload adopter media
		Preload("Pet").                  // Preload pet data
		Preload("ScheduleInterview").
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
			CreatedAt:      app.CreatedAt.Format(time.RFC3339), // Format the date as needed
		})
	}

	return c.JSON(responses)
}
//new
func GetAdoptionSubmissionsByShelterAndStatus(c *fiber.Ctx) error {
	shelterID := c.Params("shelter_id")
	status := c.Query("status") // Example: ?status=approved

	var submissions []models.AdoptionSubmission
	result := middleware.DBConn.Debug().
		Preload("Adopter").
		Preload("Adopter.AdopterMedia").
		Preload("Pet").
		Preload("Pet.PetMedia").
		Preload("ScheduleInterview").
		Where("shelter_id = ? AND status = ?", shelterID, status).
		Find(&submissions)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching submissions",
			Data:    result.Error.Error(),
		})
	}

	if len(submissions) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "No adoption submissions found",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data": fiber.Map{
			"submissions": submissions,
		},
	})
}



func GetApplicationByApplicationID(c *fiber.Ctx) error {
	applicationID := c.Params("application_id")

	var adoptionSubmission models.AdoptionSubmission
	infoResult := middleware.DBConn.Debug().Where("application_id = ?", applicationID).
		Preload("Adopter").
		Preload("Adopter.AdopterMedia").
		Preload("Pet").
		Preload("Pet.PetMedia").
		First(&adoptionSubmission)

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

	// Manually fetch the ApplicationPhotos using ImageID
	var appPhotos models.ApplicationPhotos
	photoResult := middleware.DBConn.Debug().Where("image_id = ?", adoptionSubmission.ImageID).First(&appPhotos)
	if photoResult.Error != nil && !errors.Is(photoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "500",
			Message: "Failed to retrieve application photos",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adoption details retrieved successfully",
		"data": fiber.Map{
			"info":              adoptionSubmission,
			"applicationPhotos": appPhotos,
		},
	})
}

func SetInterviewSchedule(c *fiber.Ctx) error {
	applicationID := c.Params("application_id")

	var application models.AdoptionSubmission
	Result := middleware.DBConn.Debug().Where("application_id = ?", applicationID).First(&application)
	if Result.Error != nil {
		if errors.Is(Result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "404",
				Message: "Application not found",
				Data:    nil,
			})
		}
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    Result.Error,
		})
	}

	// Check if interview schedule already exists for this application
	var existingInterview models.ScheduleInterview
	if err := middleware.DBConn.Where("application_id = ?", application.ApplicationID).First(&existingInterview).Error; err == nil {
		// Found existing interview schedule
		return c.JSON(response.ShelterResponseModel{
			RetCode: "409",
			Message: "Interview schedule already exists for this application",
			Data:    existingInterview,
		})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other DB error
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error checking existing interview",
			Data:    err.Error(),
		})
	}

	// Update the selected application status to 'interview'
	application.Status = "interview"
	if err := middleware.DBConn.Save(&application).Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to update application status",
			Data:    err,
		})
	}

	// Set all other applications for the same pet to 'in queue'
	if err := middleware.DBConn.Model(&models.AdoptionSubmission{}).
		Where("pet_id = ? AND application_id != ?", application.PetID, application.ApplicationID).
		Update("status", "in queue").Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to update other applications",
			Data:    err,
		})
	}

	// Update pet status to 'pending'
	if err := middleware.DBConn.Model(&models.PetInfo{}).
		Where("pet_id = ?", application.PetID).
		Update("status", "pending").Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to update pet status",
			Data:    err,
		})
	}

	// Parse JSON input
	type InterviewInput struct {
		InterviewDate  string `json:"interview_date"`  // Format: YYYY-MM-DD
		InterviewTime  string `json:"interview_time"`  // Format: HH:MM:SS
		InterviewNotes string `json:"interview_notes"`
	}

	var input InterviewInput
	if err := c.BodyParser(&input); err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Invalid JSON input",
			Data:    err.Error(),
		})
	}

	// Convert strings to time.Time
	interviewDate, err := time.Parse("2006-01-02", input.InterviewDate)
	if err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Invalid interview_date format. Use YYYY-MM-DD",
		})
	}

	interviewTime, err := time.Parse("15:04:05", input.InterviewTime)
	if err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Invalid interview_time format. Use HH:MM:SS",
		})
	}

	// Save interview schedule
	newInterview := models.ScheduleInterview{
		ApplicationID:  application.ApplicationID,
		ShelterID:      application.ShelterID,
		AdopterID:      application.AdopterID,
		InterviewDate:  interviewDate,
		InterviewTime:  interviewTime.Format("15:04:05"), // Convert time.Time to string
		InterviewNotes: input.InterviewNotes,
		CreatedAt:      time.Now(),
	}

	if err := middleware.DBConn.Create(&newInterview).Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to create interview schedule",
			Data:    err,
		})
	}


	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Interview scheduled. Application marked as 'interview'. Others set to 'in queue'. Pet status set to 'pending'.",
		Data:    newInterview,
	})
}


func RejectApplication(c *fiber.Ctx) error {
	applicationID := c.Params("application_id")

	// Parse incoming JSON body
	var body struct {
		ReasonForRejection []string `json:"reason_for_rejection"` // Accept multiple rejection reasons
	}
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	// Fetch the application
	var application models.AdoptionSubmission
	result := middleware.DBConn.Debug().Preload("ScheduleInterview").Where("application_id = ?", applicationID).First(&application)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "404",
				Message: "Application not found",
				Data:    nil,
			})
		}
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    result.Error,
		})
	}

	// Combine reasons into one string (e.g., comma-separated)
	reasonStr := strings.Join(body.ReasonForRejection, ", ")

	// Set status and reasons
	application.Status = "rejected"
	application.ReasonForRejection = reasonStr

	if err := middleware.DBConn.Save(&application).Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to archive application",
			Data:    err,
		})
	}

	log.Println("RejectApplication called with ID:", applicationID)
log.Println("Body:", string(c.Body()))

	// Update interview status to 'rejected'
	if err := middleware.DBConn.Model(&models.ScheduleInterview{}).
		Where("application_id = ?", applicationID).
		Update("interview_status", "rejected").Error; err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Failed to update interview status",
			Data:    err,
		})
	}

	// Update pet status if it's 'pending'
	var pet models.PetInfo
	if err := middleware.DBConn.Debug().Where("pet_id = ?", application.PetID).First(&pet).Error; err == nil {
		if pet.Status == "pending" {
			pet.Status = "available"
			if err := middleware.DBConn.Save(&pet).Error; err != nil {
				return c.JSON(response.ShelterResponseModel{
					RetCode: "500",
					Message: "Failed to update pet status",
					Data:    err,
				})
			}
		}
	}

	// Update other in-queue applications to pending
	var count int64
	middleware.DBConn.Model(&models.AdoptionSubmission{}).
		Where("pet_id = ? AND application_id != ? AND status = ?", application.PetID, application.ApplicationID, "in queue").
		Count(&count)

	if count > 0 {
		if err := middleware.DBConn.Model(&models.AdoptionSubmission{}).
			Where("pet_id = ? AND application_id != ? AND status = ?", application.PetID, application.ApplicationID, "in queue").
			Update("status", "pending").Error; err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "500",
				Message: "Failed to update other applications",
				Data:    err,
			})
		}
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Application rejected. Reasons saved. Interview rejected. Pet and other apps updated.",
		Data:    application,
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

	var available, unavailable, adopted int64

	middleware.DBConn.Debug().Model(&models.PetInfo{}).Where("shelter_id = ? AND status = ?", shelterID, "available").Count(&available)
	middleware.DBConn.Debug().Model(&models.PetInfo{}).Where("shelter_id = ? AND status = ?", shelterID, "unavailable").Count(&unavailable)
	middleware.DBConn.Debug().Model(&models.PetInfo{}).Where("shelter_id = ? AND status = ?", shelterID, "adopted").Count(&adopted)

	return c.Status(200).JSON(fiber.Map{
		"message": "Counts fetched successfully",
		"data": fiber.Map{
			"shelter_id":  shelterID,
			"available":   available,
			"unavailable": unavailable,
			"adopted":     adopted,
		},
	})
}

func CountApplicantsByPetId(c *fiber.Ctx) error {
	petId := c.Params("pet_id")

	var count int64
	result := middleware.DBConn.Debug().
		Model(&models.AdoptionSubmission{}).
		Where("pet_id = ?", petId).
		Count(&count)

	if result.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    result.Error.Error(),
		})
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Count fetched successfully",
		Data:    count,
	})
}

func GetPetsWithAdoptionRequestsByShelter(c *fiber.Ctx) error {
	shelterID := c.Params("shelter_id")

	var submissions []models.AdoptionSubmission

	result := middleware.DBConn.Debug().
		Joins("JOIN petinfo ON petinfo.pet_id = adoption_submissions.pet_id").
		Where("adoption_submissions.shelter_id = ? AND petinfo.status = ? AND adoption_submissions.status = ?", shelterID, "available", "pending").
		Preload("Pet").
		Preload("Pet.PetMedia").
		Find(&submissions)

	if result.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    result.Error.Error(),
		})
	}

	// Remove duplicate pets
	petMap := make(map[uint]models.PetInfo)
	for _, sub := range submissions {
		if sub.Pet.PetID != 0 {
			petMap[sub.Pet.PetID] = sub.Pet
		}
	}

	// Convert map to slice
	pets := make([]models.PetInfo, 0)
	for _, pet := range petMap {
		pets = append(pets, pet)
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Pets with adoption requests fetched successfully",
		Data:    pets,
	})
}

func GetAdoptionApplicationsByPetID(c *fiber.Ctx) error {
	petID := c.Params("pet_id") // This is the pet_id

	// Check if pet_id is provided
	if petID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "pet_id is required",
		})
	}

	// Create a custom struct just for the response
	type AdoptionApplicationResponse struct {
		ApplicationID  uint   `json:"application_id"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		AdopterProfile string `json:"adopter_profile"`
		PetName        string `json:"pet_name"`
		Address        string `json:"address"`
		ContactNumber  string `json:"contact_number"`
		Email		  string `json:"email"`
		Status         string `json:"status"`
		CreatedAt      string `json:"created_at"`
	}

	var applications []models.AdoptionSubmission
	var responses []AdoptionApplicationResponse

	// Fetch the adoption submissions for the given pet_id and optional status
	query := middleware.DBConn.Debug().
		Where("pet_id = ?", petID). // Filter by pet_id
		Order("created_at").
		Preload("Adopter").              // Preload adopter data
		Preload("Adopter.AdopterMedia"). // Preload adopter media
		Preload("Pet").                  // Preload pet data
		Find(&applications)

	if err := query.Error; err != nil {
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
			Address: 	  app.Adopter.Address,
			ContactNumber:  app.Adopter.ContactNumber,
			Email: 		app.Adopter.Email,
			AdopterProfile: app.Adopter.AdopterMedia.AdopterProfile, // Assuming `AdopterProfile` is the correct field
			PetName:        app.Pet.PetName,                         // Assuming `PetName` is the pet's name field
			Status:         app.Status,
			CreatedAt:      app.CreatedAt.Format(time.RFC3339), // Format the date as needed
		})
	}

	return c.JSON(responses)
}
