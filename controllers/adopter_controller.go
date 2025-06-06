package controllers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Error generating token",
			Data:    err.Error(),
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

func GetAdopterInfoOnly(c *fiber.Ctx) error {
	adopterID := c.Params("adopter_id")

	// Fetch shelter info by ID
	var adopterInfo models.AdopterInfo
	infoResult := middleware.DBConn.Debug().Preload("AdopterMedia").Where("adopter_id = ?", adopterID).First(&adopterInfo)

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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adopter info retrieved successfully",
		"data":    adopterInfo,
	})
}

func GetAdopterInfoByID(c *fiber.Ctx) error {
	adopterID := c.Params("adopter_id")

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

func GetOtherPetsByAdopterID(c *fiber.Ctx) error {
	shelterID := c.Params("shelter_id")

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
			"pets": petInfo,
		},
	})
}

func CreateAdoption(c *fiber.Ctx) error {
	shelterIDStr := c.Params("shelter_id")
	petIDStr := c.Params("pet_id")
	adopterIDStr := c.Params("adopter_id")

	adopterID, err := strconv.ParseUint(adopterIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid Adopter ID"})
	}
	shelterID, err := strconv.ParseUint(shelterIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid Shelter ID"})
	}
	petID, err := strconv.ParseUint(petIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid Pet ID"})
	}

	getBase64 := func(field string) (string, error) {
		file, err := c.FormFile(field)
		if err != nil {
			return "", err
		}
		f, err := file.Open()
		if err != nil {
			return "", err
		}
		defer f.Close()

		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(bytes), nil
	}

	// Required ID images
	validID, err := getBase64("adopter_valid_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Adopter valid ID is required"})
	}
	altID, err := getBase64("alt_valid_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Alternative valid ID is required"})
	}

	// Optional home images
	home1, _ := getBase64("home_image1") // front of house
	home2, _ := getBase64("home_image2") // living room
	home3, _ := getBase64("home_image3") // kitchen
	home4, _ := getBase64("home_image4")
	home5, _ := getBase64("home_image5")
	home6, _ := getBase64("home_image6")
	home7, _ := getBase64("home_image7")
	home8, _ := getBase64("home_image8")

	// Save photos
	photos := models.ApplicationPhotos{
		AdopterIDType:  c.FormValue("adopter_id_type"),
		AdopterValidID: validID,
		AltIDType:      c.FormValue("alt_id_type"),
		AltValidID:     altID,
		HomeImage1:     home1,
		HomeImage2:     home2,
		HomeImage3:     home3,
		HomeImage4:     home4,
		HomeImage5:     home5,
		HomeImage6:     home6,
		HomeImage7:     home7,
		HomeImage8:     home8,
	}

	if err := middleware.DBConn.Debug().Create(&photos).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to save application photos", "error": err.Error()})
	}

	// Create adoption record
	adoption := models.AdoptionSubmission{
		ShelterID:           uint(shelterID),
		PetID:               uint(petID),
		AdopterID:           uint(adopterID),
		AltFName:            c.FormValue("alt_f_name"),
		AltLName:            c.FormValue("alt_l_name"),
		Relationship:        c.FormValue("relationship"),
		AltContactNumber:    c.FormValue("alt_contact_number"),
		AltEmail:            c.FormValue("alt_email"),
		ReasonForAdoption:   c.FormValue("reason_for_adoption"),
		IdealPetDescription: c.FormValue("ideal_pet_description"),
		HousingSituation:    c.FormValue("housing_situation"),
		PetsAtHome:          c.FormValue("pets_at_home"),
		Allergies:           c.FormValue("allergies"),
		FamilySupport:       c.FormValue("family_support"),
		PastPets:            c.FormValue("past_pets"),
		InterviewSetting:    c.FormValue("interview_setting"),
		Status:              "pending",
		CreatedAt:           time.Now(),
		ImageID:             photos.ImageID,
	}

	if err := middleware.DBConn.Debug().Create(&adoption).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to save adoption submission", "error": err.Error()})
	}

	var pet models.PetInfo
	if err := middleware.DBConn.First(&pet, petID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Pet not found"})
	}
	// Check if the pet is already adopted or unavailable
	if pet.Status == "adopted" {
		return c.Status(400).JSON(fiber.Map{"message": "Pet has already been adopted"})
	}

	// Count current adoption submissions for this pet
	var adoptionCount int64
	if err := middleware.DBConn.Model(&models.AdoptionSubmission{}).
		Where("pet_id = ?", petID).
		Count(&adoptionCount).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to count adoptions"})
	}

	// If 3 or more adoptions, set pet status to "pending"
	if adoptionCount >= 3 {
		pet.Status = "unavailable"
		if err := middleware.DBConn.Save(&pet).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "Failed to update pet status"})
		}
	}

	return c.Status(201).JSON(fiber.Map{
		"message":  "Adoption submitted successfully",
		"adoption": adoption,
	})
}

