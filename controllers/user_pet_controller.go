package controllers

import (
	"encoding/base64"
	"errors"
	"pethub_api/middleware"
	"pethub_api/models"
	"pethub_api/models/response"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllPets retrieves all pet records from the petinfo table
// GetAllPets retrieves all pet records from the petinfo table
func GetAllPets(c *fiber.Ctx) error {
	// Define a struct to include pet info and the pet_image1 field
	type PetWithImage struct {
		models.PetInfo
		PetImage1 string `json:"pet_image1"`
	}

	// Create a slice to store all pets with images
	var pets []PetWithImage

	// Fetch all pets and their images from the database
	result := middleware.DBConn.Table("petinfo").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.status =?", "available").
		Find(&pets)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving pets",
		})
	}

	// Return the list of pets with images
	return c.JSON(pets)
}

func GetAvailablePets(c *fiber.Ctx) error {
	// Define a struct to hold the combined data

	// Create a slice to store the combined data
	var pets []models.PetInfo

	// Fetch all pets with status "available" and their images from the database
	result := middleware.DBConn.Table("petinfo").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.status = ?", "available").
		Find(&pets)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving available pets",
		})
	}

	// Return the list of available pets with their images
	return c.JSON(pets)
}

func GetPetsWithTrueStatus(c *fiber.Ctx) error {
	// Define a struct to hold the pet data

	// Create a slice to store the pets with true status
	var pets []models.PetInfo

	// Fetch pets with status "true" from the database
	result := middleware.DBConn.Preload("PetMedia").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.priority_status = ?", "true").
		Find(&pets)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving pets with true status",
		})
	}

	// Return the list of pets with true status
	return c.JSON(pets)
}

func GetPetByID(c *fiber.Ctx) error {
	petID := c.Params("pet_id")

	var petInfo models.PetInfo
	infoResult := middleware.DBConn.Debug().Preload("PetMedia").Where("pet_id = ?", petID).First(&petInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "404",
			Message: "Pet not found",
			Data:    nil,
		})
	} else if infoResult.Error != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    infoResult.Error,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet info retrieved successfully",
		"data": fiber.Map{
			"pet": petInfo,
		},
	})

}

func GetAllSheltersByID(c *fiber.Ctx) error {
	ShelterID := c.Params("id")

	// Fetch shelter info by ID
	var ShelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Table("shelterinfo").Where("shelter_id = ?", ShelterID).First(&ShelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "shelter info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Fetch shelter account associated with the shelter info
	var shelterData models.ShelterInfo
	accountResult := middleware.DBConn.Where("shelter_id = ?", ShelterID).First(&shelterData)

	if accountResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter info",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter info retrieved successfully",
		"data": fiber.Map{
			"info": ShelterInfo,
		},
	})
}

func GetAllShelters(c *fiber.Ctx) error {
	var AllShelters []models.ShelterAccount
	result := middleware.DBConn.Debug().Preload("ShelterInfo").Preload("ShelterInfo.ShelterMedia").Where("reg_status = ? AND status = ?", "approved", "active").Find(&AllShelters)

	if result.Error != nil {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "500",
			Message: "Database error",
			Data:    nil,
		})
	}

	if len(AllShelters) == 0 {
		return c.JSON(response.AdopterResponseModel{
			RetCode: "404",
			Message: "No shelters found",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter search results",
		"data": fiber.Map{
			"shelters": AllShelters,
		},
	})
}

func GetShelters(c *fiber.Ctx) error {
	// Define a struct to include shelter info and the shelter_profile field
	type ShelterWithMedia struct {
		models.ShelterInfo
		ShelterProfile *string `json:"shelter_profile"`
	}

	// Create a slice to store all shelters with their profile images
	var shelters []ShelterWithMedia

	// Fetch all shelters and their profile images from the database
	result := middleware.DBConn.Table("shelterinfo").
		Select("shelterinfo.*, sheltermedia.shelter_profile").
		Joins("LEFT JOIN sheltermedia ON shelterinfo.shelter_id = sheltermedia.shelter_id").
		Find(&shelters)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving shelters",
			"error":   result.Error.Error(),
		})
	}

	// Return the list of shelters with their profile images
	return c.JSON(fiber.Map{
		"message": "Shelters retrieved successfully",
		"data":    shelters,
	})
}

