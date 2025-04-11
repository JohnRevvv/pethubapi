package routes

import (
	"pethub_api/controllers"
	shelter "pethub_api/controllers/shelter"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {

	// Admin Routes
	app.Post("/admin/register", controllers.RegisterAdmin)
	app.Post("/admin/login", controllers.AdminLogin)
	app.Get("/admin/view-all-accounts", controllers.ViewAllAccounts)
	app.Put("/admin/update-status/:type/:id", controllers.UpdateAccountStatus)
	app.Delete("/admin/delete-account/:type/:id", controllers.DeleteAccount)

	// Adopter Routes
	app.Post("/user/register", controllers.RegisterAdopter)
	app.Post("/user/login", controllers.LoginAdopter)    // Changed to POST for login
	app.Get("/user", controllers.GetAllAdopters)         // Route to get all adopters
	app.Get("/user/:id", controllers.GetAdopterInfoByID) // Route to get adopter by id
	app.Get("/users/petinfo", controllers.GetAllPets)    // Route to get all pets // user view all pet
	app.Get("/users/pets/:id", controllers.GetPetByID)   // Route to get pet by id

	// Shelter Routes
	app.Post("/shelter/register", shelter.RegisterShelter)            // Route to register a shelter
	app.Post("/shelter/login", shelter.LoginShelter)                  // Changed to POST for login
	app.Get("/user/:id/pet", shelter.GetAllPetsInfoByShelterID)       // Route to get all pets by shelter ID
	app.Get("/shelters", controllers.GetAllShelters)                  // Route to get all shelters
	app.Get("/shelter", controllers.GetShelterByName)                 // Route to get shelter by name
	app.Get("/shelter/:id", shelter.GetShelterInfoByID)               // Route to get shelter by ID
	app.Put("/shelter/:id/update-info", shelter.UpdateShelterDetails) // Route to update shelter details
	app.Post("/shelter/:id/upload-media", shelter.UploadShelterMedia) // Route to upload or update shelter media
	app.Get("/shelter/:id/petinfo", shelter.GetPetInfoByPetID)        // Route to get shelter media by ID
	app.Post("/shelter/:id/add-pet-info", shelter.AddPetInfo)
	app.Get("/shelter/:id/pets", shelter.GetAllPetsInfoByShelterID)
	app.Put("/shelter/:id/update-pet-info", shelter.UpdatePetInfo) // Route to update pet info

	// Pet Routes
	app.Post("/shelter/:id/add-pet-info", shelter.AddPetInfo)
	// app.Post("/shelter/:id/add-pet-media", shelter.AddPetMedia)

	// Questionnaire Routes
	app.Post("/questionnaires", controllers.CreateQuestionnaire) // Route to submit a questionnaire

	// Adoptionform Routes
	app.Post("/adoption-application/:pet_id", controllers.SubmitAdoptionApplication)

	// Route for getting specific adoption application by adopter_id
	app.Get("/adoption-application/:adopter_id", controllers.GetAdoptionApplication)

	// Route for getting a specific questionanire form by adopter_id
	app.Get("/questionnaire/:application_id", controllers.GetQuestionnaire)

	// Route for getting both adoption application and questionnaire form by adopter_id
	app.Get("/adoption-and-questionnaire/:adopter_id", controllers.GetAdoptionApplicationAndQuestionnaire)

	// Define the route for updating adoption and questionnaire
	app.Put("/updateAdoptionAndQuestionnaire", controllers.UpdateAdoptionAndQuestionnaire)

	app.Post("/shelter/:id/add-pet-info", shelter.AddPetInfo)
	app.Get("/allshelter", controllers.GetShelter)
	app.Get("/users/shelters/:id", controllers.GetAllSheltersByID) // shelters view all button
	app.Get("/users/profile/:id", controllers.GetAdopterInfoByID)
	app.Put("/users/:id/update-info", controllers.UpdateAdopterDetails)
	app.Post("/users/:id/upload-media", controllers.UploadAdopterMedia)

	app.Get("/adopter/:id", controllers.GetAdopterProfile) // Route to upload or update adopter media// Route to get adopter media by ID
	app.Put("/adopter/:id", controllers.EditAdopterProfile)

}
