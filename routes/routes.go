package routes

import (
	"pethub_api/controllers"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {
	//Admin

	// Adopter
	app.Post("/user/register", controllers.RegisterAdopter)
	app.Get("/user/login", controllers.LoginAdopter)

	// Shelter
	app.Post("/shelter/register", controllers.RegisterShelter)
	app.Get("/shelter/login", controllers.LoginShelter)
	app.Post("/shelter/:id/pet", controllers.AddPet) // Route for adding pet info with shelter ID in URL
}
