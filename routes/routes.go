package routes

import (
	"pethub_api/controllers"
	// shelter "pethub_api/controllers/shelter"

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
	app.Get("/user/:id/pet", controllers.GetAllPetsInfoByShelterID) 

	// Shelter
	app.Post("/shelter/register", controllers.RegisterShelter)            // Route to register a shelter
	app.Post("/shelter/login", controllers.LoginShelter)                  // Changed to POST for login 
	app.Get("/shelters", controllers.GetAllShelters)                  // Route to get all shelters
	app.Get("/shelter", controllers.GetShelterByName)                 // Route to get shelter by name
	app.Get("/shelter/:id", controllers.GetShelterInfoByID)               // Route to get shelter by ID
	app.Put("/shelter/:id/update-info", controllers.UpdateShelterDetails) // Route to update shelter details
	app.Post("/shelter/:id/upload-media", controllers.UploadShelterMedia) // Route to upload or update shelter media
	app.Get("/shelter/:id/petinfo", controllers.GetPetInfoByPetID) // Route to get shelter media by ID
	app.Put("/shelter/:id/update-pet-info", controllers.UpdatePetInfo) // Route to update pet info
	app.Put("/shelter/:id/archive-pet", controllers.SetPetStatusToArchive) // Route to archive pet info
	app.Put("/shelter/:id/unarchive-pet", controllers.SetPetStatusToUnarchive) // Route to unarchive pet info
	app.Get("/filter/:id/pets/search", controllers.FetchAndSearchPets) // Route to search pets by name
	app.Get("/shelter/archive/pets/:id/search", controllers.FetchAndSearchArchivedPets) // Route to get archived pets
	app.Get("/shelter/:id/petcount", controllers.CountPetsByShelter) // Route to get pet count by shelter ID


	app.Post("/shelter/:id/add-pet-info", controllers.AddPetInfo)
	app.Get("/allshelter", controllers.GetShelter)
	app.Get("/users/shelters/:id", controllers.GetAllSheltersByID) // shelters view all button
	app.Get("/users/profile/:id", controllers.GetAdopterInfoByID)
	app.Put("/users/:id/update-info", controllers.UpdateAdopterDetails)
	app.Post("/users/:id/upload-media", controllers.UploadAdopterMedia)

	app.Get("/adopter/:id", controllers.GetAdopterProfile) // Route to upload or update adopter media// Route to get adopter media by ID
	app.Put("/adopter/:id", controllers.EditAdopterProfile)

	//FORGOT PASSWORD
	app.Post("/shelter/forgot-password", controllers.ShelterForgotPassword) // Route to handle forgot password
	app.Post("/shelter/verify-code", controllers.ShelterVerifyResetCode)       // Route to verify code
	app.Post("/shelter/reset-password", controllers.ShelterResetPassword) // Route to reset password

	app.Post("/adopter/forgot-password", controllers.AdopterForgotPassword) // Route to handle forgot password
	app.Post("/adopter/verify-code", controllers.AdopterVerifyResetCode)       // Route to verify code
	app.Post("/adopter//reset-password", controllers.AdopterResetPassword) // Route to reset password
	
}
