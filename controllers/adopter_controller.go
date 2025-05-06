package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"pethub_api/middleware"
	"pethub_api/models"
	"pethub_api/models/response"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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
		ContactNumber string `json:"contact_number"`
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

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	// Create adopter account
	adopterAccount := models.AdopterAccount{
		Username:  requestBody.Username,
		Password:  string(hashedPassword), // Store hashed password
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
		AdopterID:     adopterAccount.AdopterID,
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

	// Check password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(adopterAccount.Password), []byte(requestBody.Password)); err != nil {
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

	// Generate JWT
	token, err := middleware.GenerateJWT(adopterAccount.AdopterID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("Failed to generate token: %v", err),
		})
	}

	// Return JWT along with account and info
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"data": fiber.Map{
			"adopter": adopterAccount,
			"info":    adopterInfo,
		},
	})
}

func GetAllAdopters(c *fiber.Ctx) error {
	// Fetch all adopter accounts
	var adopterAccounts []models.AdopterAccount
	accountResult := middleware.DBConn.Find(&adopterAccounts)

	if accountResult.Error != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "400",
			Message: "Failed to fetch adopter info",
			Data:    nil,
		})
	}

	// Fetch all adopter info
	var adopterInfos []models.AdopterInfo
	infoResult := middleware.DBConn.Preload("AdopterMedia").Find(&adopterInfos)

	if infoResult.Error != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "400",
			Message: "Failed to fetch adopter info",
			Data:    nil,
		})
	}

	// Combine accounts and info into a single response
	adopters := []fiber.Map{}
	for _, account := range adopterAccounts {
		for _, info := range adopterInfos {
			if account.AdopterID == info.AdopterID {
				adopters = append(adopters, fiber.Map{
					"adopter": account,
					"info":    info,
				})
				break
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adopters retrieved successfully",
		"data":    adopters,
	})
}

func GetAdopterInfoByID(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	// Fetch shelter info by ID
	var adopterInfo models.AdopterInfo
	infoResult := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "404",
			Message: "Adopter not found",
			Data:    nil,
		})
	} else if infoResult.Error != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    nil,
		})
	}

	// Fetch adopter media by ID
	var adopterMedia models.AdopterMedia
	mediaResult := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopterMedia)

	// If no adopter media is found, set it to null
	var mediaResponse interface{}
	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		mediaResponse = nil
	} else if mediaResult.Error != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    nil,
		})
	} else {
		// Decode Base64-encoded images
		decodedProfile, err := base64.StdEncoding.DecodeString(adopterMedia.AdopterProfile)
		if err != nil {
			return c.JSON(response.AdopterResponseModel{
				RetCode: "400",
				Message: "Failed to decode profile image",
				Data:    nil,
			})
		}

		// Include decoded images in the response
		mediaResponse = fiber.Map{
			"adopter_profile": decodedProfile,
		}
	}

	// Combine shelter info and media into a single response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adopter info retrieved successfully",
		"data": fiber.Map{
			"info":  adopterInfo,
			"media": mediaResponse,
		},
	})
}

