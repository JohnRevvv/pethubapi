package controllers

// import (
// 	"encoding/base64"
// 	"errors"
// 	"pethub_api/middleware"
// 	"pethub_api/models"

// 	"github.com/gofiber/fiber/v2"
// 	"gorm.io/gorm"
// )

// func GetShelterInfoByID(c *fiber.Ctx) error {
// 	shelterID := c.Params("id")

// 	// Fetch shelter info by ID
// 	var shelterInfo models.ShelterInfo
// 	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo)

// 	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"message": "Shelter info not found",
// 		})
// 	} else if infoResult.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Database error while fetching shelter info",
// 		})
// 	}

// 	// Fetch shelter media by ID
// 	var shelterMedia models.ShelterMedia
// 	mediaResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterMedia)

// 	// If no shelter media is found, set it to null
// 	var mediaResponse interface{}
// 	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
// 		mediaResponse = nil
// 	} else if mediaResult.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Database error while fetching shelter media",
// 		})
// 	} else {
// 		// Decode Base64-encoded images
// 		decodedProfile, err := base64.StdEncoding.DecodeString(shelterMedia.ShelterProfile)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Failed to decode profile image",
// 			})
// 		}

// 		decodedCover, err := base64.StdEncoding.DecodeString(shelterMedia.ShelterCover)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Failed to decode cover image",
// 			})
// 		}

// 		// Include decoded images in the response
// 		mediaResponse = fiber.Map{
// 			"shelter_profile": decodedProfile,
// 			"shelter_cover":   decodedCover,
// 		}
// 	}

// 	// Combine shelter info and media into a single response
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Shelter info retrieved successfully",
// 		"data": fiber.Map{
// 			"info":  shelterInfo,
// 			"media": mediaResponse,
// 		},
// 	})
// }

// func GetShelterDetailsByID(c *fiber.Ctx) error {
// 	shelterID := c.Params("id")

// 	// Fetch shelter info by ID
// 	var shelterInfo models.ShelterInfo
// 	infoResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterInfo)

// 	// If no shelter info is found, set it to null
// 	var infoResponse interface{}
// 	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
// 		infoResponse = nil
// 	} else if infoResult.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Database error while fetching shelter info",
// 		})
// 	} else {
// 		infoResponse = shelterInfo
// 	}

// 	// Fetch shelter media by ID
// 	var shelterMedia models.ShelterMedia
// 	mediaResult := middleware.DBConn.Where("shelter_id = ?", shelterID).First(&shelterMedia)

// 	// If no shelter media is found, set it to null
// 	var mediaResponse interface{}
// 	if errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
// 		mediaResponse = nil
// 	} else if mediaResult.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Database error while fetching shelter media",
// 		})
// 	} else {
// 		mediaResponse = shelterMedia
// 	}

// 	// Combine shelter info and media into a single response
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Shelter details retrieved successfully",
// 		"data": fiber.Map{
// 			"info":  infoResponse,
// 			"media": mediaResponse,
// 		},
// 	})
// }
