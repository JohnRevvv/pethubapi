package controllers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"time"

	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// func GetAllAdopters(c *fiber.Ctx) error {
// 	// Fetch all adopter accounts
// 	var adopterAccounts []models.AdopterAccount
// 	accountResult := middleware.DBConn.Find(&adopterAccounts)

// 	if accountResult.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to fetch adopter accounts",
// 		})
// 	}

// 	// Fetch all adopter info
// 	var adopterInfos []models.AdopterInfo
// 	infoResult := middleware.DBConn.Find(&adopterInfos)

// 	if infoResult.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to fetch adopter info",
// 		})
// 	}

// 	// Combine accounts and info into a single response
// 	adopters := []fiber.Map{}
// 	for _, account := range adopterAccounts {
// 		for _, info := range adopterInfos {
// 			if account.AdopterID == info.AdopterID {
// 				adopters = append(adopters, fiber.Map{
// 					"adopter": account,
// 					"info":    info,
// 				})
// 				break
// 			}
// 		}
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Adopters retrieved successfully",
// 		"data":    adopters,
// 	})
// }

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

	// Login successful, return adopter account and info
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data": fiber.Map{
			"adopter": adopterAccount,
			"info":    adopterInfo,
		},
	})
}

func GetAdopterInfoByID(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	// Fetch shelter info by ID
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

	// Fetch adopter media by ID
	var adopterMedia models.AdopterMedia
	mediaResult := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopterMedia)

	// If no adopter media is found, set it to null
	var mediaResponse interface{}
	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		mediaResponse = nil
	} else if mediaResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching shelter media",
		})
	} else {
		// Decode Base64-encoded images
		decodedProfile, err := base64.StdEncoding.DecodeString(adopterMedia.AdopterProfile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to decode profile image",
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Adopter ID parameter is missing",
		})
	}

	adopterID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid adopter ID",
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

func GetAllPetsInfoByShelterID(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Fetch pet information by shelter ID with pet_image1 from petmedia
	var petInfo []struct {
		PetID          uint      `json:"pet_id"`
		ShelterID      uint      `json:"shelter_id"`
		PetName        string    `json:"pet_name"`
		PetAge         uint      `json:"pet_age"`
		PetSex         string    `json:"pet_sex"`
		PetDescription string    `json:"pet_descriptions"`
		Status         string    `json:"status"`
		CreatedAt      time.Time `json:"created_at"`
		PetType        *string   `json:"pet_type"`
		PetImage1      *string   `json:"pet_image1"`
	}

	// Query the database with a join
	infoResult := middleware.DBConn.Table("petinfo").
		Select("petinfo.pet_id, petinfo.shelter_id, petinfo.pet_name, petinfo.pet_age, petinfo.pet_sex, petinfo.pet_descriptions, petinfo.status, petinfo.created_at, petinfo.pet_type, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.shelter_id = ?", shelterID).
		Scan(&petInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) || len(petInfo) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No pets found for the specified shelter ID",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
			"error":   infoResult.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pets retrieved successfully",
		"data":    petInfo,
	})
}

func GetShelterWithPetsByID(c *fiber.Ctx) error {
	// Retrieve the shelter ID from the URL parameter
	shelterID := c.Params("id")
	if shelterID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Shelter ID parameter is missing",
		})
	}

	// Fetch shelter information by shelter ID
	var shelterInfo struct {
		models.ShelterInfo
		ShelterCover   *string `json:"shelter_cover"`
		ShelterProfile *string `json:"shelter_profile"`
	}
	if err := middleware.DBConn.Table("shelterinfo").
		Select("shelterinfo.*, sheltermedia.shelter_cover, sheltermedia.shelter_profile").
		Joins("LEFT JOIN sheltermedia ON shelterinfo.shelter_id = sheltermedia.shelter_id").
		Where("shelterinfo.shelter_id = ?", shelterID).
		Scan(&shelterInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Shelter not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching shelter info",
			"error":   err.Error(),
		})
	}

	// Fetch pets under the shelter by shelter ID
	var pets []struct {
		models.PetInfo
		PetImage1 *string `json:"pet_image1"`
	}
	if err := middleware.DBConn.Table("petinfo").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.shelter_id = ?", shelterID).
		Scan(&pets).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching pets",
			"error":   err.Error(),
		})
	}

	// Combine shelter info and pets into a single response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter and pets retrieved successfully",
		"data": fiber.Map{
			"shelter": shelterInfo,
			"pets":    pets,
		},
	})
}

