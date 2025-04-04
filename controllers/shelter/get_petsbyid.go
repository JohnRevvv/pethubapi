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