func UpdatePetStatusToPending(c *fiber.Ctx) error {
	petID := c.Params("id")

	// Update the status of the pet to "pending"
	result := middleware.DBConn.Table("petinfo").
		Where("pet_id = ?", petID).
		Update("status", "pending")

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update pet status",
			"error":   result.Error.Error(),
		})
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pet not found",
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet status updated to pending successfully",
	})
}

func FetchAndSearchAllPetsforAdopter(c *fiber.Ctx) error {
	// Get filters from query parameters
	shelterID := c.Params("id")
	petName := c.Query("pet_name")
	petSex := c.Query("sex")
	petType := c.Query("type")
	prioritystatus := c.Query("priority_status")

	var pets []models.PetInfo

	query := middleware.DBConn.Where("status = ? AND shelter_id = ?", "available", shelterID)

	// if shelterID != "" {
	// 	query = query.Where("shelter_id = ?", shelterID)
	// }
	if petName != "" {
		query = query.Where("pet_name ILIKE ?", "%"+petName+"%")
	}
	if petSex != "" {
		query = query.Where("pet_sex = ?", petSex)
	}
	if petType != "" {
		query = query.Where("pet_type = ?", petType)
	}
	if prioritystatus != "" {
		query = query.Where("priority_status = ?", prioritystatus)
	}

	result := query.Order("priority_status DESC").Order("created_at DESC").Find(&pets)
	if result.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while searching pets",
			Data:    result.Error,
		})
	}

	if len(pets) == 0 {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "No pets found",
			Data:    nil,
		})
	}

	// Pet media logic
	petMediaMap := make(map[uint][]string)
	var petMedia []models.PetMedia
	petmediaResult := middleware.DBConn.Where("pet_id IN ?", getPetIDs(pets)).Find(&petMedia)
	if petmediaResult.Error != nil && !errors.Is(petmediaResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet media",
			Data:    petmediaResult.Error,
		})
	}

	for _, media := range petMedia {
		if media.PetImage1 != "" {
			if _, err := base64.StdEncoding.DecodeString(media.PetImage1); err == nil {
				petMediaMap[media.PetID] = append(petMediaMap[media.PetID], media.PetImage1)
			}
		}
	}

	var petResponses []PetResponse
	for _, pet := range pets {
		petResponses = append(petResponses, PetResponse{
			PetID:          pet.PetID,
			PetType:        pet.PetType,
			PetName:        pet.PetName,
			PetSex:         pet.PetSex,
			PriorityStatus: pet.PriorityStatus,
			ShelterID:      pet.ShelterID,
			PetImages:      petMediaMap[pet.PetID],
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet search results",
		"data": fiber.Map{
			"pets": petResponses,
		},
	})
}

func FetchAllPets(c *fiber.Ctx) error {
	// Get filters from query parameters
	petName := c.Query("pet_name")
	petSex := c.Query("sex")
	petType := c.Query("type")
	prioritystatus := c.Query("priority_status")

	var pets []models.PetInfo

	// Start building the query without filtering by shelter_id
	query := middleware.DBConn.Where("status = ?", "available")

	if petName != "" {
		query = query.Where("pet_name ILIKE ?", "%"+petName+"%")
	}
	if petSex != "" {
		query = query.Where("pet_sex = ?", petSex)
	}
	if petType != "" {
		query = query.Where("pet_type = ?", petType)
	}
	if prioritystatus != "" {
		query = query.Where("priority_status = ?", prioritystatus)
	}

	// Execute the query
	result := query.Order("priority_status DESC").Order("created_at DESC").Find(&pets)
	if result.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pets",
			Data:    result.Error,
		})
	}

	if len(pets) == 0 {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "No pets found",
			Data:    nil,
		})
	}

	// Pet media logic
	petMediaMap := make(map[uint][]string)
	var petMedia []models.PetMedia
	petmediaResult := middleware.DBConn.Where("pet_id IN ?", getPetIDs(pets)).Find(&petMedia)
	if petmediaResult.Error != nil && !errors.Is(petmediaResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet media",
			Data:    petmediaResult.Error,
		})
	}

	for _, media := range petMedia {
		if media.PetImage1 != "" {
			if _, err := base64.StdEncoding.DecodeString(media.PetImage1); err == nil {
				petMediaMap[media.PetID] = append(petMediaMap[media.PetID], media.PetImage1)
			}
		}
	}

	var petResponses []PetResponse
	for _, pet := range pets {
		petResponses = append(petResponses, PetResponse{
			PetID:          pet.PetID,
			PetType:        pet.PetType,
			PetName:        pet.PetName,
			PetSex:         pet.PetSex,
			PriorityStatus: pet.PriorityStatus,
			ShelterID:      pet.ShelterID,
			PetImages:      petMediaMap[pet.PetID],
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All pets retrieved successfully",
		"data": fiber.Map{
			"pets": petResponses,
		},
	})
}

