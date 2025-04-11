package main

import (
	"fmt"
	"pethub_api/middleware"
	"pethub_api/models" // Import models to access them for AutoMigrate
	"pethub_api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func init() {
	fmt.Println("STARTING SERVER...")
	fmt.Println("INITIALIZE DB CONNECTION...")
	if middleware.ConnectDB() {
		fmt.Println("DB CONNECTION FAILED!")
	} else {
		fmt.Println("DB CONNECTION SUCCESSFUL!")
		// Auto-migrate the models (ensure the tables are created in the database)
		if err := middleware.DBConn.AutoMigrate(&models.AdoptionApplication{}, &models.Questionnaires{}); err != nil {
			fmt.Println("Migration failed:", err)
		} else {
			fmt.Println("Migration successful!")
		}
	}

	// Initialize session store
	store = session.New(session.Config{
		CookieHTTPOnly: true,
		CookieSecure:   false, // Set to true in production with HTTPS
		CookieSameSite: "Lax",
	})
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: middleware.GetEnv("PROJ_NAME"),
	})

	// Make session store accessible in middleware
	middleware.SessionStore = store

	// CORS CONFIG
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins (update for production)
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// LOGGER
	app.Use(logger.New())

	// SESSION middleware to attach session object to context
	app.Use(func(c *fiber.Ctx) error {
		sess, err := middleware.SessionStore.Get(c)
		if err != nil {
			return c.Status(500).SendString("Session error")
		}
		c.Locals("session", sess)
		return c.Next()
	})

	// API ROUTES
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(204) // No Content
	})

	// Call the route handler from the routes package
	routes.AppRoutes(app)

	// Start Server
	port := middleware.GetEnv("PROJ_PORT")
	if port == "" {
		port = "5566" // Default to port 5566 if not set
	}

	app.Listen(fmt.Sprintf("0.0.0.0:%s", port)) // Bind to all network interfaces
}
