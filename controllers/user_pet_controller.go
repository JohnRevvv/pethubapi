package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"pethub_api/middleware"
	"pethub_api/models"
	"pethub_api/models/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllPets retrieves all pet records from the petinfo table
// GetAllPets retrieves all pet records from the petinfo table
func GetAllPets(c *fiber.Ctx) error {
	// Define a struct to include pet info and the pet_image1 field
	type PetWithImage struct {
		models.PetInfo
		PetImage1 string `json:"pet_image1"`
	}

	// Create a slice to store all pets with images
	var pets []PetWithImage

	// Fetch all pets and their images from the database
	result := middleware.DBConn.Table("petinfo").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.status =?", "available").
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

func GetAvailablePets(c *fiber.Ctx) error {
	// Define a struct to hold the combined data

	// Create a slice to store the combined data
	var pets []models.PetInfo

	// Fetch all pets with status "available" and their images from the database
	result := middleware.DBConn.Table("petinfo").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.status = ?", "available").
		Find(&pets)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving available pets",
		})
	}

	// Return the list of available pets with their images
	return c.JSON(pets)
}

func GetPetsWithTrueStatus(c *fiber.Ctx) error {
	// Define a struct to hold the pet data

	// Create a slice to store the pets with true status
	var pets []models.PetInfo

	// Fetch pets with status "true" from the database
	result := middleware.DBConn.Preload("PetMedia").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.priority_status = ?", "true").
		Find(&pets)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving pets with true status",
		})
	}

	// Return the list of pets with true status
	return c.JSON(pets)
}

func GetPetByID(c *fiber.Ctx) error {
	petID := c.Params("id")

	// Define a struct to hold the combined data
	type PetWithMedia struct {
		PetID           uint      `json:"pet_id"`
		ShelterID       uint      `json:"shelter_id"`
		PetName         string    `json:"pet_name"`
		PetAge          uint      `json:"pet_age"`
		PetSex          string    `json:"pet_sex"`
		PetDescriptions string    `json:"pet_descriptions"`
		Status          string    `json:"status"`
		CreatedAt       time.Time `json:"created_at"`
		PetType         *string   `json:"pet_type"`
		PetImage1       *string   `json:"pet_image1"` // Nullable column for pet image
	}

	// Fetch pet info with its image by ID
	var pet PetWithMedia
	result := middleware.DBConn.Table("petinfo").
		Select("petinfo.*, petmedia.pet_image1").
		Joins("LEFT JOIN petmedia ON petinfo.pet_id = petmedia.pet_id").
		Where("petinfo.pet_id = ?", petID).
		First(&pet)

	// Handle errors
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pet info not found",
		})
	} else if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Return the pet info with its image
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet info retrieved successfully",
		"data":    pet,
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
	// Define a struct to include shelter info and the shelter_profile field
	type ShelterWithMedia struct {
		models.ShelterInfo
		ShelterProfile *string `json:"shelter_profile"`
	}

	// Create a slice to store all shelters with their profile images
	var shelters []ShelterWithMedia

	// Fetch all shelters and their profile images from the database
	result := middleware.DBConn.Table("shelterinfo").
		Select("shelterinfo.*, sheltermedia.shelter_profile").
		Joins("LEFT JOIN sheltermedia ON shelterinfo.shelter_id = sheltermedia.shelter_id").
		Find(&shelters)

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving shelters",
			"error":   result.Error.Error(),
		})
	}

	// Return the list of shelters with their profile images
	return c.JSON(fiber.Map{
		"message": "Shelters retrieved successfully",
		"data":    shelters,
	})
}

func GetSheltertry(c *fiber.Ctx) error {
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

func UpdatePetStatusToPending(c *fiber.Ctx) error {
	petID := c.Params("id")

	// Update the status of the pet to "pending"
	result := middleware.DBConn.Table("petinfo").
		Where("pet_id = ?", petID).
		Update("status", "pending")

	// Handle errors
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update pet status",
			"error":   result.Error.Error(),
		})
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pet not found",
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet status updated to pending successfully",
	})
}

