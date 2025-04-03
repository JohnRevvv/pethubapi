package controllers

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"pethub_api/middleware"
	"pethub_api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UploadAdopterMedia(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	// Parse adopter ID
	parsedAdopterID, err := strconv.ParseUint(adopterID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid adopter ID",
		})
	}
	// Fetch existing adopter media or create a new one
	var adopterMedia models.AdopterMedia
	mediaResult := middleware.DBConn.Where("adopter_id = ?", parsedAdopterID).First(&adopterMedia)

	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		// Create a new adopter media record if not found
		adopterMedia = models.AdopterMedia{AdopterID: uint(parsedAdopterID)}
		middleware.DBConn.Create(&adopterMedia)
	} else if mediaResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching adopter media",
		})
	}

	// Handle profile image upload
	profileFile, err := c.FormFile("profile")
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
			"profile": adopterMedia.AdopterProfile,
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
