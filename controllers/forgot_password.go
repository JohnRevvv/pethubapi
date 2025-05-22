package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"sync"
	"time"

	"pethub_api/middleware"
	"pethub_api/models"
	"pethub_api/models/response"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendEmail(toEmail string, code string) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	subject := "Password Reset Verification Code"
	body := fmt.Sprintf(`
		Your password reset verification code is: %s

		This code will expire in 5 minutes.

		If you did not request this, please ignore the message.
	`, code)

	msg := []byte("Subject: " + subject + "\r\n\r\n" + body)
	to := []string{toEmail}

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		log.Println("Send email error:", err)
	}
	return err
}

var resetCodeStore = struct {
	sync.RWMutex
	codes map[string]struct {
		Code      string
		ExpiresAt time.Time
	}
}{codes: make(map[string]struct {
	Code      string
	ExpiresAt time.Time
})}

// ForgotPassword - step 1
func ShelterForgotPassword(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"shelter_email"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	var shelterInfo models.ShelterInfo
	result := middleware.DBConn.Where("shelter_email = ?", req.Email).First(&shelterInfo)
	if result.Error == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(response.ResponseModel{
			RetCode: "404", Message: "Email not found",
		})
	} else if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500", Message: "Database error",
		})
	}

	code := GenerateRandomCode()
	if err := SendEmail(req.Email, code); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500", Message: "Failed to send email",
		})
	}

	resetCodeStore.Lock()
	resetCodeStore.codes[req.Email] = struct {
		Code      string
		ExpiresAt time.Time
	}{Code: code, ExpiresAt: time.Now().Add(5 * time.Minute)}
	resetCodeStore.Unlock()

	return c.JSON(response.ResponseModel{
		RetCode: "200", Message: "Verification code sent to your email",
	})
}

// VerifyResetCode - step 2
func ShelterVerifyResetCode(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"shelter_email"`
		Code  string `json:"code"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	resetCodeStore.RLock()
	data, exists := resetCodeStore.codes[req.Email]
	resetCodeStore.RUnlock()

	if !exists || time.Now().After(data.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ResponseModel{
			RetCode: "401", Message: "Code expired or not found",
		})
	}

	if data.Code != req.Code {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ResponseModel{
			RetCode: "401", Message: "Incorrect code",
		})
	}

	resetCodeStore.Lock()
	resetCodeStore.codes[req.Email] = struct {
		Code      string
		ExpiresAt time.Time
	}{"VERIFIED", time.Now().Add(15 * time.Minute)}
	resetCodeStore.Unlock()

	return c.JSON(response.ResponseModel{
		RetCode: "200", Message: "Code verified. You can now reset your password.",
	})
}

// ResetPassword - step 3
func ShelterResetPassword(c *fiber.Ctx) error {
	type Request struct {
		Email           string `json:"shelter_email"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400", Message: "Passwords do not match",
		})
	}

	resetCodeStore.RLock()
	data, exists := resetCodeStore.codes[req.Email]
	resetCodeStore.RUnlock()

	if !exists || data.Code != "VERIFIED" || time.Now().After(data.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ResponseModel{
			RetCode: "401", Message: "Unauthorized or session expired",
		})
	}

	var shelterInfo models.ShelterInfo
	if err := middleware.DBConn.Where("shelter_email = ?", req.Email).First(&shelterInfo).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ResponseModel{
			RetCode: "404", Message: "User not found",
		})
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500", Message: "Password encryption failed",
		})
	}

	if err := middleware.DBConn.Model(&models.ShelterAccount{}).
		Where("shelter_id = ?", shelterInfo.ShelterID).
		Update("password", string(hashedPwd)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Password update failed",
			Data:    err.Error(),
		})
	}

	resetCodeStore.Lock()
	delete(resetCodeStore.codes, req.Email)
	resetCodeStore.Unlock()

	return c.JSON(response.ResponseModel{
		RetCode: "200", Message: "Password has been reset successfully",
	})
}

// *****************************************************************************
// ******************************** ADOPTER ************************************
// *****************************************************************************
// ForgotPassword - step 1
func AdopterForgotPassword(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	var adopterInfo models.AdopterInfo
	result := middleware.DBConn.Where("email = ?", req.Email).First(&adopterInfo)
	if result.Error == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(response.ResponseModel{
			RetCode: "404", Message: "Email not found",
		})
	} else if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500", Message: "Database error",
		})
	}

	code := GenerateRandomCode()
	if err := SendEmail(req.Email, code); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500", Message: "Failed to send email",
		})
	}

	resetCodeStore.Lock()
	resetCodeStore.codes[req.Email] = struct {
		Code      string
		ExpiresAt time.Time
	}{Code: code, ExpiresAt: time.Now().Add(5 * time.Minute)}
	resetCodeStore.Unlock()

	return c.JSON(response.ResponseModel{
		RetCode: "200", Message: "Verification code sent to your email",
	})
}

// VerifyResetCode - step 2
func AdopterVerifyResetCode(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	resetCodeStore.RLock()
	data, exists := resetCodeStore.codes[req.Email]
	resetCodeStore.RUnlock()

	if !exists || time.Now().After(data.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ResponseModel{
			RetCode: "401", Message: "Code expired or not found",
		})
	}

	if data.Code != req.Code {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ResponseModel{
			RetCode: "401", Message: "Incorrect code",
		})
	}

	resetCodeStore.Lock()
	resetCodeStore.codes[req.Email] = struct {
		Code      string
		ExpiresAt time.Time
	}{"VERIFIED", time.Now().Add(15 * time.Minute)}
	resetCodeStore.Unlock()

	return c.JSON(response.ResponseModel{
		RetCode: "200", Message: "Code verified. You can now reset your password.",
	})
}

// ResetPassword - step 3
func AdopterResetPassword(c *fiber.Ctx) error {
	type Request struct {
		Email           string `json:"email"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	if req.NewPassword != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400", Message: "Passwords do not match",
		})
	}

	resetCodeStore.RLock()
	data, exists := resetCodeStore.codes[req.Email]
	resetCodeStore.RUnlock()

	if !exists || data.Code != "VERIFIED" || time.Now().After(data.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ResponseModel{
			RetCode: "401", Message: "Unauthorized or session expired",
		})
	}

	var adopterInfo models.AdopterInfo
	if err := middleware.DBConn.Where("email = ?", req.Email).First(&adopterInfo).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ResponseModel{
			RetCode: "404", Message: "User not found",
		})
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500", Message: "Password encryption failed",
		})
	}

	if err := middleware.DBConn.Model(&models.AdopterAccount{}).
		Where("adopter_id = ?", adopterInfo.AdopterID).
		Update("password", string(hashedPwd)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Password update failed",
			Data:    err.Error(),
		})
	}

	resetCodeStore.Lock()
	delete(resetCodeStore.codes, req.Email)
	resetCodeStore.Unlock()

	return c.JSON(response.ResponseModel{
		RetCode: "200", Message: "Password has been reset successfully",
	})
}
