package middleware

import (
	"fmt"
	"pethub_api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
	DBErr  error
)

// ConnectDB initializes the connection to the PostgreSQL database using
// environment variables for configuration and assigns the connection
// to the global variable DBConn. It returns true if there was an error
// establishing the connection, otherwise false.
func ConnectDB() bool {
	// Database Confg
	dns := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s TimeZone=%s",
		GetEnv("DB_HOST"), GetEnv("DB_PORT"), GetEnv("DB_NAME"),
		GetEnv("DB_UNME"), GetEnv("DB_PWRD"), GetEnv("DB_SSLM"),
		GetEnv("DB_TMEZ"))

	DBConn, DBErr = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if DBErr != nil {
		fmt.Printf("Database connection error: %v\n", DBErr) // Debugging log
		return true
	}

	// Debugging log to confirm successful connection
	fmt.Println("Database connection established successfully")

	// Auto-migrate models
	DBConn.AutoMigrate(&models.AdopterAccount{},
		&models.AdopterInfo{},
		&models.ShelterAccount{},
		&models.ShelterInfo{},
		&models.ShelterMedia{},
		&models.PetInfo{},
		&models.PetMedia{})
	return false
}
