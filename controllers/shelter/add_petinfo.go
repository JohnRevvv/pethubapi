package controllers

import (
	"fmt"
	"encoding/base64"
	"io/ioutil"
	"pethub_api/models"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var DBConn *gorm.DB // Assume DB is initialized elsewhere

// Helper function to convert image file to Base64 string
func ConvertImageToBase64(c *fiber.Ctx, fieldName string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		if err.Error() == "file not found" {
			// No file uploaded, return empty string
			return "", nil
		}
		return "", fmt.Errorf("failed to get file for field '%s': %v", fieldName, err)
	}

	fileContent, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file for field '%s': %v", fieldName, err)
	}
	defer fileContent.Close()

	fileBytes, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return "", fmt.Errorf("failed to read file for field '%s': %v", fieldName, err)
	}

	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

// Updated function for adding pet info and media
func AddPetInfo(c *fiber.Ctx) error {
	shelterID := c.Params("id")
	parsedShelterID, err := strconv.Atoi(shelterID)
	if err != nil || parsedShelterID < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid shelter ID provided",
		})
	}

	// Safe conversion from int to uint
	pet := models.PetInfo{ShelterID: uint(parsedShelterID)}

	// Parse request body
	if err := c.BodyParser(&pet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body into pet data",
		})
	}

// Validate required fields
if pet.PetName == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "message": "Missing required field: pet_name",
    })
}

if pet.PetAge == 0 {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "message": "Missing required field: pet_age",
    })
}

if pet.AgeType == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "message": "Missing required field: age_type",
    })
}

if pet.PetSex == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "message": "Missing required field: pet_sex",
    })
}

if pet.PetDescriptions == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "message": "Missing required field: pet_descriptions",
    })
}


	// Save pet to database
	if err := DBConn.Create(&pet).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add pet information to the database",
		})
	}

	// Process images
	petMedia := models.PetMedia{PetID: pet.PetID}
	imageFields := []string{"pet_image1", "pet_image2", "pet_image3", "pet_image4"}
	imagePointers := []*string{&petMedia.PetImage1, &petMedia.PetImage2, &petMedia.PetImage3, &petMedia.PetImage4}

	for i, field := range imageFields {
		base64Str, err := ConvertImageToBase64(c, field)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error while processing image: " + err.Error(),
			})
		}
		*imagePointers[i] = base64Str
	}

	// Save media to database
	if err := DBConn.Create(&petMedia).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save pet media to the database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pet and media added successfully",
		"data": fiber.Map{
			"pet":   pet,
			"media": petMedia,
		},
	})
}