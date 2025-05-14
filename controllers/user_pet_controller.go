package controllers

import (
	"encoding/base64"
	"errors"
	"pethub_api/middleware"
	"pethub_api/models"
	"pethub_api/models/response"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

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
	result := middleware.DBConn.Debug().Preload("ShelterInfo.ShelterMedia").Where("reg_status = ? AND status = ?", "approved", "active").Find(&AllShelters)

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

//working
func FetchAllPets(c *fiber.Ctx) error {
	// Get and sanitize filters from query parameters
	petName := strings.TrimSpace(c.Query("pet_name"))
	petSex := strings.TrimSpace(c.Query("sex"))
	petType := strings.TrimSpace(c.Query("type"))
	priority := strings.ToLower(strings.TrimSpace(c.Query("priority_status"))) // expecting "true" or "false"

	var pets []models.PetInfo

	// Start with base query: status = 'available'
	query := middleware.DBConn.Where("status = ?", "available")

	// Apply filters
	if petName != "" {
		query = query.Where("pet_name ILIKE ?", "%"+petName+"%")
	}
	if petSex != "" {
		query = query.Where("pet_sex = ?", petSex)
	}
	if petType != "" {
		query = query.Where("pet_type = ?", petType)
	}
	if priority == "true" {
		query = query.Where("priority_status = ?", true)
	} else if priority == "false" {
		query = query.Where("priority_status = ?", false)
	}

	// Execute the query
	result := query.Order("created_at DESC").Find(&pets)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pets",
			Data:    result.Error,
		})
	}

	if len(pets) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "No pets found",
			Data:    nil,
		})
	}

	// Fetch pet images
	petMediaMap := make(map[uint][]string)
	var petMedia []models.PetMedia
	petmediaResult := middleware.DBConn.Where("pet_id IN ?", getPetIDs(pets)).Find(&petMedia)
	if petmediaResult.Error != nil && !errors.Is(petmediaResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ShelterResponseModel{
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

	// Prepare response
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
