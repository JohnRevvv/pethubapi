package controllers

import (
	"encoding/base64"
	"io"
	"pethub_api/middleware"
	"pethub_api/models"
	"strconv"
	"time"

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

	// Parse form values
	petAge, _ := strconv.Atoi(c.FormValue("pet_age")) // Convert age to int

	requestBody := struct {
		PetType         string `json:"pet_type"`
		PetName         string `json:"pet_name"`
		PetAge          int    `json:"pet_age"`
		AgeType         string `json:"age_type"`
		PetSex          string `json:"pet_sex"`
		PetDescriptions string `json:"pet_descriptions"`
		PetImage1       string `json:"pet_image1"`

	}{
		PetType:         c.FormValue("pet_type"),
		PetName:         c.FormValue("pet_name"),
		PetAge:          petAge,
		AgeType:         c.FormValue("age_type"),
		PetSex:          c.FormValue("pet_sex"),
		PetDescriptions: c.FormValue("pet_descriptions"),
		PetImage1:       c.FormValue("pet_image1"),
	}

	// Create PetInfo instance
	petInfo := models.PetInfo{
		ShelterID:       uint(shelterID),
		PetType:         requestBody.PetType,
		PetName:         requestBody.PetName,
		PetAge:          requestBody.PetAge,
		AgeType:         requestBody.AgeType,
		PetSex:          requestBody.PetSex,
		PetDescriptions: requestBody.PetDescriptions,
		CreatedAt:       time.Now(),
	}

	// Database transaction
	tx := middleware.DBConn.Begin()
	if err := tx.Create(&petInfo).Error; err != nil {
		tx.Rollback()
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
