package controllers

import (
	"encoding/base64"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
)

func CreateAdoptionSubmission(c *fiber.Ctx) error {
	petIDStr := c.Params("pet_id")
	petID, err := strconv.Atoi(petIDStr)
	if err != nil || petID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or missing petID"})
	}

	adopterID, err := middleware.GetAdopterIDFromJWT(c)
	if err != nil || adopterID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized or invalid token"})
	}

	// Form fields
	altFName := c.FormValue("altFName")
	altLName := c.FormValue("altLName")
	relationship := c.FormValue("relationship")
	altContactNumber := c.FormValue("altContactNumber")
	altEmail := c.FormValue("altEmail")
	altValidID := c.FormValue("altValidID")
	validID := c.FormValue("validID")

	petType := c.FormValue("petType")
	shelterAnimal := c.FormValue("shelterAnimal")
	idealPetDescription := c.FormValue("idealPetDescription")
	housingSituation := c.FormValue("housingSituation")
	petsAtHome := c.FormValue("petsAtHome")
	allergies := c.FormValue("allergies")
	familySupport := c.FormValue("familySupport")
	pastPets := c.FormValue("pastPets")
	interviewSetting := c.FormValue("interviewSetting")

	// Basic input validation
	if altFName == "" || altLName == "" || relationship == "" || altContactNumber == "" || altEmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Alternate contact fields are required.",
		})
	}

	if petType == "" || idealPetDescription == "" || housingSituation == "" || petsAtHome == "" || allergies == "" || familySupport == "" || pastPets == "" || interviewSetting == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Please complete all required questionnaire fields.",
		})
	}

	// Contact number validation
	if len(altContactNumber) < 11 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Contact numbers must be at least 11 digits.",
		})
	}
	for _, ch := range altContactNumber {
		if ch < '0' || ch > '9' {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Contact number must only contain digits.",
			})
		}
	}

	// Email format validation (basic)
	if !strings.Contains(altEmail, "@") || !strings.Contains(altEmail, ".") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid alternate email format.",
		})
	}

	// Duplicate altEmail check (globally unique)
	var existingSubmission models.AdoptionSubmission
	if err := middleware.DBConn.Where("alt_email = ?", altEmail).First(&existingSubmission).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Alternate email is already used in a previous submission. Please use a different one.",
		})
	}

	tx := middleware.DBConn.Begin()
	if tx == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to begin DB transaction"})
	}

	submission := models.AdoptionSubmission{
		AdopterID:           adopterID,
		PetID:               uint(petID),
		AltFName:            altFName,
		AltLName:            altLName,
		Relationship:        relationship,
		AltContactNumber:    altContactNumber,
		AltEmail:            altEmail,
		AltValidID:          altValidID,
		ValidID:             validID,
		PetType:             petType,
		ShelterAnimal:       shelterAnimal,
		IdealPetDescription: idealPetDescription,
		HousingSituation:    housingSituation,
		PetsAtHome:          petsAtHome,
		Allergies:           allergies,
		FamilySupport:       familySupport,
		PastPets:            pastPets,
		InterviewSetting:    interviewSetting,
	}

	if err := tx.Create(&submission).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save submission"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	// General Home Photos or PDF Upload (max total size 15MB)
	var totalHomeUploadSize int64
	if homeFiles := form.File["homeFiles"]; len(homeFiles) > 0 {
		for _, file := range homeFiles {
			totalHomeUploadSize += file.Size
			if totalHomeUploadSize > 15*1024*1024 {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Total home photo/PDF size exceeds 15MB"})
			}

			base64Data, err := encodeFileToBase64(file)
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode home file"})
			}
			photo := models.ApplicationPhoto{
				ApplicationID: submission.ApplicationID,
				PhotoType:     "Home Photo or PDF",
				Base64Data:    base64Data,
				UploadedAt:    time.Now(),
			}
			if err := tx.Create(&photo).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save home photo or PDF"})
			}
		}
	}

	// Valid ID photo
	if validIDHeader, err := c.FormFile("validIDPhoto"); err == nil {
		if validIDHeader.Size > 8*1024*1024 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Valid ID file size exceeds 8MB"})
		}
		base64Data, err := encodeFileToBase64(validIDHeader)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode Valid ID photo"})
		}
		photo := models.ApplicationPhoto{
			ApplicationID: submission.ApplicationID,
			PhotoType:     "Valid ID",
			Base64Data:    base64Data,
			UploadedAt:    time.Now(),
		}
		if err := tx.Create(&photo).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save Valid ID photo"})
		}
	}

	// Alternate Valid ID photo
	if altValidIDHeader, err := c.FormFile("altValidIDPhoto"); err == nil {
		if altValidIDHeader.Size > 8*1024*1024 {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Alternate Valid ID file size exceeds 8MB"})
		}
		base64Data, err := encodeFileToBase64(altValidIDHeader)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encode Alternate Valid ID photo"})
		}
		photo := models.ApplicationPhoto{
			ApplicationID: submission.ApplicationID,
			PhotoType:     "Alternate Valid ID",
			Base64Data:    base64Data,
			UploadedAt:    time.Now(),
		}
		if err := tx.Create(&photo).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save Alternate Valid ID photo"})
		}
	}

	tx.Commit()
	return c.JSON(fiber.Map{
		"message":    "Adoption submission successful",
		"submission": submission,
	})
}

func encodeFileToBase64(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileData), nil
}

