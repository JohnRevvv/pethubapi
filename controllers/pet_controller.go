package controllers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"pethub_api/middleware"
	"pethub_api/models"
	"pethub_api/models/response"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Struct to combine PetInfo and PetMedia based on pet_id
type PetResponse struct {
	PetID           uint     `json:"pet_id"`
	PetType         string   `json:"pet_type"`
	PetName         string   `json:"pet_name"`
	PetAge          int      `json:"pet_age"`
	AgeType         string   `json:"age_type"`
	PetSex          string   `json:"pet_sex"`
	PetSize         string   `json:"pet_size"`
	PriorityStatus  bool     `json:"priority_status"`
	PetDescriptions string   `json:"pet_descriptions"`
	ShelterID       uint     `json:"shelter_id"`
	PetImages       []string `json:"pet_image1"`
}

func AddPetInfo(c *fiber.Ctx) error {
	// Get ShelterID from route parameters
	shelterIDParam := c.Params("id")
	shelterID, err := strconv.ParseUint(shelterIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Shelter ID",
		})
	}

	// Parse form values
	petAge, _ := strconv.Atoi(c.FormValue("pet_age")) // Convert age to int

	requestBody := struct {
		PetType         string `json:"pet_type"`
		PetName         string `json:"pet_name"`
		PetAge          int    `json:"pet_age"`
		AgeType         string `json:"age_type"`
		PetSex          string `json:"pet_sex"`
		PetSize         string `json:"pet_size"`
		PetDescriptions string `json:"pet_descriptions"`
		PriorityStatus  bool   `json:"priority_status"`
		PetImage1       string `json:"pet_image1"`
	}{
		PetType:         c.FormValue("pet_type"),
		PetName:         c.FormValue("pet_name"),
		PetAge:          petAge,
		AgeType:         c.FormValue("age_type"),
		PetSex:          c.FormValue("pet_sex"),
		PetSize:         c.FormValue("pet_size"),
		PetDescriptions: c.FormValue("pet_descriptions"),
		PriorityStatus:  c.FormValue("priority_status") == "1",
		PetImage1:       c.FormValue("pet_image1"),
	}

	// Create PetInfo instance
	petInfo := models.PetInfo{
		ShelterID:       uint(shelterID),
		PetType:         requestBody.PetType,
		PetName:         requestBody.PetName,
		PetAge:          requestBody.PetAge,
		AgeType:         requestBody.AgeType,
		PetSex:          requestBody.PetSex,
		PetSize:         requestBody.PetSize,
		PriorityStatus:  requestBody.PriorityStatus,
		PetDescriptions: requestBody.PetDescriptions,
		CreatedAt:       time.Now(),
	}

	// Database transaction
	tx := middleware.DBConn.Begin()
	if err := tx.Create(&petInfo).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add pet",
		})
	}

	// Process Images
	petMedia := models.PetMedia{PetID: petInfo.PetID}
	petMedia.PetImage1 = processImage(c, "pet_image1", requestBody.PetImage1)

	if err := tx.Create(&petMedia).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add pet media",
		})
	}

	tx.Commit()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pet added successfully",
		"data": fiber.Map{
			"pet_info": petInfo,
			"image":    petMedia,
		},
	})
}

// processImage handles file uploads and Base64 strings
func processImage(c *fiber.Ctx, formKey, base64Str string) string {
	uploadedFile, err := c.FormFile(formKey)
	if err == nil {
		file, err := uploadedFile.Open()
		if err == nil {
			defer file.Close()
			fileBytes, _ := io.ReadAll(file)
			return base64.StdEncoding.EncodeToString(fileBytes)
		}
	}
	return base64Str
}

func FetchAndSearchPets(c *fiber.Ctx) error {
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
			PetID:     pet.PetID,
			PetType:   pet.PetType,
			PetName:   pet.PetName,
			PetSex:    pet.PetSex,
			PriorityStatus: pet.PriorityStatus,
			ShelterID: pet.ShelterID,
			PetImages: petMediaMap[pet.PetID],
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet search results",
		"data": fiber.Map{
			"pets": petResponses,
		},
	})
}

func FetchAndSearchArchivedPets(c *fiber.Ctx) error {
	// Get filters from query parameters
	shelterID := c.Params("id")
	petName := c.Query("pet_name")
	petSex := c.Query("sex")
	petType := c.Query("type")

	var pets []models.PetInfo

	query := middleware.DBConn.Where("status = ? AND shelter_id = ?", "archived", shelterID)

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

	result := query.Order("created_at DESC").Find(&pets)
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
			PetID:     pet.PetID,
			PetType:   pet.PetType,
			PetName:   pet.PetName,
			PetSex:    pet.PetSex,
			ShelterID: pet.ShelterID,
			PetImages: petMediaMap[pet.PetID],
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet search results",
		"data": fiber.Map{
			"pets": petResponses,
		},
	})
}

// Old fetching pet info
// func GetAllPetsInfoByShelterID(c *fiber.Ctx) error {
// 	shelterID := c.Params("id")