func FetchAndSearchAllPetsforAdopter(c *fiber.Ctx) error {
	// Get filters from query parameters
	shelterID := c.Params("id")
	petName := c.Query("pet_name")
	petSex := c.Query("sex")
	petType := c.Query("type")
	prioritystatus := c.Query("priority_status")

	var pets []models.PetInfo

	query := middleware.DBConn.Where("status = ? AND shelter_id = ?", "available", shelterID)

	// if shelterID != "" {
	// 	query = query.Where("shelter_id = ?", shelterID)
	// }
	if petName != "" {
		query = query.Where("pet_name ILIKE ?", "%"+petName+"%")
	}
	if petSex != "" {
		query = query.Where("pet_sex = ?", petSex)
	}
	if petType != "" {
		query = query.Where("pet_type = ?", petType)
	}
	if prioritystatus != "" {
		query = query.Where("priority_status = ?", prioritystatus)
	}

	result := query.Order("priority_status DESC").Order("created_at DESC").Find(&pets)
	if result.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while searching pets",
			Data:    result.Error,
		})
	}

	if len(pets) == 0 {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "No pets found",
			Data:    nil,
		})
	}

	// Pet media logic
	petMediaMap := make(map[uint][]string)
	var petMedia []models.PetMedia
	petmediaResult := middleware.DBConn.Where("pet_id IN ?", getPetIDs(pets)).Find(&petMedia)
	if petmediaResult.Error != nil && !errors.Is(petmediaResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet media",
			Data:    petmediaResult.Error,
		})
	}

	for _, media := range petMedia {
		if media.PetImage1 != "" {
			if _, err := base64.StdEncoding.DecodeString(media.PetImage1); err == nil {
				petMediaMap[media.PetID] = append(petMediaMap[media.PetID], media.PetImage1)
			}
		}
	}

	var petResponses []PetResponse
	for _, pet := range pets {
		petResponses = append(petResponses, PetResponse{
			PetID:          pet.PetID,
			PetType:        pet.PetType,
			PetName:        pet.PetName,
			PetSex:         pet.PetSex,
			PriorityStatus: pet.PriorityStatus,
			ShelterID:      pet.ShelterID,
			PetImages:      petMediaMap[pet.PetID],
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet search results",
		"data": fiber.Map{
			"pets": petResponses,
		},
	})
}

func FetchAllPets(c *fiber.Ctx) error {
	// Get filters from query parameters
	petName := c.Query("pet_name")
	petSex := c.Query("sex")
	petType := c.Query("type")
	prioritystatus := c.Query("priority_status")

	var pets []models.PetInfo

	// Start building the query without filtering by shelter_id
	query := middleware.DBConn.Where("status = ?", "available")

	if petName != "" {
		query = query.Where("pet_name ILIKE ?", "%"+petName+"%")
	}
	if petSex != "" {
		query = query.Where("pet_sex = ?", petSex)
	}
	if petType != "" {
		query = query.Where("pet_type = ?", petType)
	}
	if prioritystatus != "" {
		query = query.Where("priority_status = ?", prioritystatus)
	}

	// Execute the query
	result := query.Order("priority_status DESC").Order("created_at DESC").Find(&pets)
	if result.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pets",
			Data:    result.Error,
		})
	}

	if len(pets) == 0 {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "No pets found",
			Data:    nil,
		})
	}

	// Pet media logic
	petMediaMap := make(map[uint][]string)
	var petMedia []models.PetMedia
	petmediaResult := middleware.DBConn.Where("pet_id IN ?", getPetIDs(pets)).Find(&petMedia)
	if petmediaResult.Error != nil && !errors.Is(petmediaResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet media",
			Data:    petmediaResult.Error,
		})
	}

	for _, media := range petMedia {
		if media.PetImage1 != "" {
			if _, err := base64.StdEncoding.DecodeString(media.PetImage1); err == nil {
				petMediaMap[media.PetID] = append(petMediaMap[media.PetID], media.PetImage1)
			}
		}
	}

	var petResponses []PetResponse
	for _, pet := range pets {
		petResponses = append(petResponses, PetResponse{
			PetID:          pet.PetID,
			PetType:        pet.PetType,
			PetName:        pet.PetName,
			PetSex:         pet.PetSex,
			PriorityStatus: pet.PriorityStatus,
			ShelterID:      pet.ShelterID,
			PetImages:      petMediaMap[pet.PetID],
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All pets retrieved successfully",
		"data": fiber.Map{
			"pets": petResponses,
		},
	})
}