func UpdateAdopterDetails(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	// Fetch adopter info by ID
	var adopterInfo models.AdopterInfo
	infoResult := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Adopter info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching adopter info",
		})
	}

	// Parse JSON body for adopter info updates
	var updateRequest models.AdopterInfo
	if err := c.BodyParser(&updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Convert struct to a map for updating only non-empty fields
	updateData := map[string]interface{}{}
	if updateRequest.FirstName != "" {
		updateData["first_name"] = updateRequest.FirstName
	}
	if updateRequest.LastName != "" {
		updateData["last_name"] = updateRequest.LastName
	}
	if updateRequest.Address != "" {
		updateData["address"] = updateRequest.Address
	}
	if updateRequest.ContactNumber != "" {
		updateData["contact_number"] = updateRequest.ContactNumber
	}
	if updateRequest.Email != "" {
		updateData["email"] = updateRequest.Email
	}

	// Debugging log
	fmt.Printf("Update Data: %+v\n", updateData)

	if len(updateData) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No fields to update",
		})
	}

	// Update adopter info fields
	if err := middleware.DBConn.Model(&adopterInfo).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update adopter info",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adopter details updated successfully",
		"data":    adopterInfo,
	})
}

func UploadAdopterMedia(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	// Parse shelter ID
	parsedAdopterID, err := strconv.ParseUint(adopterID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid adopter ID",
		})
	}
	// Fetch existing shelter media or create a new one
	var adopterMedia models.AdopterMedia
	mediaResult := middleware.DBConn.Where("adopter_id = ?", parsedAdopterID).First(&adopterMedia)

	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		// Create a new shelter media record if not found
		adopterMedia = models.AdopterMedia{AdopterID: uint(parsedAdopterID)}
		middleware.DBConn.Create(&adopterMedia)
	} else if mediaResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching adopter media",
		})
	}

	// Handle profile image upload
	profileFile, err := c.FormFile("adopter_profile")
	if err == nil {
		fileContent, err := profileFile.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to open profile image",
			})
		}
		defer fileContent.Close()

		fileBytes, err := ioutil.ReadAll(fileContent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to read profile image",
			})
		}
		adopterMedia.AdopterProfile = base64.StdEncoding.EncodeToString(fileBytes) // Replace old image
	}

	// Explicitly update fields with WHERE condition
	updateResult := middleware.DBConn.Model(&models.AdopterMedia{}).
		Where("adopter_id = ?", parsedAdopterID).
		Updates(map[string]interface{}{
			"adopter_profile": adopterMedia.AdopterProfile,
		})

	if updateResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update adopter media",
			"error":   updateResult.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adopter media uploaded/updated successfully",
		"data":    adopterMedia,
	})
}

