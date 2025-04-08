package controllers

// import (
// 	"encoding/base64"
// 	"errors"
// 	"io/ioutil"
// 	"pethub_api/middleware"
// 	"pethub_api/models"
// 	"strconv"

// 	"github.com/gofiber/fiber/v2"
// 	"gorm.io/gorm"
// )

// func UploadShelterMedia(c *fiber.Ctx) error {
// 	shelterID := c.Params("id")

// 	// Parse shelter ID
// 	parsedShelterID, err := strconv.ParseUint(shelterID, 10, 32)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Invalid shelter ID",
// 		})
// 	}
// 	// Fetch existing shelter media or create a new one
// 	var shelterMedia models.ShelterMedia
// 	mediaResult := middleware.DBConn.Where("shelter_id = ?", parsedShelterID).First(&shelterMedia)

// 	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
// 		// Create a new shelter media record if not found
// 		shelterMedia = models.ShelterMedia{ShelterID: uint(parsedShelterID)}
// 		middleware.DBConn.Create(&shelterMedia)
// 	} else if mediaResult.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Database error while fetching shelter media",
// 		})
// 	}

// 	// Handle profile image upload
// 	profileFile, err := c.FormFile("shelter_profile")
// 	if err == nil {
// 		fileContent, err := profileFile.Open()
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Failed to open profile image",
// 			})
// 		}
// 		defer fileContent.Close()

// 		fileBytes, err := ioutil.ReadAll(fileContent)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Failed to read profile image",
// 			})
// 		}
// 		shelterMedia.ShelterProfile = base64.StdEncoding.EncodeToString(fileBytes) // Replace old image
// 	}

// 	// Handle cover image upload
// 	coverFile, err := c.FormFile("shelter_cover")
// 	if err == nil {
// 		fileContent, err := coverFile.Open()
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Failed to open cover image",
// 			})
// 		}
// 		defer fileContent.Close()

// 		fileBytes, err := ioutil.ReadAll(fileContent)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Failed to read cover image",
// 			})
// 		}
// 		shelterMedia.ShelterCover = base64.StdEncoding.EncodeToString(fileBytes) // Replace old image
// 	}

// 	// Explicitly update fields with WHERE condition
// 	updateResult := middleware.DBConn.Model(&models.ShelterMedia{}).
// 		Where("shelter_id = ?", parsedShelterID).
// 		Updates(map[string]interface{}{
// 			"shelter_profile": shelterMedia.ShelterProfile,
// 			"shelter_cover":   shelterMedia.ShelterCover,
// 		})

// 	if updateResult.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to update shelter media",
// 			"error":   updateResult.Error.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Shelter media uploaded/updated successfully",
// 		"data":    shelterMedia,
// 	})
// }
