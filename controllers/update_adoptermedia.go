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

	// Parse shelter ID
	parsedAdopterID, err := strconv.ParseUint(adopterID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid adopter ID",
		})
	}
	// Fetch existing shelter media or create a new one
	var adopterMedia models.AdopterMedia
	mediaResult := middleware.DBConn.Where("adopter_id = ?", parsedAdopterID).First(&adopterMedia)

	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		// Create a new shelter media record if not found
		adopterMedia = models.AdopterMedia{AdopterID: uint(parsedAdopterID)}
		middleware.DBConn.Create(&adopterMedia)
	} else if mediaResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching adopter media",
		})
	}

	// Handle profile image upload
	profileFile, err := c.FormFile("adopter_profile")
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
			"adopter_profile": adopterMedia.AdopterProfile,
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

func GetAdopterProfile(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	var adopter models.AdopterInfo
	var profile models.AdopterMedia

	if err := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopter).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Adopter not found"})
	}

	middleware.DBConn.Where("adopter_id = ?", adopterID).First(&profile)

	return c.JSON(fiber.Map{
		"adopter":         adopter,
		"adopter_profile": profile.AdopterProfile,
	})
}

func EditAdopterProfile(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	// Define a struct for the fields you want to update
	var updatedProfile models.AdopterInfo

	// Parse the body into the struct
	if err := c.BodyParser(&updatedProfile); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid data"})
	}

	// Check if adopter exists
	var adopter models.AdopterInfo
	if err := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopter).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Adopter not found"})
	}

	// Update adopter info (only the fields we care about)
	if err := middleware.DBConn.Model(&adopter).Updates(map[string]interface{}{
		"first_name":     updatedProfile.FirstName,
		"last_name":      updatedProfile.LastName,
		"contact_number": updatedProfile.ContactNumber,
		"email":          updatedProfile.Email,
		"address":        updatedProfile.Address,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update adopter info"})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"adopter": adopter,
	})
}
