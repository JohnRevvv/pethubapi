package controllers

// import (
// 	"encoding/base64"
// 	"io/ioutil"
// 	"pethub_api/middleware"
// 	"pethub_api/models"
// 	"strconv"

// 	"github.com/gofiber/fiber/v2"
// )
// func AddPetMedia(c *fiber.Ctx) error {
// 	petID := c.Params("id")
// 	parsedPetID, err := strconv.ParseUint(petID, 10, 32)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Invalid pet ID",
// 		})
// 	}

// 	imageFile, err := c.FormFile("pet_image_1")
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "No image file uploaded",
// 		})
// 	}

// 	fileContent, err := imageFile.Open()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to open image file",
// 		})
// 	}
// 	defer fileContent.Close()

// 	fileBytes, err := ioutil.ReadAll(fileContent)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to read image file",
// 		})
// 	}

// 	encodedImage := base64.StdEncoding.EncodeToString(fileBytes)

// 	petMedia := models.PetMedia{
// 		PetID:     uint(parsedPetID),
// 		PetImage1: encodedImage,
// 	}

// 	if err := middleware.DBConn.Create(&petMedia).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to save pet media",
// 			"error":   err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "Pet media uploaded successfully",
// 		"data":    petMedia,
// 	})
// }
