package routes

import (
	"pethub_api/controllers"
	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {
	// SAMPLE ENDPOINT
	app.Post("/registeradopter", controllers.RegisterAdopter)
	app.Get("/loginadopter", controllers.LoginAdopter)

	// CREATE YOUR ENDPOINTS HERE

	// --------------------------
}
