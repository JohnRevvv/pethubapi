package routes

import (
	"pethub_api/controllers"
	shelter "pethub_api/controllers/shelter"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {
	// Admin
	app.Post("/admin/register", controllers.RegisterAdmin) // Route to register an admin

	// Adopter
	app.Post("/user/register", controllers.RegisterAdopter)
	app.Post("/user/login", controllers.LoginAdopter) // Changed to POST for login
	app.Get("/user", controllers.GetAllAdopters)      // Route to get all adopters
	app.Get("/user/:id", controllers.GetAdopterByID)  // Route to get adopter by ID

	// Shelter
	app.Post("/shelter/register", shelter.RegisterShelter)            // Route to register a shelter
	app.Post("/shelter/login", shelter.LoginShelter)                  // Changed to POST for login      
	app.Get("/shelter/:id/pets", shelter.GetAllPetsByShelterID)   // Route to get all pets by shelter ID
	app.Get("/shelters", controllers.GetAllShelters)                  // Route to get all shelters
	app.Get("/shelter", controllers.GetShelterByName)                 // Route to get shelter by name
	app.Get("/shelter/:id", shelter.GetShelterInfoByID)               // Route to get shelter by ID
	app.Put("/shelter/:id/update-info", shelter.UpdateShelterDetails)      // Route to update shelter details
	app.Post("/shelter/:id/upload-media", shelter.UploadShelterMedia) // Route to upload or update shelter media

	//pet
	app.Post("/shelter/:id/add-pet-info", shelter.AddPetInfo) 
	//app.Post("/shelter/:id/add-pet-media", shelter.AddPetMedia) 
}
