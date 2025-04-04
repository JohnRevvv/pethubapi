package controllers

import (
	"pethub_api/middleware"
	"pethub_api/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SubmitAdoptionApplication handles the submission of the adoption application
func SubmitAdoptionApplication(c *fiber.Ctx) error {
	// Get pet_id from URL
	petIDParam := c.Params("pet_id")
	petID, err := strconv.Atoi(petIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid pet ID",
		})
	}

	// Parse request body
	requestBody := struct {
		AdopterID        uint   `json:"adopter_id"`
		AltFName         string `json:"alt_f_name"`
		AltLName         string `json:"alt_l_name"`
		Relationship     string `json:"relationship"`
		AltContactNumber string `json:"alt_contact_number"`
		AltEmail         string `json:"alt_email"`
		HouseFile        string `json:"housefile"`
		ValidID          string `json:"valid_id"`
		PreferredDate    string `json:"preferred_date"`
		PreferredTime    string `json:"preferred_time"`
	}{}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Convert PreferredDate and PreferredTime from string to time.Time
	preferredDate, err := time.Parse("2006-01-02", requestBody.PreferredDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid date format. Use YYYY-MM-DD.",
		})
	}

	preferredTime, err := time.Parse("15:04", requestBody.PreferredTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid time format. Use HH:MM (24-hour format).",
		})
	}

	// Create new adoption application
	adoptionApplication := models.AdoptionApplication{
		PetID:            uint(petID),
		AdopterID:        requestBody.AdopterID,
		AltFName:         requestBody.AltFName,
		AltLName:         requestBody.AltLName,
		Relationship:     requestBody.Relationship,
		AltContactNumber: requestBody.AltContactNumber,
		AltEmail:         requestBody.AltEmail,
		HouseFile:        requestBody.HouseFile,
		ValidID:          requestBody.ValidID,
		PreferredDate:    preferredDate,
		PreferredTime:    preferredTime,
		Status:           "Pending", // Default status
	}

	// Insert into database
	if err := middleware.DBConn.Create(&adoptionApplication).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to submit adoption application",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Adoption application submitted successfully",
		"data": fiber.Map{
			"application": adoptionApplication,
		},
	})
}

// GetAdoptionApplication retrieves all adoption applications for a specific adopter
func GetAdoptionApplication(c *fiber.Ctx) error {
	// Get adopter_id from the URL parameters
	adopterIDParam := c.Params("adopter_id")
	adopterID, err := strconv.Atoi(adopterIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid adopter ID",
		})
	}

	// Fetch all adoption applications for the given adopter_id
	var applications []models.AdoptionApplication
	if err := middleware.DBConn.Where("adopter_id = ?", adopterID).Find(&applications).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No applications found for the given adopter",
		})
	}

	// Create a slice to hold the responses for each application
	var response []fiber.Map

	// Loop through each application and construct the response
	for _, application := range applications {
		response = append(response, fiber.Map{
			"adopter_id":         application.AdopterID,
			"application_id":     application.ApplicationID,
			"pet_id":             application.PetID,
			"alt_f_name":         application.AltFName,
			"alt_l_name":         application.AltLName,
			"relationship":       application.Relationship,
			"alt_contact_number": application.AltContactNumber,
			"alt_email":          application.AltEmail,
			"housefile":          application.HouseFile,
			"valid_id":           application.ValidID,
			"preferred_date":     application.PreferredDate,
			"preferred_time":     application.PreferredTime,
			"status":             application.Status,
		})
	}

	// Return all adoption applications for the adopter
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"adopter_id":   adopterID,
		"applications": response,
	})
}