func GetApplicationByAdopterID(c *fiber.Ctx) error {
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

func GetAdoptionApplicationsByAdopterID(c *fiber.Ctx) error {
	adopterID := c.Params("adopter_id")

	if adopterID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "adopter_id is required",
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
		Where("adopter_id = ?", adopterID). // Filter by adopter_id
		Preload("Adopter").
		Preload("Adopter.AdopterMedia"). // Preload adopter media
		Preload("Pet").
		Preload("Pet.PetMedia"). // Preload pet media
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
			AdopterProfile: app.Adopter.AdopterMedia.AdopterProfile, // Assuming `AdopterProfile` is the field for the profile image
			PetName:        app.Pet.PetName,                         // Assuming `PetName` is the pet's name field
			Status:         app.Status,                              // Assuming `Status` is the field for application status
		})
	}

	return c.JSON(fiber.Map{
		"message": "Adoption applications retrieved successfully",
		"data":    responses,
	})
}

func GetAdoptionApplicationsByPetID2(c *fiber.Ctx) error {
	petID := c.Params("pet_id")

	if petID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "pet_id is required",
		})
	}

	// Create a custom struct just for the response
	type PetDetailsResponse struct {
		PetImage       string `json:"pet_image1"`
		PetName        string `json:"pet_name"`
		PetSex         string `json:"pet_sex"`
		PetAge         int    `json:"pet_age"`
		PetSize        string `json:"pet_size"`
		PetDescription string `json:"pet_descriptions"`
	}

	var pet models.PetInfo
	if err := middleware.DBConn.Debug().
		Preload("PetMedia").
		Where("pet_id = ?", petID).
		First(&pet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(response.ResponseModel{
				RetCode: "404",
				Message: "Pet not found",
				Data:    nil,
			})
		}
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch pet details",
			Data:    err.Error(),
		})
	}

	// Format the response to include only the required fields
	response := PetDetailsResponse{
		PetImage:       pet.PetMedia.PetImage1, // Assuming PetImage1 is the field for the pet's image
		PetName:        pet.PetName,
		PetSex:         pet.PetSex,
		PetAge:         pet.PetAge,
		PetSize:        pet.PetSize,
		PetDescription: pet.PetDescriptions,
	}

	return c.JSON(fiber.Map{
		"message": "Pet details retrieved successfully",
		"data":    response,
	})
}

func GetAdoptionSubmissionStatusByApplicationID(c *fiber.Ctx) error {
	applicationID := c.Params("application_id")

	if applicationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "application_id is required",
		})
	}

	// Define a struct to hold the status
	type SubmissionStatusResponse struct {
		Status string `json:"status"`
	}

	var submission models.AdoptionSubmission
	if err := middleware.DBConn.Debug().
		Where("application_id = ?", applicationID).
		First(&submission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(response.ResponseModel{
				RetCode: "404",
				Message: "Submission not found",
				Data:    nil,
			})
		}
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch submission status",
			Data:    err.Error(),
		})
	}

	// Return the status
	return c.JSON(fiber.Map{
		"message": "Submission status retrieved successfully",
		"data": SubmissionStatusResponse{
			Status: submission.Status,
		},
	})
}