func GetAdoptionNotifications(c *fiber.Ctx) error {
	adopterIDStr := c.Params("adopter_id")
	adopterID, err := strconv.Atoi(adopterIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid adopter ID",
		})
	}

	// Step 1: Fetch all applications for this adopter
	var applications []models.AdoptionSubmission
	if err := middleware.DBConn.Preload("Pet").
		Where("adopter_id = ?", adopterID).
		Order("created_at DESC").
		Find(&applications).Error; err != nil {
		log.Println("Error fetching applications:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch adoption applications",
		})
	}

	// Adjustment: If no applications, return empty notifications
	if len(applications) == 0 {
		return c.JSON(fiber.Map{
			"notifications": []models.Notification{},
		})
	}

	// Step 2: Generate notification if needed
	for _, app := range applications {
		var notifType, notifStatus, notifCategory, notifTitle string

		if app.Status == "application_reject" || app.Status == "interview_reject" || app.Status == "approved_reject" {
			notifType = "rejected"
			notifStatus = "rejected"
			notifCategory = "rejected"
			switch app.Status {
			case "application_reject":
				notifTitle = "Application Rejected"
			case "interview_reject":
				notifTitle = "Interview Rejected"
			case "approved_reject":
				notifTitle = "Approval Rejected"
			}
		} else {
			switch app.Status {
			case "pending":
				notifType = "application"
				notifStatus = "pending"
				notifCategory = "inprogress"
				notifTitle = "Application Pending"
			case "in queue":
				notifType = "application"
				notifStatus = "in queue"
				notifCategory = "inprogress"
				notifTitle = "Application In Queue"
			case "interview":
				notifType = "interview"
				notifStatus = "interview"
				notifCategory = "inprogress"
				notifTitle = "Interview Scheduled"
			case "approved":
				notifType = "approved"
				notifStatus = "approved"
				notifCategory = "approved"
				notifTitle = "Application Approved"
			case "completed":
				notifType = "completed"
				notifStatus = "completed"
				notifCategory = "completed"
				notifTitle = "Adoption Completed"
			default:
				notifType = "application"
				notifStatus = app.Status
				notifCategory = "inprogress"
				notifTitle = "Application Update"
			}
		}

		message := fmt.Sprintf("Your adoption application for %s is now '%s'.", app.Pet.PetName, strings.Title(app.Status))

		// Check if notification already exists (no deleted filter)
		var existingNotif models.Notification
		err := middleware.DBConn.
			Where("adopter_id = ? AND pet_id = ? AND status = ?",
				adopterID, app.PetID, notifStatus).
			First(&existingNotif).Error

		if err == nil {
			continue // Notification already exists
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Error checking existing notification:", err)
			continue
		}

		// Create new notification
		notif := models.Notification{
			AdopterID: uint(adopterID),
			PetID:     uint(app.PetID),
			Title:     notifTitle,
			Message:   message,
			Type:      notifType,
			Status:    notifStatus,
			Category:  notifCategory,
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		if err := middleware.DBConn.Create(&notif).Error; err != nil {
			log.Println("Error creating notification:", err)
		}
	}

	// Step 3: Return all notifications for this adopter (no deleted filter)
	var notifications []models.Notification
	if err := middleware.DBConn.
		Where("adopter_id = ?", adopterID).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		log.Println("Error fetching notifications:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notifications",
		})
	}

	return c.JSON(fiber.Map{
		"notifications": notifications,
	})
}

func EditAdopterProfileOLD(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	// Define a struct for the fields you want to update
	var updatedProfile models.AdopterInfo

	// Parse the body into the struct
	if err := c.BodyParser(&updatedProfile); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	// Check if adopter exists
	var adopter models.AdopterInfo
	if err := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopter).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Adopter not found"})
	}

	// Update adopter info (only the fields we care about)
	if err := middleware.DBConn.Model(&adopter).Updates(map[string]interface{}{
		"first_name":     updatedProfile.FirstName,
		"last_name":      updatedProfile.LastName,
		"contact_number": updatedProfile.ContactNumber,
		"email":          updatedProfile.Email,
		"address":        updatedProfile.Address,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update adopter info"})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"adopter": adopter,
	})
}

