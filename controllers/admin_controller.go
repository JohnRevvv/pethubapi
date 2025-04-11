package controllers

import (
	"pethub_api/middleware"
	"pethub_api/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(c *fiber.Ctx) error {
	// Parse raw JSON input
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Validate input
	if loginRequest.Username == "" || loginRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	// Look for the admin account in the database
	var admin models.AdminAccount
	if err := middleware.DBConn.Where("username = ?", loginRequest.Username).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid admin credentials",
		})
	}

	// Compare the hashed password with the input password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid admin credentials",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"admin":   admin,
	})
}

func RegisterAdmin(c *fiber.Ctx) error {
	// Parse raw JSON input
	var adminRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&adminRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Check if both username and password are provided
	if adminRequest.Username == "" || adminRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	// Check if the admin account already exists
	var existingAdmin models.AdminAccount
	if err := middleware.DBConn.Where("username = ?", adminRequest.Username).First(&existingAdmin).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Admin account with this username already exists",
		})
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create a new admin account
	admin := models.AdminAccount{
		Username: adminRequest.Username,
		Password: string(hashedPassword), // Store the hashed password
	}

	// Save the admin account to the database
	if err := middleware.DBConn.Create(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create admin account",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Admin account created successfully",
		"admin":   admin,
	})
}

func ViewAllAccounts(c *fiber.Ctx) error {
	var adopters []models.AdopterAccount
	var shelters []models.ShelterAccount

	// Get all adopters
	if err := middleware.DBConn.Find(&adopters).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve adopters",
		})
	}

	// Get all shelters
	if err := middleware.DBConn.Find(&shelters).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve shelters",
		})
	}

	// Return the results
	return c.JSON(fiber.Map{
		"adopters": adopters,
		"shelters": shelters,
	})
}

// UpdateAccountStatus allows the admin to update the status of an adopter or shelter account
func UpdateAccountStatus(c *fiber.Ctx) error {
	accountType := c.Params("type")
	id := c.Params("id")

	var input struct {
		Status    string `json:"status"`
		RegStatus string `json:"reg_status"` // Add reg_status input for shelters
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if accountType == "adopter" {
		var adopter models.AdopterAccount
		if err := middleware.DBConn.Where("adopter_id = ?", id).First(&adopter).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Adopter not found"})
		}
		if input.Status != "" {
			adopter.Status = input.Status
		}
		if err := middleware.DBConn.Save(&adopter).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update adopter status"})
		}
	} else if accountType == "shelter" {
		var shelter models.ShelterAccount
		if err := middleware.DBConn.Where("shelter_id = ?", id).First(&shelter).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Shelter not found"})
		}
		if input.Status != "" {
			shelter.Status = input.Status
		}
		if input.RegStatus != "" {
			shelter.RegStatus = input.RegStatus
		}
		if err := middleware.DBConn.Save(&shelter).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update shelter status"})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid account type"})
	}

	return c.JSON(fiber.Map{"message": "Account status updated successfully"})
}

// DeleteAccount handles the deletion of an adopter or shelter account
func DeleteAccount(c *fiber.Ctx) error {
	accountType := c.Params("type")
	id := c.Params("id")

	// Check if accountType is valid
	if accountType != "adopter" && accountType != "shelter" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account type",
		})
	}

	// Delete the corresponding account based on the account type
	if accountType == "adopter" {
		var adopter models.AdopterAccount
		if err := middleware.DBConn.Where("adopter_id = ?", id).First(&adopter).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Adopter not found",
			})
		}

		// Delete adopter-related records in other tables (if needed) and then delete the adopter account
		if err := middleware.DBConn.Delete(&adopter).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete adopter account",
			})
		}
	} else if accountType == "shelter" {
		var shelterInfo models.ShelterInfo
		var shelterAccount models.ShelterAccount

		// Delete shelterinfo record first to avoid foreign key violation
		if err := middleware.DBConn.Where("shelter_id = ?", id).First(&shelterInfo).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Shelter info not found",
			})
		}

		if err := middleware.DBConn.Delete(&shelterInfo).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete shelter info",
			})
		}

		// Now, delete the shelter account record
		if err := middleware.DBConn.Where("shelter_id = ?", id).First(&shelterAccount).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Shelter account not found",
			})
		}

		if err := middleware.DBConn.Delete(&shelterAccount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete shelter account",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid account type",
		})
	}

	// If all deletions are successful, return a success message
	return c.JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}
