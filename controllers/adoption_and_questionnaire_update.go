package controllers

import (
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"strings"

	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
)

// convertFileToBase64 encodes uploaded file to base64 string
func convertFileToBase64(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	if _, err := src.Read(buf); err != nil {
		return "", err
	}

	return base64Encode(buf), nil
}

// base64Encode is a helper for encoding byte slice to base64
func base64Encode(data []byte) string {
	return strings.TrimRight(base64.StdEncoding.EncodeToString(data), "\n")
}

// UpdateAdoptionAndQuestionnaire handles updating both models
func UpdateAdoptionAndQuestionnaire(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid form data",
			"details": err.Error(),
		})
	}

	appID := c.FormValue("application_id")
	if appID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing application_id",
		})
	}

	// Update Adoption Application
	var application models.AdoptionApplication
	if err := middleware.DBConn.Where("application_id = ?", appID).First(&application).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Adoption Application with ID %s not found", appID),
		})
	}
	newStatus := c.FormValue("status")
	if newStatus != "" {
		application.Status = newStatus
		middleware.DBConn.Save(&application)
	}

	// Update Questionnaire
	var questionnaire models.Questionnaires
	if err := middleware.DBConn.Where("application_id = ?", appID).First(&questionnaire).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Questionnaire for Application ID %s not found", appID),
		})
	}

	// Convert booleans to yes/no
	yesNo := func(val string) string {
		if val == "true" {
			return "yes"
		}
		return "no"
	}

	// Handle home photo updates
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
		questionnaire.HomePhotos = strings.Join(homePhotoBase64s, ",")
	}

	// Handle valid ID file update
	if file, err := c.FormFile("valid_id"); err == nil {
		if encoded, err := convertFileToBase64(file); err == nil {
			questionnaire.ValidID = encoded
		}
	}

	// Set updated fields (if provided)
	set := func(field *string, value string) {
		if value != "" {
			*field = value
		}
	}

	set(&questionnaire.PetType, c.FormValue("pet_type"))
	set(&questionnaire.SpecificShelterAnimal, yesNo(c.FormValue("specific_shelter_animal")))
	set(&questionnaire.IdealPetDescription, c.FormValue("ideal_pet_description"))
	set(&questionnaire.BuildingType, c.FormValue("building_type"))
	set(&questionnaire.Rent, yesNo(c.FormValue("rent")))
	set(&questionnaire.PetMovePlan, c.FormValue("pet_move_plan"))
	set(&questionnaire.HouseholdComposition, c.FormValue("household_composition"))
	set(&questionnaire.AllergiesToAnimals, yesNo(c.FormValue("allergies_to_animals")))
	set(&questionnaire.CareResponsibility, c.FormValue("care_responsibility"))
	set(&questionnaire.FinancialResponsibility, c.FormValue("financial_responsibility"))
	set(&questionnaire.VacationCarePlan, c.FormValue("vacation_care_plan"))
	set(&questionnaire.AloneTime, c.FormValue("alone_time"))
	set(&questionnaire.IntroductionPlan, c.FormValue("introduction_plan"))
	set(&questionnaire.FamilySupport, yesNo(c.FormValue("family_support")))
	set(&questionnaire.FamilySupportExplanation, c.FormValue("family_support_explanation"))
	set(&questionnaire.OtherPets, yesNo(c.FormValue("other_pets")))
	set(&questionnaire.PastPets, yesNo(c.FormValue("past_pets")))
	set(&questionnaire.IDProof, c.FormValue("id_proof"))
	set(&questionnaire.PreferredInterviewSetting, c.FormValue("preferred_interview_setting"))

	// Save changes
	if err := middleware.DBConn.Save(&questionnaire).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update questionnaire",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":               "Adoption application and questionnaire updated successfully",
		"updated_application":   application,
		"updated_questionnaire": questionnaire,
	})
}
