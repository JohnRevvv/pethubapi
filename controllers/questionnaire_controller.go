package controllers

import (
	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
)

// CreateQuestionnaire handles new questionnaire submissions for Fiber
func CreateQuestionnaire(c *fiber.Ctx) error {
	var questionnaire models.Questionnaires
	var application models.AdoptionApplication

	// Get logged-in adopter's ID from middleware/session
	adopterID, exists := c.Locals("adopter_id").(uint)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Fetch the latest adoption application linked to the adopter
	if err := middleware.DBConn.Where("adopter_id = ?", adopterID).Order("created_at desc").First(&application).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No active adoption application found for this adopter"})
	}

	// Bind request data to questionnaire struct
	if err := c.BodyParser(&questionnaire); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data", "details": err.Error()})
	}

	// Assign the correct ApplicationID from the fetched adoption application
	questionnaire.ApplicationID = application.ApplicationID

	// Save the questionnaire
	if err := middleware.DBConn.Create(&questionnaire).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create questionnaire", "details": err.Error()})
	}

	// Optionally update the application status (if needed)
	application.Status = "Questionnaire Submitted"
	if err := middleware.DBConn.Save(&application).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update application status", "details": err.Error()})
	}

	// Respond with success
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":       "Questionnaire submitted successfully",
		"questionnaire": questionnaire,
	})
}