func AdoptionApplication(c *fiber.Ctx) error {
	adopterID := c.Params("adopter_id")
	petID := c.Params("pet_id")

	// Check if adopterID or petID is invalid (non-numeric)
	adopterUint, err := strconv.ParseUint(adopterID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Adopter ID",
		})
	}

	petUint, err := strconv.ParseUint(petID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Pet ID",
		})
	}

	// Helper to convert file to base64
	convertFileToBase64 := func(fieldName string) (string, error) {
		fileHeader, err := c.FormFile(fieldName)
		if err != nil {
			return "", nil // File is optional
		}
		file, err := fileHeader.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, file); err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	}

	// Convert files to base64
	houseFileBase64, _ := convertFileToBase64("housefile")
	altValidIDBase64, _ := convertFileToBase64("alt_valid_id")
	homePhotosBase64, _ := convertFileToBase64("home_photos")
	validIDBase64, _ := convertFileToBase64("valid_id")

	// Parse form values
	requestBody := struct {
		AltFName                  string
		AltLName                  string
		Relationship              string
		AltContactNumber          string
		AltEmail                  string
		PetType                   string
		SpecificShelterAnimal     string
		IdealPetDescription       string
		BuildingType              string
		Rent                      string
		PetMovePlan               string
		HouseholdComposition      string
		AllergiesToAnimals        string
		CareResponsibility        string
		FinancialResponsibility   string
		VacationCarePlan          string
		AloneTime                 string
		IntroductionPlan          string
		FamilySupport             string
		FamilySupportExplanation  string
		OtherPets                 string
		PastPets                  string
		IDType                    string
		PreferredInterviewSetting string
	}{
		AltFName:                  c.FormValue("alt_f_name"),
		AltLName:                  c.FormValue("alt_l_name"),
		Relationship:              c.FormValue("relationship"),
		AltContactNumber:          c.FormValue("alt_contact_number"),
		AltEmail:                  c.FormValue("alt_email"),
		PetType:                   c.FormValue("pet_type"),
		SpecificShelterAnimal:     c.FormValue("specific_shelter_animal"),
		IdealPetDescription:       c.FormValue("ideal_pet_description"),
		BuildingType:              c.FormValue("building_type"),
		Rent:                      c.FormValue("rent"),
		PetMovePlan:               c.FormValue("pet_move_plan"),
		HouseholdComposition:      c.FormValue("household_composition"),
		AllergiesToAnimals:        c.FormValue("allergies_to_animals"),
		CareResponsibility:        c.FormValue("care_responsibility"),
		FinancialResponsibility:   c.FormValue("financial_responsibility"),
		VacationCarePlan:          c.FormValue("vacation_care_plan"),
		AloneTime:                 c.FormValue("alone_time"),
		IntroductionPlan:          c.FormValue("introduction_plan"),
		FamilySupport:             c.FormValue("family_support"),
		FamilySupportExplanation:  c.FormValue("family_support_explanation"),
		OtherPets:                 c.FormValue("other_pets"),
		PastPets:                  c.FormValue("past_pets"),
		IDType:                    c.FormValue("id_type"),
		PreferredInterviewSetting: c.FormValue("preferred_interview_setting"),
	}

	// Create adoption application
	adoptionApplication := models.AdoptionApplication{
		PetID:            uint(petUint),
		AdopterID:        uint(adopterUint),
		AltFName:         requestBody.AltFName,
		AltLName:         requestBody.AltLName,
		Relationship:     requestBody.Relationship,
		AltContactNumber: requestBody.AltContactNumber,
		AltEmail:         requestBody.AltEmail,
		HouseFile:        houseFileBase64,
		ValidID:          altValidIDBase64,
	}

	tx := middleware.DBConn.Begin()
	if err := tx.Create(&adoptionApplication).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create adoption application",
		})
	}

	// Save questionnaire
	questionnaire := models.Questionnaires{
		ApplicationID:             adoptionApplication.ApplicationID,
		PetType:                   requestBody.PetType,
		SpecificShelterAnimal:     requestBody.SpecificShelterAnimal,
		IdealPetDescription:       requestBody.IdealPetDescription,
		BuildingType:              requestBody.BuildingType,
		Rent:                      requestBody.Rent,
		PetMovePlan:               requestBody.PetMovePlan,
		HouseholdComposition:      requestBody.HouseholdComposition,
		AllergiesToAnimals:        requestBody.AllergiesToAnimals,
		CareResponsibility:        requestBody.CareResponsibility,
		FinancialResponsibility:   requestBody.FinancialResponsibility,
		VacationCarePlan:          requestBody.VacationCarePlan,
		AloneTime:                 requestBody.AloneTime,
		IntroductionPlan:          requestBody.IntroductionPlan,
		FamilySupport:             requestBody.FamilySupport,
		FamilySupportExplanation:  requestBody.FamilySupportExplanation,
		OtherPets:                 requestBody.OtherPets,
		PastPets:                  requestBody.PastPets,
		HomePhotos:                homePhotosBase64,
		ValidID:                   validIDBase64,
		IDType:                    requestBody.IDType,
		PreferredInterviewSetting: requestBody.PreferredInterviewSetting,
	}

	if err := tx.Create(&questionnaire).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save questionnaire",
		})
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adoption application submitted successfully",
	})
}
