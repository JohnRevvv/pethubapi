package controllers

import (
	"encoding/base64"

	"io/ioutil"
	"pethub_api/middleware"
	"pethub_api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddPetMedia(c *fiber.Ctx) error {
	petID := c.Params("id")
	parsedPetID, err := strconv.ParseUint(petID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid pet ID",
		})
	}

	var petMedia models.PetMedia
	petMedia.PetID = uint(parsedPetID)

	for i := 1; i <= 4; i++ {
		imageFile, err := c.FormFile("pet_image_" + strconv.Itoa(i))
		if err == nil {
			fileContent, err := imageFile.Open()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to open pet image " + strconv.Itoa(i),
				})
			}
			defer fileContent.Close()

			fileBytes, err := ioutil.ReadAll(fileContent)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to read pet image " + strconv.Itoa(i),
				})
			}

			encodedImage := base64.StdEncoding.EncodeToString(fileBytes)

			switch i {
			case 1:
				petMedia.PetImage1 = encodedImage
				// case 2:
				//     petMedia.PetImage2 = encodedImage
				//  case 3:
				//    petMedia.PetImage3 = encodedImage
				// case 4:
				// petMedia.PetImage4 = encodedImage
			}
		}
	}

	if err := middleware.DBConn.Create(&petMedia).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to save pet media",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pet media uploaded successfully",
		"data":    petMedia,
	})
}
