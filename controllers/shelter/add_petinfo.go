package controllers

import (
	"encoding/base64"
	"io"
	"pethub_api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// AddPetInfo handles adding pet information and associated media
func AddPetInfo(c *fiber.Ctx) error {
	// Get ShelterID from route parameters
	shelterIDParam := c.Params("id")
	shelterID, err := strconv.ParseUint(shelterIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Shelter ID",
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
			"message": "Failed to add pet",
		})
	}

	// Process Images
	petMedia := models.PetMedia{PetID: petInfo.PetID}
	petMedia.PetImage1 = processImage(c, "pet_image1", requestBody.PetImage1)

	if err := tx.Create(&petMedia).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add pet media",
		})
	}

	tx.Commit()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pet added successfully",
		"data": fiber.Map{
			"pet_info": petInfo,
			"image":    petMedia,
		},
	})
}

// processImage handles file uploads and Base64 strings
func processImage(c *fiber.Ctx, formKey, base64Str string) string {
	uploadedFile, err := c.FormFile(formKey)
	if err == nil {
		file, err := uploadedFile.Open()
		if err == nil {
			defer file.Close()
			fileBytes, _ := io.ReadAll(file)
			return base64.StdEncoding.EncodeToString(fileBytes)
		}
	}
	return base64Str
}
