package routes

import (
	"pethub_api/controllers"
	adopter "pethub_api/controllers/adopter"
	adoptionform "pethub_api/controllers/adoptionform"
	shelter "pethub_api/controllers/shelter"
	"pethub_api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {

	// Middleware for JWT validation

	// protcted routes?
	//jwt := app.Group("/api", middleware.ValidateJWTMiddleware)

	//app.Post("/adoption-application/:pet_id", middleware.ValidateJWTMiddleware, adoptionform.SubmitAdoptionApplication)

	//jwt.Post("/questionnaires", adoptionform.CreateQuestionnaire)

	// Adopter Routes
	app.Post("/user/register", adopter.RegisterAdopter)
	app.Post("/user/login", adopter.LoginAdopter)    // Changed to POST for login
	app.Get("/user", adopter.GetAllAdopters)         // Route to get all adopters
	app.Get("/user/:id", adopter.GetAdopterInfoByID) // Route to get adopter by id
	app.Get("/users/petinfo", adopter.GetAllPets)    // Route to get all pets // user view all pet
	app.Get("/users/pets/:id", adopter.GetPetByID)   // Route to get pet by id
	app.Get("/user/:id/pet", shelter.GetAllPetsInfoByShelterID)

	// Shelter
	app.Post("/shelter/register", controllers.RegisterShelter)            // Route to register a shelter
	app.Post("/shelter/login", controllers.LoginShelter)                  // Changed to POST for login      // Route to get all pets by shelter ID
	app.Get("/shelters", controllers.GetAllShelters)                      // Route to get all shelters
	app.Get("/shelter", controllers.GetShelterByName)                     // Route to get shelter by name
	app.Get("/shelter/:id", controllers.GetShelterInfoByID)               // Route to get shelter by ID
	app.Put("/shelter/:id/update-info", controllers.UpdateShelterDetails) // Route to update shelter details
	app.Post("/shelter/:id/upload-media", controllers.UploadShelterMedia) // Route to upload or update shelter media
	app.Get("/shelter/:id/petinfo", shelter.GetPetInfoByPetID)            // Route to get shelter media by ID
	app.Post("/shelter/:id/add-pet-info", shelter.AddPetInfo)
	app.Put("/shelter/:id/update-pet-info", shelter.UpdatePetInfo)     // Route to update pet info
	app.Put("/shelter/:id/archive-pet", shelter.SetPetStatusToArchive) // Route to archive pet info

	// search
	app.Get("/filter/:id/pets/search", shelter.FetchAndSearchPets) // Route to search pets by name

	// Pet Routes
	app.Post("/shelter/:id/add-pet-info", shelter.AddPetInfo)
	// app.Post("/shelter/:id/add-pet-media", shelter.AddPetMedia)

	// Questionnaire Routes
	//app.Post("/questionnaires", adoptionform.CreateQuestionnaire) // Route to submit a questionnaire

	// Adoptionform Routes
	//app.Post("/adoption-application/:pet_id", adoptionform.SubmitAdoptionApplication)

	// Route for getting specific adoption application by adopter_id
	//app.Get("/adoption-application/:adopter_id", adoptionform.GetAdoptionApplication)

	// Route for getting a specific questionanire form by adopter_id
	//app.Get("/questionnaire/:application_id", adoptionform.GetQuestionnaire)

	// Route for getting both adoption application and questionnaire form by adopter_id
	//app.Get("/adoption-and-questionnaire/:adopter_id", adoptionform.GetAdoptionApplicationAndQuestionnaire)

	// Define the route for updating adoption and questionnaire
	//app.Put("/updateAdoptionAndQuestionnaire", adoptionform.UpdateAdoptionAndQuestionnaire)

	app.Post("/shelter/:id/add-pet-info", shelter.AddPetInfo)
	app.Get("/allshelter", adopter.GetShelter)
	app.Get("/users/shelters/:id", adopter.GetAllSheltersByID) // shelters view all button
	app.Get("/users/profile/:id", adopter.GetAdopterInfoByID)
	app.Put("/users/:id/update-info", adopter.UpdateAdopterDetails)
	app.Post("/users/:id/upload-media", adopter.UploadAdopterMedia)

	app.Get("/adopter/:id", adopter.GetAdopterProfile) // Route to upload or update adopter media// Route to get adopter media by ID
	app.Put("/adopter/:id", adopter.EditAdopterProfile)

	// pakyu
	app.Post("/submission/:pet_id", middleware.ValidateJWTMiddleware, adoptionform.CreateAdoptionSubmission)

}
