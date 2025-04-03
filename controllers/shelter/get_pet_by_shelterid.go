package controllers

import (
	"pethub_api/models"
	"pethub_api/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllPetsByShelterID(c *fiber.Ctx) error {
	// Retrieve the shelter ID from the URL parameter
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Shelter ID parameter is missing",
		})
	}

	shelterID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shelter ID",
		})
	}

	// Fetch all pets associated with the shelter ID
	var pets []models.PetInfo
	result := middleware.DBConn.Where("shelter_id = ?", shelterID).Find(&pets)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch pets",
		})
	}

	if len(pets) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No pets found for the specified shelter ID",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pets retrieved successfully",
		"data":    pets,
	})
}
