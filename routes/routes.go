package routes

import (
	"pethub_api/controllers"
	shelter "pethub_api/controllers/shelter"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {

	// Admin Routes
	app.Post("/admin/register", controllers.RegisterAdmin) // Route to register an admin

	// Adopter Routes
	app.Post("/user/register", controllers.RegisterAdopter)
	app.Post("/user/login", controllers.LoginAdopter) // Changed to POST for login
	app.Get("/user", controllers.GetAllAdopters)      // Route to get all adopters
	app.Get("/user/:id", controllers.GetAdopterByID)  // Route to get adopter by ID

	// Shelter Routes
	app.Post("/shelter/register", shelter.RegisterShelter)            // Route to register a shelter
	app.Post("/shelter/login", shelter.LoginShelter)                  // Changed to POST for login
	app.Get("/shelter/:id/pets", shelter.GetAllPetsByShelterID)       // Route to get all pets by shelter ID
	app.Get("/shelters", controllers.GetAllShelters)                  // Route to get all shelters
	app.Get("/shelter", controllers.GetShelterByName)                 // Route to get shelter by name
	app.Get("/shelter/:id", shelter.GetShelterInfoByID)               // Route to get shelter by ID
	app.Put("/shelter/:id/update-info", shelter.UpdateShelterDetails) // Route to update shelter details
	app.Post("/shelter/:id/upload-media", shelter.UploadShelterMedia) // Route to upload or update shelter media

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

}