// 	// Fetch pet info for the given shelter
// 	var petInfo []models.PetInfo
// 	infoResult := middleware.DBConn.Where("shelter_id = ? AND status = ?", shelterID , "available").Order("created_at DESC").Find(&petInfo)

// 	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
// 		return c.JSON(response.ShelterResponseModel{
// 			RetCode: "404",
// 			Message: "Pet info not found",
// 			Data:    nil,
// 		})

// 	} else if infoResult.Error != nil {
// 		return c.JSON(response.ShelterResponseModel{
// 				RetCode: "500",
// 				Message: "Database error while fetching pet info",
// 				Data:    infoResult.Error,
// 			})
// 	}

// 	// Prepare a map to hold pet media by pet_id
// 	petMediaMap := make(map[uint][]string)

// 	// Fetch pet media for each pet and store them by pet_id
// 	var petMedia []models.PetMedia
// 	petmediaResult := middleware.DBConn.Where("pet_id IN ?", getPetIDs(petInfo)).Find(&petMedia)

// 	if petmediaResult.Error != nil && !errors.Is(petmediaResult.Error, gorm.ErrRecordNotFound) {
// 		return c.JSON(response.ShelterResponseModel{
// 			RetCode: "500",
// 			Message: "Database error while fetching pet media",
// 			Data:    petmediaResult.Error,
// 		})
// 	}

// 	// Map pet media by pet_id
// 	for _, media := range petMedia {
// 		// Ensure the image is valid Base64 before adding
// 		if media.PetImage1 != "" {
// 			_, err := base64.StdEncoding.DecodeString(media.PetImage1)
// 			if err == nil {
// 				petMediaMap[media.PetID] = append(petMediaMap[media.PetID], media.PetImage1)
// 			}
// 		}
// 	}

// 	// Combine pet info and media into a single response
// 	var petResponses []PetResponse
// 	for _, pet := range petInfo {
// 		// Create response for each pet by combining pet info and media
// 		petResponse := PetResponse{
// 			PetID:     pet.PetID,
// 			PetName:   pet.PetName,
// 			PetAge:    pet.PetAge,
// 			AgeType:   pet.AgeType,
// 			PetSex:    pet.PetSex,
// 			ShelterID: pet.ShelterID,
// 			PetImages: petMediaMap[pet.PetID], // Get media for this pet
// 		}
// 		petResponses = append(petResponses, petResponse)
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Shelter info retrieved successfully",
// 		"data": fiber.Map{
// 			"pets": petResponses,
// 		},
// 	})
// }

// Helper function to extract pet IDs from the petInfo slice
func getPetIDs(pets []models.PetInfo) []uint {
	var petIDs []uint
	for _, pet := range pets {
		petIDs = append(petIDs, pet.PetID)
	}
	return petIDs
}

func GetPetInfoByPetID(c *fiber.Ctx) error {
	// Get the pet_id from the URL params
	petID := c.Params("id")

	// Fetch pet info for the given pet_id
	var petInfo models.PetInfo
	infoResult := middleware.DBConn.Where("pet_id = ?", petID).First(&petInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "Pet info not found",
			Data:    nil,
		})

	} else if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet info",
			Data:    infoResult.Error,
		})
	}

	// Prepare a map to hold pet media by pet_id
	var petMedia []models.PetMedia
	petmediaResult := middleware.DBConn.Where("pet_id = ?", petID).Find(&petMedia)

	if petmediaResult.Error != nil && !errors.Is(petmediaResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet media",
			Data:    petmediaResult.Error,
		})
	}

	// Prepare pet media for this pet
	var petImages []string
	for _, media := range petMedia {
		// Ensure the image is valid Base64 before adding
		if media.PetImage1 != "" {
			_, err := base64.StdEncoding.DecodeString(media.PetImage1)
			if err == nil {
				petImages = append(petImages, media.PetImage1)
			}
		}
	}

	// Create a response for the pet by combining pet info and media
	petResponse := PetResponse{
		PetID:           petInfo.PetID,
		PetType:         petInfo.PetType,
		PetName:         petInfo.PetName,
		PetAge:          petInfo.PetAge,
		AgeType:         petInfo.AgeType,
		PetSex:          petInfo.PetSex,
		PetSize:         petInfo.PetSize,
		PetDescriptions: petInfo.PetDescriptions,
		PriorityStatus:  petInfo.PriorityStatus,
		ShelterID:       petInfo.ShelterID,
		PetImages:       petImages, // Attach pet images
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pet info retrieved successfully",
		"data": fiber.Map{
			"pet": petResponse,
		},
	})
}