func SubmitReport(c *fiber.Ctx) error {
	// Parse the shelter_id and adopter_id from the route parameters
	shelterIDStr := c.Params("shelter_id")
	adopterIDStr := c.Params("adopter_id")

	shelterID, err := strconv.Atoi(shelterIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Shelter ID",
		})
	}

	adopterID, err := strconv.Atoi(adopterIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Adopter ID",
		})
	}

	// Parse the request body to get the reason and description
	requestBody := struct {
		Reason      string `json:"reason"`
		Description string `json:"description"`
	}{}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate that reason and description are not empty
	if requestBody.Reason == "" || requestBody.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Reason and description are required",
		})
	}

	// Insert the report into the database
	report := models.Report{
		ShelterID:   shelterID,
		AdopterID:   adopterID,
		Reason:      requestBody.Reason,
		Description: requestBody.Description,
		Status:      "reported", // Automatically set the status to "reported"
		CreatedAt:   time.Now(),
	}

	if err := middleware.DBConn.Debug().Create(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to submit report",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Report submitted successfully",
		"data":    report,
	})
}

func ShowPetsByAdopterID(c *fiber.Ctx) error {
	adopterID := c.Params("adopter_id")

	if adopterID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "adopter_id is required",
		})
	}

	type AdoptedListResponse struct {
		ApplicationID uint   `json:"application_id"`
		PetImage      string `json:"pet_image1"`
		PetName       string `json:"pet_name"`
		ShelterName   string `json:"shelter_name"`
		Status        string `json:"status"`
	}

	var submissions []models.AdoptionSubmission
	if err := middleware.DBConn.Where("adopter_id = ?", adopterID).Find(&submissions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"retCode": "500",
			"message": "Failed to fetch adoption submissions",
			"data":    err.Error(),
		})
	}

	if len(submissions) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"retCode": "404",
			"message": "No pets found for the adopter",
			"data":    nil,
		})
	}

	var results []AdoptedListResponse

	for _, submission := range submissions {
		var pet models.PetInfo
		if err := middleware.DBConn.Where("pet_id = ?", submission.PetID).First(&pet).Error; err != nil {
			continue
		}

		var petMedia models.PetMedia
		if err := middleware.DBConn.Where("pet_id = ?", submission.PetID).First(&petMedia).Error; err != nil {
			petMedia.PetImage1 = ""
		}

		var shelter models.ShelterInfo
		if err := middleware.DBConn.Where("shelter_id = ?", pet.ShelterID).First(&shelter).Error; err != nil {
			shelter.ShelterName = ""
		}

		results = append(results, AdoptedListResponse{
			ApplicationID: submission.ApplicationID,
			PetImage:      petMedia.PetImage1,
			PetName:       pet.PetName,
			ShelterName:   shelter.ShelterName,
			Status:        submission.Status,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Adopted pets retrieved successfully",
		"data":    results,
	})
}

func CheckApplicationExists(c *fiber.Ctx) error {
	petIdStr := c.Query("petId")
	adopterIdStr := c.Query("adopterId")

	if petIdStr == "" || adopterIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "petId and adopterId query parameters are required",
		})
	}

	petId, err := strconv.ParseUint(petIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid petId parameter",
		})
	}

	adopterId, err := strconv.ParseUint(adopterIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid adopterId parameter",
		})
	}

	var submissions []models.AdoptionSubmission
	result := middleware.DBConn.Find(&submissions)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}

	petExists := false
	adopterExists := false
	petAndAdopterMatch := false
	var matchedApplicationID uint = 0

	rejectedStatuses := map[string]bool{
		"application_reject": true,
		"interview_reject":   true,
		"approved_reject":    true,
	}

	for _, submission := range submissions {
		if submission.PetID == uint(petId) {
			petExists = true
		}
		if submission.AdopterID == uint(adopterId) {
			if !rejectedStatuses[submission.Status] {
				adopterExists = true
				if submission.PetID == uint(petId) {
					petAndAdopterMatch = true
					matchedApplicationID = submission.ApplicationID
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"pet_exists":      petExists,
		"adopter_exists":  adopterExists,
		"pet_and_adopter": petAndAdopterMatch,
		"application_id":  matchedApplicationID,
	})

}