func AddPetAdoption(c *fiber.Ctx) error {
	// Get adopter ID from route param
	adopterIDParam := c.Params("id")
	adopterID, err := strconv.ParseUint(adopterIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid adopter ID",
		})
	}

	// Parse pet_id from JSON body
	var body struct {
		PetID uint `json:"pet_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate Adopter exists
	var adopter models.AdopterInfo
	if err := middleware.DBConn.First(&adopter, uint(adopterID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Adopter not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching adopter",
		})
	}

	// Validate Pet exists
	var pet models.PetInfo
	if err := middleware.DBConn.First(&pet, body.PetID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Pet not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching pet",
		})
	}

	// Save to AdoptedPets table (prevent duplicate entry)
	adoptedPet := models.AdoptedPet{
		AdopterID: adopter.AdopterID,
		PetID:     pet.PetID,
	}

	if err := middleware.DBConn.
		Where("adopter_id = ? AND pet_id = ?", adoptedPet.AdopterID, adoptedPet.PetID).
		FirstOrCreate(&adoptedPet).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to record adoption",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pet successfully adopted",
		"data":    adoptedPet,
	})
}

func GetPetsByAdopterID(c *fiber.Ctx) error {
	// Retrieve the adopter ID from the URL parameter
	idParam := c.Params("id")
	if idParam == "" {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "404",
			Message: "ID not found",
			Data:    nil,
		})
	}

	adopterID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "400",
			Message: "Invalid ID",
			Data:    nil,
		})
	}

	// Fetch all adopted pets for the adopter
	var adoptedPets []struct {
		AdoptedID uint   `json:"adopted_id"`
		PetID     uint   `json:"pet_id"`
		PetName   string `json:"pet_name"`
		PetAge    int    `json:"pet_age"`
		AgeType   string `json:"age_type"`
		PetSex    string `json:"pet_sex"`
		PetType   string `json:"pet_type"`
		Status    string `json:"status"`
		PetImage1 string `json:"pet_image1"`
	}

	// Query the database
	if err := middleware.DBConn.Table("adopterpets").
		Select("adopterpets.adopted_id, adopterpets.pet_id, petinfo.pet_name, petinfo.pet_age, petinfo.age_type, petinfo.pet_sex, petinfo.pet_type, petinfo.status, petmedia.pet_image1").
		Joins("JOIN petinfo ON adopterpets.pet_id = petinfo.pet_id").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("adopterpets.adopter_id = ?", adopterID).
		Scan(&adoptedPets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch adopted pets",
			"details": err.Error(),
		})
	}

	// Check if no pets were found
	if len(adoptedPets) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No adopted pets found for the specified adopter ID",
		})
	}

	// Return the adopted pets with their information
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adopted pets retrieved successfully",
		"data":    adoptedPets,
	})
}

func GetShelterWithPetsByID(c *fiber.Ctx) error {
	shelterID := c.Params("id")
	if shelterID == "" {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "404",
			Message: "Shelter ID is missing",
		})
	}

	var shelterInfo models.ShelterInfo
	ShelterResult := middleware.DBConn.Debug().Preload("ShelterMedia").Where("shelter_id = ?", shelterID).First(&shelterInfo)
	if errors.Is(ShelterResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "404",
			Message: "Shelter not found",
			Data:    nil,
		})
	} else if ShelterResult.Error != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    nil,
		})
	}

	var petInfo []models.PetInfo
	PetResult := middleware.DBConn.Debug().Preload("PetMedia").Where("shelter_id = ? AND status = ?", shelterID, "available").Find(&petInfo)
	if errors.Is(PetResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "404",
			Message: "Pets not found",
			Data:    nil,
		})
	} else if PetResult.Error != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter and pets retrieved successfully",
		"data": fiber.Map{
			"shelter": shelterInfo,
			"pets":    petInfo,
		},
	})
}

func CreateAdoptionSubmission(c *fiber.Ctx) error {
	petIDStr := c.Params("pet_id")
	petID, err := strconv.Atoi(petIDStr)
	if err != nil || petID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or missing petID"})
	}

	var pet models.PetInfo
	if err := middleware.DBConn.Where("pet_id = ?", petID).First(&pet).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pet not found"})
	}

	adopterID, err := middleware.GetAdopterIDFromJWT(c)
	if err != nil || adopterID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized or invalid token"})
	}

	// Form fields
	altFName := c.FormValue("altFName")
	altLName := c.FormValue("altLName")
	relationship := c.FormValue("relationship")
	altContactNumber := c.FormValue("altContactNumber")
	altEmail := c.FormValue("altEmail")
	altValidID := c.FormValue("altValidID")
	validID := c.FormValue("validID")
	reasonForAdoption := c.FormValue("reasonForAdoption")
	idealPetDescription := c.FormValue("idealPetDescription")
	housingSituation := c.FormValue("housingSituation")
	petsAtHome := c.FormValue("petsAtHome")
	allergies := c.FormValue("allergies")
	familySupport := c.FormValue("familySupport")
	pastPets := c.FormValue("pastPets")
	interviewSetting := c.FormValue("interviewSetting")

	if altFName == "" || altLName == "" || relationship == "" || altContactNumber == "" || altEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Alternate contact fields are required."})
	}

	if reasonForAdoption == "" || idealPetDescription == "" || housingSituation == "" || petsAtHome == "" || allergies == "" || familySupport == "" || pastPets == "" || interviewSetting == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Please complete all required questionnaire fields."})
	}

	if len(altContactNumber) < 11 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Contact numbers must be at least 11 digits."})
	}
	for _, ch := range altContactNumber {
		if ch < '0' || ch > '9' {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Contact number must only contain digits."})
		}
	}

	if !strings.Contains(altEmail, "@") || !strings.Contains(altEmail, ".") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid alternate email format."})
	}

	var existingSubmission models.AdoptionSubmission
	if err := middleware.DBConn.Where("alt_email = ?", altEmail).First(&existingSubmission).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Alternate email is already used in a previous submission. Please use a different one."})
	}

	tx := middleware.DBConn.Begin()
	if tx == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to begin DB transaction"})
	}

	submission := models.AdoptionSubmission{
		AdopterID:           adopterID,
		PetID:               uint(petID),
		ShelterID:           pet.ShelterID,
		AltFName:            altFName,
		AltLName:            altLName,
		Relationship:        relationship,
		AltContactNumber:    altContactNumber,
		AltEmail:            altEmail,
		AltValidID:          altValidID,
		ValidID:             validID,
		ReasonForAdoption:   reasonForAdoption,
		IdealPetDescription: idealPetDescription,
		HousingSituation:    housingSituation,
		PetsAtHome:          petsAtHome,
		Allergies:           allergies,
		FamilySupport:       familySupport,
		PastPets:            pastPets,
		InterviewSetting:    interviewSetting,
	}

	if err := tx.Create(&submission).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save submission"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	var totalHomeUploadSize int64
	if homeFiles := form.File["homeFiles"]; len(homeFiles) > 0 {
		for _, file := range homeFiles {
			totalHomeUploadSize += file.Size
			if totalHomeUploadSize > 15*1024*1024 {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Total home photo/PDF size exceeds 15MB"})
			}
			base64Data, err := encodeFileToBase64(file)
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode home file"})
			}
			photo := models.ApplicationPhoto{
				ApplicationID: submission.ApplicationID,
				PhotoType:     "Home Photo or PDF",
				Base64Data:    base64Data,
				UploadedAt:    time.Now(),
			}
			if err := tx.Create(&photo).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save home photo or PDF"})
			}
		}
	}

	if validIDHeader, err := c.FormFile("validIDPhoto"); err == nil {
		if validIDHeader.Size > 8*1024*1024 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Valid ID file size exceeds 8MB"})
		}
		base64Data, err := encodeFileToBase64(validIDHeader)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode Valid ID photo"})
		}
		photo := models.ApplicationPhoto{
			ApplicationID: submission.ApplicationID,
			PhotoType:     "Valid ID",
			Base64Data:    base64Data,
			UploadedAt:    time.Now(),
		}
		if err := tx.Create(&photo).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save Valid ID photo"})
		}
	}

	if altValidIDHeader, err := c.FormFile("altValidIDPhoto"); err == nil {
		if altValidIDHeader.Size > 8*1024*1024 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Alternate Valid ID file size exceeds 8MB"})
		}
		base64Data, err := encodeFileToBase64(altValidIDHeader)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode Alternate Valid ID photo"})
		}
		photo := models.ApplicationPhoto{
			ApplicationID: submission.ApplicationID,
			PhotoType:     "Alternate Valid ID",
			Base64Data:    base64Data,
			UploadedAt:    time.Now(),
		}
		if err := tx.Create(&photo).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save Alternate Valid ID photo"})
		}
	}

	tx.Commit()
	return c.JSON(fiber.Map{
		"message":    "Adoption submission successful",
		"submission": submission,
		"shelter_id": pet.ShelterID,
	})
}

func encodeFileToBase64(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileData), nil
}