func UpdatePetInfo(c *fiber.Ctx) error {
	petID := c.Params("id")

	// Fetch the existing PetInfo record
	var petInfo models.PetInfo
	infoResult := middleware.DBConn.Where("pet_id = ?", petID).First(&petInfo)
	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "Pet info not found",
			Data:    nil,
		})

	} else if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet info",
			Data:    infoResult.Error,
		})
	}

	// Get form values for text fields
	petInfo.PetName = c.FormValue("pet_name")
	petInfo.PetType = c.FormValue("pet_type")
	petInfo.PetSex = c.FormValue("pet_sex")
	petInfo.PetSize = c.FormValue("pet_size")
	petInfo.PetDescriptions = c.FormValue("pet_descriptions")
	petInfo.AgeType = c.FormValue("age_type")

	// Handle age field (assuming it's an integer)
	ageStr := c.FormValue("pet_age")
	petAge, err := strconv.Atoi(ageStr)
	if err != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Invalid pet age",
			Data:    err,
		})
	}
	petInfo.PetAge = petAge

	// Update PetInfo in the database
	middleware.DBConn.Table("petinfo").Where("pet_id = ?", petID).Updates(&petInfo)

	// Update PetMedia if image is provided
	fileHeader, err := c.FormFile("pet_image1")
	if err != nil {
		// If no image was uploaded, continue without updating pet image
		if err != fiber.ErrBadRequest {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "400",
				Message: "Failed to get image file",
				Data:    err,
			})
		}
	} else {
		// Open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "400",
				Message: "Failed to open image file",
				Data:    err,
			})
		}
		defer file.Close()

		// Read the file into a buffer
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "400",
				Message: "Failed to read image file",
				Data:    err,
			})
		}

		// Encode the image to base64
		base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())

		// Save to the petmedia table
		var petMedia models.PetMedia
		err = middleware.DBConn.Debug().Table("petmedia").Where("pet_id = ?", petID).First(&petMedia).Error
		if err != nil {
			// Handle missing record for petMedia
			petMedia.PetID = petInfo.PetID
		}

		// Update or insert the pet image
		petMedia.PetImage1 = base64Image
		middleware.DBConn.Table("petmedia").Where("pet_id = ?", petID).Updates(&petMedia)
	}

	// Return a successful response
	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Pet info and media updated successfully",
		Data:    petInfo,
	})
}

func UpdatePriorityStatus(c *fiber.Ctx) error {
	petID := c.Params("id")

	// Fetch the existing PetInfo record
	var petInfo models.PetInfo
	infoResult := middleware.DBConn.Where("pet_id = ?", petID).First(&petInfo)
	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "404",
			Message: "Pet info not found",
			Data:    nil,
		})

	} else if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet info",
			Data:    infoResult.Error,
		})
	}

	// Update the priority status
	petInfo.PriorityStatus = !petInfo.PriorityStatus // Toggle the priority status

	middleware.DBConn.Table("petinfo").Where("pet_id = ?", petID).Updates(&petInfo)

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Pet priority status updated successfully",
		Data:    petInfo,
	})
}

func SetPetStatusToArchive(c *fiber.Ctx) error {
	petID := c.Params("id") // Get pet ID from URL parameter

	// Check if pet exists in the database
	var pet models.PetInfo
	result := middleware.DBConn.First(&pet, petID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "404",
				Message: "Pet not found",
				Data:    nil,
			})
		}
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet",
			Data:    result.Error,
		})
	}

	// Check if the pet is already archived
	if pet.Status == "archived" {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Pet is already archived",
			Data:    nil,
		})
	}

	// Update pet status to 'archived'
	pet.Status = "archived"
	updateResult := middleware.DBConn.Save(&pet)

	if updateResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while updating pet status",
			Data:    updateResult.Error,
		})
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Pet status updated to archived successfully",
		Data:    pet,
	})
}

func SetPetStatusToUnarchive(c *fiber.Ctx) error {
	petID := c.Params("id") // Get pet ID from URL parameter

	// Check if pet exists in the database
	var pet models.PetInfo
	result := middleware.DBConn.First(&pet, petID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(response.ShelterResponseModel{
				RetCode: "404",
				Message: "Pet not found",
				Data:    nil,
			})
		}
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while fetching pet",
			Data:    result.Error,
		})
	}

	// Check if the pet is already archived
	if pet.Status == "available" {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "400",
			Message: "Pet is not archived",
			Data:    nil,
		})
	}

	// Update pet status to 'archived'
	pet.Status = "available"
	updateResult := middleware.DBConn.Save(&pet)

	if updateResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while updating pet status",
			Data:    updateResult.Error,
		})
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Pet status updated to unarchived successfully",
		Data:    pet,
	})
}

func DeletePetInfo(c *fiber.Ctx) error {
	petID := c.Params("id")

	// Delete the pet info
	infoResult := middleware.DBConn.Where("pet_id = ?", petID).Delete(&models.PetInfo{})
	if infoResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while deleting pet info",
			Data:    infoResult.Error,
		})
	}

	// Delete the pet media
	mediaResult := middleware.DBConn.Where("pet_id = ?", petID).Delete(&models.PetMedia{})
	if mediaResult.Error != nil {
		return c.JSON(response.ShelterResponseModel{
			RetCode: "500",
			Message: "Database error while deleting pet media",
			Data:    mediaResult.Error,
		})
	}

	return c.JSON(response.ShelterResponseModel{
		RetCode: "200",
		Message: "Pet info and media deleted successfully",
		Data:    nil,
	})
}
