package main

import (
	"fmt"
	"pethub_api/middleware"
	"pethub_api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	} else {
		fmt.Println(".env file loaded successfully")
	}
	fmt.Println("STARTING SERVER...")
	fmt.Println("INITIALIZE DB CONNECTION...")

	if middleware.ConnectDB() {
		fmt.Println("DB CONNECTION FAILED!")
	} else {
		fmt.Println("DB CONNECTION SUCCESSFUL!")

		// NOTE: AutoMigrate disabled to prevent schema changes in production
		// If you need to re-enable migration, uncomment the lines below:
		/*
			if err := middleware.DBConn.AutoMigrate(&models.AdoptionApplication{}, &models.Questionnaires{}); err != nil {
				fmt.Println("Migration failed:", err)
			} else {
				fmt.Println("Migration successful!")
			}
		*/
	}
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: middleware.GetEnv("PROJ_NAME"),
	})

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Logger
	app.Use(logger.New())

	// No session middleware required now ✅

	// Favicon route
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(204)
	})

	// Load app routes
	routes.AppRoutes(app)

	// Start server
	port := middleware.GetEnv("PROJ_PORT")
	if port == "" {
		port = "5566"
	}
	app.Listen(fmt.Sprintf("0.0.0.0:%s", port))
}
