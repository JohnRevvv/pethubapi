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

func GetPetInfoByPetID(c *fiber.Ctx) error {
	petIDParam := c.Params("id")
	if petIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Pet ID parameter is missing",
		})
	}

	petID, err := strconv.Atoi(petIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pet ID",
		})
	}

	var pet models.PetInfo
	result := middleware.DBConn.Preload("PetMedia").First(&pet, petID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pet not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet retrieved successfully",
		"data":    pet,
	})
}
