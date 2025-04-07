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
	app.Post("/user/login", controllers.LoginAdopter)    // Changed to POST for login
	app.Get("/user", controllers.GetAllAdopters)         // Route to get all adopters
	app.Get("/user/:id", controllers.GetAdopterInfoByID) // Route to get adopter by id
	app.Get("/users/petinfo", controllers.GetAllPets)    // Route to get all pets // user view all pet
	app.Get("/users/pets/:id", controllers.GetPetByID)   // Route to get pet by id

	// Shelter
	app.Post("/shelter/register", shelter.RegisterShelter)            // Route to register a shelter
	app.Post("/shelter/login", shelter.LoginShelter)                  // Changed to POST for login
	app.Get("/user/:id/pet", shelter.GetAllPetsInfoByShelterID)       // Route to get all pets by shelter ID
	app.Get("/shelters", controllers.GetAllShelters)                  // Route to get all shelters
	app.Get("/shelter", controllers.GetShelterByName)                 // Route to get shelter by name
	app.Get("/shelter/:id", shelter.GetShelterInfoByID)               // Route to get shelter by ID
	app.Put("/shelter/:id/update-info", shelter.UpdateShelterDetails) // Route to update shelter details
	app.Post("/shelter/:id/upload-media", shelter.UploadShelterMedia) // Route to upload or update shelter media
	app.Get("/shelter/:id/petinfo", shelter.GetPetInfoByPetID) // Route to get shelter media by ID
	app.Post("/shelter/:id/add-pet-info", shelter.AddPetInfo)
	app.Get("/shelter/:id/pets", shelter.GetAllPetsInfoByShelterID)
	app.Put("/shelter/:id/update-pet-info", shelter.UpdatePetInfo) // Route to update pet info

	app.Post("/shelter/:id/add-pet-info", shelter.AddPetInfo)
	app.Get("/allshelter", controllers.GetShelter)
	app.Get("/users/shelters/:id", controllers.GetAllSheltersByID) // shelters view all button
	app.Get("/users/profile/:id", controllers.GetAdopterInfoByID)
	app.Put("/users/:id/update-info", controllers.UpdateAdopterDetails)
	app.Post("/users/:id/upload-media", controllers.UploadAdopterMedia)

	app.Get("/adopter/:id", controllers.GetAdopterProfile) // Route to upload or update adopter media// Route to get adopter media by ID
	app.Put("/adopter/:id", controllers.EditAdopterProfile)
}
