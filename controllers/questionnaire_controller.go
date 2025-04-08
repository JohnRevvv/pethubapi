package controllers

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
)

// Converts a file to base64 string
func convertFileToBase64(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	if _, err := io.ReadFull(src, buf); err != nil {
		return "", fmt.Errorf("failed to read uploaded file: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(buf)
	return encoded, nil
}

// CreateQuestionnaire handles the form submission and uploads
func CreateQuestionnaire(c *fiber.Ctx) error {
	applicationID := c.FormValue("application_id")
	petType := c.FormValue("pet_type")

	if applicationID == "" || petType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: application_id or pet_type",
		})
	}

	var application models.AdoptionApplication
	if err := middleware.DBConn.First(&application, applicationID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Application with ID %s not found", applicationID),
		})
	}

	// Prevent duplicate questionnaire
	var existing models.Questionnaires
	if err := middleware.DBConn.Where("application_id = ?", applicationID).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message":                   "Questionnaire already submitted for this application",
			"existing_questionnaire_id": existing.QuestionID,
		})
	}

	// Handle multiple home_photos as base64
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse multipart form",
			"details": err.Error(),
		})
	}

	var homePhotoBase64s []string
	if files, ok := form.File["home_photos"]; ok {
		for _, file := range files {
			encoded, err := convertFileToBase64(file)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error":   "Failed to encode home photo",
					"details": err.Error(),
				})
			}
			homePhotoBase64s = append(homePhotoBase64s, encoded)
		}
	}

	// Handle valid_id file as base64
	var validIDBase64 string
	if file, err := c.FormFile("valid_id"); err == nil {
		encoded, err := convertFileToBase64(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to encode ID proof",
				"details": err.Error(),
			})
		}
		validIDBase64 = encoded
	}

	// Convert "true"/"false" to "yes"/"no"
	yesNo := func(val string) string {
		if val == "true" {
			return "yes"
		}
		return "no"
	}

	questionnaire := models.Questionnaires{
		ApplicationID:            stringToUint(applicationID),
		PetType:                  petType,
		SpecificShelterAnimal:    yesNo(c.FormValue("specific_shelter_animal")),
		IdealPetDescription:      c.FormValue("ideal_pet_description"),
		BuildingType:             c.FormValue("building_type"),
		Rent:                     yesNo(c.FormValue("rent")),
		PetMovePlan:              c.FormValue("pet_move_plan"),
		HouseholdComposition:     c.FormValue("household_composition"),
		AllergiesToAnimals:       yesNo(c.FormValue("allergies_to_animals")),
		CareResponsibility:       c.FormValue("care_responsibility"),
		FinancialResponsibility:  c.FormValue("financial_responsibility"),
		VacationCarePlan:         c.FormValue("vacation_care_plan"),
		AloneTime:                c.FormValue("alone_time"),
		IntroductionPlan:         c.FormValue("introduction_plan"),
		FamilySupport:            yesNo(c.FormValue("family_support")),
		FamilySupportExplanation: c.FormValue("family_support_explanation"),
		OtherPets:                yesNo(c.FormValue("other_pets")),
		PastPets:                 yesNo(c.FormValue("past_pets")),
		IDProof:                  c.FormValue("id_proof"),
		ZoomInterviewDate:        c.FormValue("zoom_interview_date"),
		ZoomInterviewTime:        c.FormValue("zoom_interview_time"),
		ShelterVisit:             yesNo(c.FormValue("shelter_visit")),
		HomePhotos:               strings.Join(homePhotoBase64s, ","), // Joined base64 images
		ValidID:                  validIDBase64,                       // Base64 ID proof
	}

	// Save to database
	if err := middleware.DBConn.Create(&questionnaire).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to save questionnaire",
			"details": err.Error(),
		})
	}

	// Update application status
	application.Status = "Questionnaire Submitted"
	middleware.DBConn.Save(&application)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":          "Questionnaire submitted successfully",
		"questionnaire_id": questionnaire.QuestionID,
	})
}

func stringToUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}

func GetQuestionnaire(c *fiber.Ctx) error {
	applicationID := c.Params("application_id")

	var questionnaire models.Questionnaires
	if err := middleware.DBConn.Where("application_id = ?", applicationID).First(&questionnaire).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Questionnaire for application ID %s not found", applicationID),
		})
	}

	return c.Status(fiber.StatusOK).JSON(questionnaire)
}