func EditAdopterProfile(c *fiber.Ctx) error {
	adopterIdParam := c.Params("adopter_id")
	adopterId, err := strconv.Atoi(adopterIdParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.AdopterResponseModel{
			RetCode: "400",
			Message: "Invalid adopter ID",
			Data:    nil,
		})
	}

	var adopterInfo models.AdopterInfo
	infoResult := middleware.DBConn.Where("adopter_id = ?", adopterId).First(&adopterInfo)
	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "404",
			Message: "Adopter not found",
			Data:    nil,
		})
	}
	if infoResult.Error != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    nil,
		})
	}

	// Update adopter info fields
	adopterInfo.FirstName = c.FormValue("first_name")
	adopterInfo.LastName = c.FormValue("last_name")
	adopterInfo.Address = c.FormValue("address")
	adopterInfo.ContactNumber = c.FormValue("contact_number")
	adopterInfo.Email = c.FormValue("email")
	adopterInfo.Occupation = c.FormValue("occupation")
	adopterInfo.CivilStatus = c.FormValue("civil_status")
	adopterInfo.SocialMedia = c.FormValue("social_media")

	middleware.DBConn.Debug().Where("adopter_id = ?", adopterId).Updates(&adopterInfo)

	// Load or initialize media record
	var adopterMedia models.AdopterMedia
	mediaResult := middleware.DBConn.Where("adopter_id = ?", adopterId).First(&adopterMedia)

	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		// Initialize a new record if not found
		adopterMedia = models.AdopterMedia{
			AdopterID: uint(adopterId),
		}
	}

	// Handle image upload
	imageFile, err := c.FormFile("adopter_profile")
	if err == nil {
		file, _ := imageFile.Open()
		defer file.Close()
		buf := new(bytes.Buffer)
		io.Copy(buf, file)

		newImageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

		if newImageBase64 != "" {
			adopterMedia.AdopterProfile = newImageBase64
			middleware.DBConn.Debug().Save(&adopterMedia)
		}
	}

	return c.JSON(response.AdopterResponseModel{
		RetCode: "200",
		Message: "Adopter profile updated successfully",
		Data:    adopterInfo,
	})
}

func SetNotificationReadStatus(c *fiber.Ctx) error {
	notificationID := c.Params("id")

	var body struct {
		IsRead bool `json:"is_read"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var notif models.Notification
	if err := middleware.DBConn.First(&notif, notificationID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notification not found",
		})
	}

	notif.IsRead = body.IsRead
	if err := middleware.DBConn.Save(&notif).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update notification",
		})
	}

	status := "read"
	if !body.IsRead {
		status = "unread"
	}
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Notification marked as %s", status),
	})
}

func CountUnreadNotifications(c *fiber.Ctx) error {
	adopterID := c.Params("adopter_id")

	var count int64
	if err := middleware.DBConn.Model(&models.Notification{}).
		Where("adopter_id = ? AND is_read = false", adopterID).
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count unread notifications",
		})
	}

	return c.JSON(fiber.Map{
		"unread_count": count,
	})
}

func GetNotificationByID(c *fiber.Ctx) error {
	notificationID := c.Params("id")

	var notif models.Notification
	if err := middleware.DBConn.First(&notif, notificationID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notification not found",
		})
	}

	// Mark as read if not already
	if !notif.IsRead {
		notif.IsRead = true
		if err := middleware.DBConn.Save(&notif).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to mark notification as read",
			})
		}
	}

	return c.JSON(fiber.Map{
		"notification": notif,
	})
}

func GetAdopterProfile(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	var adopter models.AdopterInfo
	var profile models.AdopterMedia

	if err := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopter).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Adopter not found"})
	}

	middleware.DBConn.Where("adopter_id = ?", adopterID).First(&profile)

	return c.JSON(fiber.Map{
		"adopter":         adopter,
		"adopter_profile": profile.AdopterProfile,
	})
}

func DeleteAllNotifications(c *fiber.Ctx) error {
	adopterIDStr := c.Params("adopter_id")
	adopterID, err := strconv.Atoi(adopterIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid adopter ID",
		})
	}

	// Hard delete all notifications for this adopter
	if err := middleware.DBConn.
		Where("adopter_id = ?", adopterID).
		Delete(&models.Notification{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete notifications",
		})
	}

	return c.JSON(fiber.Map{
		"message": "All notifications deleted successfully",
	})
}
