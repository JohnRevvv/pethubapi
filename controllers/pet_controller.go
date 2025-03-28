package controllers

import (
	"fmt"
	"pethub_api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var DB *gorm.DB // Assume DB is initialized elsewhere

func AddPet(c *fiber.Ctx) error {
	if DB == nil {
		fmt.Println("Database connection is not initialized") // Debugging log
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection is not initialized",
		})
	}

	pet := new(models.PetInfo)

	// Retrieve the shelter ID from the URL parameter
	idParam := c.Params("id")
	if idParam == "" {
		fmt.Println("Shelter ID parameter is missing") // Debugging log
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Shelter ID parameter is missing",
		})
	}

	shelterID, err := strconv.Atoi(idParam)
	if err != nil {
		fmt.Printf("Error parsing shelter ID: %v\n", err) // Debugging log
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shelter ID",
		})
	}

	// Parse the request body into the pet struct
	if err := c.BodyParser(pet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Ensure the shelter ID in the request matches the URL parameter
	if pet.ShelterID != 0 && pet.ShelterID != shelterID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Shelter ID in request body does not match URL parameter",
		})
	}

	// Set the shelter ID from the URL parameter
	pet.ShelterID = shelterID

	// Validate required fields
	if pet.PetName == "" || pet.PetAge == 0 || pet.PetSex == "" || pet.PetDescriptions == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	// Save the pet information to the database
	if err := DB.Create(&pet).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add pet information",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(pet)
}

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
	result := DB.Where("shelter_id = ?", shelterID).Find(&pets)

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
