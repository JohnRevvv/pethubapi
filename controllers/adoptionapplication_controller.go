package controllers

import (
	"encoding/base64" // for Base64 encoding
	"io/ioutil"       // for reading file content
	"mime/multipart"  // for handling file uploads

	"pethub_api/middleware"
	"pethub_api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Helper function to convert uploaded file to base64
func ConvertFileToBase64(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	base64String := base64.StdEncoding.EncodeToString(fileBytes)
	return base64String, nil
}

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

	// Parse form fields from the request body
	adopterID := c.FormValue("adopter_id")
	altFName := c.FormValue("alt_f_name")
	altLName := c.FormValue("alt_l_name")
	relationship := c.FormValue("relationship")
	altContactNumber := c.FormValue("alt_contact_number")
	altEmail := c.FormValue("alt_email")
	preferredDate := c.FormValue("preferred_date")
	preferredTime := c.FormValue("preferred_time")

	// Convert adopterID from string to uint
	adopterIDInt, err := strconv.Atoi(adopterID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid adopter ID",
		})
	}

	// Get the uploaded house file from form-data
	houseFile, err := c.FormFile("housefile")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get housefile",
		})
	}

	// Convert the house file to base64
	houseFileBase64, err := ConvertFileToBase64(houseFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to convert housefile to base64",
		})
	}

	// Get the uploaded valid ID file from form-data
	validIDFile, err := c.FormFile("valid_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get valid_id",
		})
	}

	// Convert the valid ID file to base64
	validIDBase64, err := ConvertFileToBase64(validIDFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to convert valid_id to base64",
		})
	}

	// Create a new adoption application
	adoptionApplication := models.AdoptionApplication{
		PetID:            uint(petID),
		AdopterID:        uint(adopterIDInt),
		AltFName:         altFName,
		AltLName:         altLName,
		Relationship:     relationship,
		AltContactNumber: altContactNumber,
		AltEmail:         altEmail,
		HouseFile:        houseFileBase64, // Store Base64 encoded house file here
		ValidID:          validIDBase64,   // Store Base64 encoded valid ID here
		PreferredDate:    preferredDate,   // Store as string
		PreferredTime:    preferredTime,   // Store as string
		Status:           "Pending",       // Default status
	}

	// Insert into the database
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
			"preferred_date":     application.PreferredDate, // Now string
			"preferred_time":     application.PreferredTime, // Now string
			"status":             application.Status,
		})
	}

	// Return all adoption applications for the adopter
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"adopter_id":   adopterID,
		"applications": response,
	})
}

// GetAdoptionApplicationAndQuestionnaire retrieves all adoption applications and their questionnaires for a given adopter_id
func GetAdoptionApplicationAndQuestionnaire(c *fiber.Ctx) error {
	adopterID := c.Params("adopter_id")

	var applications []models.AdoptionApplication
	if err := middleware.DBConn.Where("adopter_id = ?", adopterID).Find(&applications).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve adoption applications",
		})
	}

	if len(applications) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No adoption applications found for this adopter",
		})
	}

	type ApplicationWithQuestionnaire struct {
		Application   models.AdoptionApplication `json:"application"`
		Questionnaire *models.Questionnaires     `json:"questionnaire,omitempty"`
	}

	var result []ApplicationWithQuestionnaire

	for _, app := range applications {
		var questionnaire models.Questionnaires
		err := middleware.DBConn.Where("application_id = ?", app.ApplicationID).First(&questionnaire).Error

		if err == nil {
			result = append(result, ApplicationWithQuestionnaire{
				Application:   app,
				Questionnaire: &questionnaire,
			})
		} else {
			result = append(result, ApplicationWithQuestionnaire{
				Application:   app,
				Questionnaire: nil,
			})
		}
	}

	return c.JSON(result)
}
