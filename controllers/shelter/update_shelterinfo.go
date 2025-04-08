package controllers

// import (
// 	"errors"
// 	"fmt"
// 	"pethub_api/middleware"
// 	"pethub_api/models"

// 	"github.com/gofiber/fiber/v2"
// 	"gorm.io/gorm"
// )

// // =======================================================
// // ================UPDATE SHELTER DETAILS=================
// // =======================================================
// func UpdateShelterDetails(c *fiber.Ctx) error {
// 	shelterID := c.Params("id")

// 	// Fetch shelter info by ID
// 	var shelterInfo models.ShelterInfo
// 	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo)

// 	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"message": "Shelter info not found",
// 		})
// 	} else if infoResult.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Database error while fetching shelter info",
// 		})
// 	}

// 	// Parse JSON body for shelter info updates
// 	var updateRequest models.ShelterInfo
// 	if err := c.BodyParser(&updateRequest); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Invalid request body",
// 		})
// 	}

// 	// Convert struct to a map for updating only non-empty fields
// 	updateData := map[string]interface{}{}
// 	if updateRequest.ShelterName != "" {
// 		updateData["shelter_name"] = updateRequest.ShelterName
// 	}
// 	if updateRequest.ShelterAddress != "" {
// 		updateData["shelter_address"] = updateRequest.ShelterAddress
// 	}
// 	if updateRequest.ShelterLandmark != "" {
// 		updateData["shelter_landmark"] = updateRequest.ShelterLandmark
// 	}
// 	if updateRequest.ShelterContact != "" {
// 		updateData["shelter_contact"] = updateRequest.ShelterContact
// 	}
// 	if updateRequest.ShelterEmail != "" {
// 		updateData["shelter_email"] = updateRequest.ShelterEmail
// 	}
// 	if updateRequest.ShelterOwner != "" {
// 		updateData["shelter_owner"] = updateRequest.ShelterOwner
// 	}
// 	if updateRequest.ShelterDescription != "" {
// 		updateData["shelter_description"] = updateRequest.ShelterDescription
// 	}
// 	if updateRequest.ShelterSocial != "" {
// 		updateData["shelter_social"] = updateRequest.ShelterSocial
// 	}

// 	// Debugging log
// 	fmt.Printf("Update Data: %+v\n", updateData)

// 	if len(updateData) == 0 {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "No fields to update",
// 		})
// 	}

// 	// Update shelter info fields
// 	if err := middleware.DBConn.Model(&shelterInfo).Updates(updateData).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to update shelter info",
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Shelter details updated successfully",
// 		"data":    shelterInfo,
// 	})
// }
