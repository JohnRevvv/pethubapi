package controllers

import (
	"errors"
	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllShelters(c *fiber.Ctx) error {
	// Fetch all shelter accounts
	var shelterAccounts []models.ShelterAccount
	accountResult := middleware.DBConn.Find(&shelterAccounts)

	if accountResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter accounts",
		})
	}

	// Fetch all shelter info
	var shelterInfos []models.ShelterInfo
	infoResult := middleware.DBConn.Find(&shelterInfos)

	if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter info",
		})
	}

	// Combine accounts and info into a single response
	shelters := []fiber.Map{}
	for _, account := range shelterAccounts {
		for _, info := range shelterInfos {
			if account.ShelterID == info.ShelterID {
				shelters = append(shelters, fiber.Map{
					"shelter": account,
					"info":    info,
				})
				break
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelters retrieved successfully",
		"data":    shelters,
	})
}

func GetShelterByName(c *fiber.Ctx) error {
	shelterName := c.Query("shelter_name")
	if shelterName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Shelter name query parameter is missing",
		})
	}

	// Fetch shelter info by name
	var shelterInfo models.ShelterInfo
	infoResult := middleware.DBConn.Where("shelter_name = ?", shelterName).First(&shelterInfo)

	if errors.Is(infoResult.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Shelter not found",
		})
	} else if infoResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
		})
	}

	// Fetch shelter account associated with the shelter info
	var shelterAccount models.ShelterAccount
	accountResult := middleware.DBConn.Where("id = ?", shelterInfo.ShelterID).First(&shelterAccount)

	if accountResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch shelter account",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shelter retrieved successfully",
		"data": fiber.Map{
			"shelter": shelterAccount,
			"info":    shelterInfo,
		},
	})
}