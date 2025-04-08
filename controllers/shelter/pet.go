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

// ================================================
//          GET ALL PETS BY SHELTER ID
// ================================================

func GetAllPetsInfoByShelterID(c *fiber.Ctx) error {
	shelterID := c.Params("id")

	// Fetch pet info for the given shelter
	var petInfo []models.PetInfo
	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).Find(&petInfo)

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

// ================================================
//          GET SINGLE PET BY PET ID
// ================================================

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

// ================================================
//
//	UPDATE PET BY PET ID
//
// ================================================
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
		PetInfo  models.PetInfo  `json:"pet_info"`
		PetMedia models.PetMedia `json:"pet_media"`
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
	middleware.DBConn.Updates(&petInfo)

	// Update PetMedia if provided
	if updateData.PetMedia.PetImage1 != "" {
		var petMedia models.PetMedia
		err := middleware.DBConn.Where("pet_id = ?", petID).First(&petMedia).Error

		if err == nil {
			// Update existing pet media
			middleware.DBConn.Model(&petMedia).Updates(models.PetMedia{
				PetImage1: updateData.PetMedia.PetImage1,
			})
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// Handle missing pet media (if needed)
			// Example: You could create new PetMedia if allowed
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error while checking pet media",
			})
		}
	}

	// Save updated PetInf

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet info and media updated successfully",
		"data": fiber.Map{
			"pet_info":  petInfo,
			"pet_media": updateData.PetMedia,
		},
	})
}
