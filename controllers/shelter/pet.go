package controllers

import (
	"encoding/base64"
	"errors"
	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Struct to combine PetInfo and PetMedia based on pet_id
type PetResponse struct {
	PetID     uint     `json:"pet_id"`
	PetName   string   `json:"pet_name"`
	PetAge    int      `json:"pet_age"`
	AgeType   string   `json:"age_type"`
	PetSex    string   `json:"pet_sex"`
	ShelterID uint     `json:"shelter_id"`
	PetImages []string `json:"pet_image1"`
}

func GetAllPetsInfoByShelterID(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Fetch pet info for the given shelter
	var petInfo []models.PetInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).Order("created_at DESC").Find(&petInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pet info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching pet info",
		})
	}

	// Prepare a map to hold pet media by pet_id
	petMediaMap := make(map[uint][]string)

	// Fetch pet media for each pet and store them by pet_id
	var petMedia []models.PetMedia
	petmediaResult := middleware.DBConn.Where("pet_id IN ?", getPetIDs(petInfo)).Find(&petMedia)

	if petmediaResult.Error != nil && !errors.Is(petmediaResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching pet media",
		})
	}

	// Map pet media by pet_id
	for _, media := range petMedia {
		// Ensure the image is valid Base64 before adding
		if media.PetImage1 != "" {
			_, err := base64.StdEncoding.DecodeString(media.PetImage1)
			if err == nil {
				petMediaMap[media.PetID] = append(petMediaMap[media.PetID], media.PetImage1)
			}
		}
	}

	// Combine pet info and media into a single response
	var petResponses []PetResponse
	for _, pet := range petInfo {
		// Create response for each pet by combining pet info and media
		petResponse := PetResponse{
			PetID:     pet.PetID,
			PetName:   pet.PetName,
			PetAge:    pet.PetAge,
			AgeType:   pet.AgeType,
			PetSex:    pet.PetSex,
			ShelterID: pet.ShelterID,
			PetImages: petMediaMap[pet.PetID], // Get media for this pet
		}
		petResponses = append(petResponses, petResponse)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter info retrieved successfully",
		"data": fiber.Map{
			"pets": petResponses,
		},
	})
}

// Helper function to extract pet IDs from the petInfo slice
func getPetIDs(pets []models.PetInfo) []uint {
	var petIDs []uint
	for _, pet := range pets {
		petIDs = append(petIDs, pet.PetID)
	}
	return petIDs
}

// Struct to combine PetInfo and PetMedia based on pet_id
type PetInfoResponse struct {
	PetID           uint     `json:"pet_id"`
	PetType         string   `json:"pet_type"`
	PetName         string   `json:"pet_name"`
	PetAge          int      `json:"pet_age"`
	AgeType         string   `json:"age_type"`
	PetSex          string   `json:"pet_sex"`
	PetDescriptions string   `json:"pet_descriptions"`
	ShelterID       uint     `json:"shelter_id"`
	PetImages       []string `json:"pet_image1"`
}

func GetPetInfoByPetID(c *fiber.Ctx) error {
	// Get the pet_id from the URL params
	petID := c.Params("id")

	// Fetch pet info for the given pet_id
	var petInfo models.PetInfo
	infoResult := middleware.DBConn.Where("pet_id = ?", petID).First(&petInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pet info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching pet info",
		})
	}

	// Prepare a map to hold pet media by pet_id
	var petMedia []models.PetMedia
	petmediaResult := middleware.DBConn.Where("pet_id = ?", petID).Find(&petMedia)

	if petmediaResult.Error != nil && !errors.Is(petmediaResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching pet media",
		})
	}

	// Prepare pet media for this pet
	var petImages []string
	for _, media := range petMedia {
		// Ensure the image is valid Base64 before adding
		if media.PetImage1 != "" {
			_, err := base64.StdEncoding.DecodeString(media.PetImage1)
			if err == nil {
				petImages = append(petImages, media.PetImage1)
			}
		}
	}

	// Create a response for the pet by combining pet info and media
	petResponse := PetInfoResponse{
		PetID:           petInfo.PetID,
		PetType:         petInfo.PetType,
		PetName:         petInfo.PetName,
		PetAge:          petInfo.PetAge,
		AgeType:         petInfo.AgeType,
		PetSex:          petInfo.PetSex,
		PetDescriptions: petInfo.PetDescriptions,
		ShelterID:       petInfo.ShelterID,
		PetImages:       petImages, // Attach pet images
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet info retrieved successfully",
		"data": fiber.Map{
			"pet": petResponse,
		},
	})
}

// UpdatePetInfo and PetMedia
func UpdatePetInfo(c *fiber.Ctx) error {
	petID := c.Params("id")

	// Fetch the existing PetInfo record
	var petInfo models.PetInfo
	infoResult := middleware.DBConn.Where("pet_id = ?", petID).First(&petInfo)
	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pet info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching pet info",
		})
	}

	// Parse the request body for updates
	var updateData struct {
		PetInfo  models.PetInfo  `json:"petinfo"`
		PetMedia models.PetMedia `json:"petmedia"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Update PetInfo fields only if changed
	if updateData.PetInfo.PetName != petInfo.PetName {
		petInfo.PetName = updateData.PetInfo.PetName
	}
	if updateData.PetInfo.PetType != petInfo.PetType {
		petInfo.PetType = updateData.PetInfo.PetType
	}
	if updateData.PetInfo.PetSex != petInfo.PetSex {
		petInfo.PetSex = updateData.PetInfo.PetSex
	}
	if updateData.PetInfo.PetAge != petInfo.PetAge {
		petInfo.PetAge = updateData.PetInfo.PetAge
	}
	if updateData.PetInfo.AgeType != petInfo.AgeType {
		petInfo.AgeType = updateData.PetInfo.AgeType
	}
	if updateData.PetInfo.PetDescriptions != petInfo.PetDescriptions {
		petInfo.PetDescriptions = updateData.PetInfo.PetDescriptions
	}

	// Update the PetInfo in the database
	middleware.DBConn.Table("petinfo").Updates(&petInfo)

	// Update PetMedia if provided
	var petMedia models.PetMedia
	err := middleware.DBConn.Debug().Table("petmedia").Where("pet_id = ?", petID).First(&petMedia).Error

	if err == nil {
		// Encode the image to Base64
		encodedImage := updateData.PetMedia.PetImage1 // Assuming PetImage1 is the base64 string sent from the client

		// Update existing record with the Base64-encoded image
		petMedia.PetImage1 = encodedImage // Replace the old image with the new one

		// Update the PetMedia in the database
		middleware.DBConn.Table("petmedia").Where("pet_id = ?", petID).Update("pet_image1", petMedia.PetImage1)

	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Pet media record not found
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Pet media not found. Cannot update non-existent record.",
		})
	} else {
		// Database error during update
		return c.JSON(fiber.Map{
			"message": "Database error while updating pet media",
			"error":   err.Error(),
		})
	}

	// Return the successful response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet info and media updated successfully",
		"data": fiber.Map{
			"petinfo": petInfo,
		},
	})
}
