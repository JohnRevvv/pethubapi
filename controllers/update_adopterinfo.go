package controllers

import (
	"errors"
	"fmt"
	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UpdateAdopterDetails(c *fiber.Ctx) error {
	adopterID := c.Params("id")

	// Fetch adopter info by ID
	var adopterInfo models.AdopterInfo
	infoResult := middleware.DBConn.Where("adopter_id = ?", adopterID).First(&adopterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Adopter info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error while fetching adopter info",
		})
	}

	// Parse JSON body for adopter info updates
	var updateRequest models.AdopterInfo
	if err := c.BodyParser(&updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Convert struct to a map for updating only non-empty fields
	updateData := map[string]interface{}{}
	if updateRequest.FirstName != "" {
		updateData["first_name"] = updateRequest.FirstName
	}
	if updateRequest.LastName != "" {
		updateData["last_name"] = updateRequest.LastName
	}
	if updateRequest.Address != "" {
		updateData["address"] = updateRequest.Address
	}
	if updateRequest.ContactNumber != "" {
		updateData["contact_number"] = updateRequest.ContactNumber
	}
	if updateRequest.Email != "" {
		updateData["email"] = updateRequest.Email
	}

	// Debugging log
	fmt.Printf("Update Data: %+v\n", updateData)

	if len(updateData) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No fields to update",
		})
	}

	// Update adopter info fields
	if err := middleware.DBConn.Model(&adopterInfo).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update adopter info",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Adopter details updated successfully",
		"data":    adopterInfo,
	})
}
