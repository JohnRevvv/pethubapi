package controllers

import (
	"errors"
	"fmt"
	"pethub_api/middleware"
	"pethub_api/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PetInfo struct {
	PetID          uint      `json:"pet_id"`
	ShelterID      uint      `json:"shelter_id"`
	PetName        string    `json:"pet_name"`
	PetAge         uint      `json:"pet_age"`
	PetSex         string    `json:"pet_sex"`
	PetDescription string    `json:"pet_descriptions"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	PetType        *string   `json:"pet_type"` // Nullable column
	PetSize        string    `json:"pet_size"`
	PriorityStatus bool      `json:"priority_status"`
}

// GetAllPets retrieves all pet records from the petinfo table
func GetAllPets(c *fiber.Ctx) error {
	// Define a struct to include pet info and the pet_image1 field
	type PetWithImage struct {
		PetInfo
		PetImage1 string `json:"pet_image1"`
	}

	// Create a slice to store all pets with images
	var pets []PetWithImage

	// Fetch all pets and their images from the database
	result := middleware.DBConn.Table("petinfo").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Find(&pets)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving pets",
		})
	}

	// Return the list of pets with images
	return c.JSON(pets)
}

func GetPetByID(c *fiber.Ctx) error {
	// Retrieve and sanitize the petID parameter
	petIDParam := c.Params("id")
	petID, err := strconv.Atoi(petIDParam) // Convert to integer
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid pet ID format",
			"error":   err.Error(),
		})
	}

	// Define a struct to include pet info and the pet_image1 field
	var petInfo struct {
		PetID          uint      `json:"pet_id"`
		ShelterID      uint      `json:"shelter_id"`
		PetName        string    `json:"pet_name"`
		PetAge         uint      `json:"pet_age"`
		PetSex         string    `json:"pet_sex"`
		PetDescription string    `json:"pet_descriptions"`
		Status         string    `json:"status"`
		CreatedAt      time.Time `json:"created_at"`
		PetType        *string   `json:"pet_type"`
		PetSize        string    `json:"pet_size"`
		PriorityStatus bool      `json:"priority_status"`
		PetImage1      *string   `json:"pet_image1"`
	}

	// Fetch pet info and its image from the database
	infoResult := middleware.DBConn.Table("petinfo").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.pet_id = ?", petID).
		Scan(&petInfo)

	// Handle errors
	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pet info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
			"error":   infoResult.Error.Error(),
		})
	}

	// Return the pet info with its image
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet info retrieved successfully",
		"data":    petInfo,
	})
}

func GetAllSheltersByID(c *fiber.Ctx) error {
	ShelterID := c.Params("id")

	// Fetch shelter info by ID
	var ShelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Table("shelterinfo").Where("shelter_id = ?", ShelterID).First(&ShelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "shelter info not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Fetch shelter account associated with the shelter info
	var shelterData models.ShelterInfo
	accountResult := middleware.DBConn.Where("shelter_id = ?", ShelterID).First(&shelterData)

	if accountResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter info",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter info retrieved successfully",
		"data": fiber.Map{
			"info": ShelterInfo,
		},
	})
}

func GetShelter(c *fiber.Ctx) error {
	// Create a slice to store all shelter
	var shelter []models.ShelterInfo

	// Fetch all shelter from the database
	result := middleware.DBConn.Table("shelterinfo").Find(&shelter)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving pets",
		})
	}

	// Return the list of pets
	return c.JSON(shelter)
}

type AdoptionRequest struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Address          string `json:"address"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Occupation       string `json:"occupation"`
	SocialMedia      string `json:"social_media"`
	CivilStatus      string `json:"civil_status"`
	Sex              string `json:"sex"`
	Birthdate        string `json:"birthdate"`
	HasAdoptedBefore string `json:"has_adopted_before"`
	IdealPet         string `json:"ideal_pet"`
	BuildingType     string `json:"building_type"`
	RentStatus       string `json:"rent_status"`
	MovePlan         string `json:"move_plan"`
	LivingWith       string `json:"living_with"`
	Allergy          string `json:"allergy"`
	PetCare          string `json:"pet_care"`
	PetNeeds         string `json:"pet_needs"`
	VacationPlan     string `json:"vacation_plan"`
	FamilySupport    string `json:"family_support"`
	HasOtherPets     string `json:"has_other_pets"`
	HasPastPets      string `json:"has_past_pets"`
	PetId            int    `json:"pet_id"`

	// Image uploads
	FrontOfHouse string `json:"front_of_house"`
	StreetPhoto  string `json:"street_photo"`
	LivingRoom   string `json:"living_room"`
	DiningArea   string `json:"dining_area"`
	Kitchen      string `json:"kitchen"`
	Bedroom      string `json:"bedroom"`
	HouseWindow  string `json:"window"`
	FrontYard    string `json:"front_yard"`
	Backyard     string `json:"backyard"`
}

// SubmitAdoptionRequest handles the adoption request submission
func SubmitAdoptionRequest(c *fiber.Ctx) error {
	// Log the raw request body for debugging
	body := c.Body()
	fmt.Println("Raw Request Body:", string(body))

	// Initialize the struct
	var adoptionRequest AdoptionRequest // Use the correct struct from the current package

	// Parse JSON request body
	if err := c.BodyParser(&adoptionRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if adoptionRequest.Birthdate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Birthdate is required.",
		})
	}

	// Parse and validate birthdate
	birthdate, err := time.Parse("2006-01-02", adoptionRequest.Birthdate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid birthdate format. Use YYYY-MM-DD.",
			"error":   err.Error(),
		})
	}
	adoptionRequest.Birthdate = birthdate.Format("2006-01-02") // Format the birthdate as a string

	// Save the request to the database
	result := middleware.DBConn.Table("adoption_requests").Create(&adoptionRequest)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to submit adoption request",
			"error":   result.Error.Error(),
		})
	}

	// Success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Adoption request submitted successfully",
	})
}

// GetAdopterInfoByID retrieves adopter info and media by ID
// for User
