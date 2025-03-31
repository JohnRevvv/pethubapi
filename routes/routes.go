package routes

import (
	"pethub_api/controllers"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {
	//Admin
	app.Post("/admin/register", controllers.RegisterAdmin) // Route to register an admin

	// Adopter
	app.Post("/user/register", controllers.RegisterAdopter)
	app.Post("/user/login", controllers.LoginAdopter)
	app.Get("/user", controllers.GetAllAdopters)     // Route to get all adopters
	app.Get("/user/:id", controllers.GetAdopterByID) // Route to get adopter by ID

	// Shelter
	app.Post("/shelter/register", controllers.RegisterShelter)
	app.Post("/shelter/login", controllers.LoginShelter)
	app.Post("/shelter/:id/pet", controllers.AddPet)                   // Route for adding pet info with shelter ID in URL
	app.Get("/shelter/:id/pets", controllers.GetAllPetsByShelterID)    // Route to get all pets by shelter ID
	app.Get("/shelters", controllers.GetAllShelters)                   // Route to get all shelters
	app.Get("/shelter", controllers.GetShelterByName)                  // Route to get shelter by name
	app.Get("/shelter/getinfo/:id", controllers.GetShelterInfoByID)    // Route to get shelter info by ID
	app.Get("/shelter/details/:id", controllers.GetShelterDetailsByID) // Route to get shelter details by ID
	app.Put("/shelter/update/:id", controllers.UpdateShelterDetails)   // Route to update shelter details
}
