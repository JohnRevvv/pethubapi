package routes

import (
	"pethub_api/controllers"
	"pethub_api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {

	pethubRoutes := app.Group("/api", middleware.JWTMiddleware())
	// ---------------- Admin Routes ----------------
	app.Post("/admin/register", controllers.RegisterAdmin)
	//app.Post("/admin/login", controllers.LoginAdmin)
	app.Get("/admin/getallpendingrequest", controllers.GetAllPendingRequests)
	app.Get("/admin/getalladopters", controllers.GetAllAdopters)
	app.Get("/admin/getallshelters", controllers.GetAllShelters)
	app.Post("/admin/updateregstatus", controllers.UpdateRegistrationStatus)
	app.Post("/admin/updateshelterstatus", controllers.UpdateShelterStatus)
	app.Post("/admin/updateadopterstatus", controllers.UpdateAdopterStatus)
	//try
	app.Get("/admin/getallshelterstry", controllers.GetAllSheltersAdmintry) // Route to get all shelters by id
	app.Put("/admin/shelters/:id/approve", controllers.ApproveShelterRegStatus)

	// =====================
	// Public Routes (No Auth Required)
	app.Post("/user/register", controllers.RegisterAdopter)
	app.Post("/user/login", controllers.LoginAdopter)
	// =====================
	// Protected Routes Group (Auth Required)// Adopter routes
	// =====================

	pethubRoutes.Get("/user", controllers.GetAllAdopters)
	pethubRoutes.Get("/user/:id", controllers.GetAdopterInfoByID)
	pethubRoutes.Get("/users/profile/:id", controllers.GetAdopterInfoByID)
	pethubRoutes.Put("/users/:id/update-info", controllers.UpdateAdopterDetails)
	pethubRoutes.Post("/users/:id/upload-media", controllers.UploadAdopterMedia)
	pethubRoutes.Get("/adopter/:id", controllers.GetAdopterProfile)
	pethubRoutes.Put("/adopter/:id", controllers.EditAdopterProfile)
	pethubRoutes.Post("/submission/:pet_id", controllers.CreateAdoptionSubmission)

	// Adopter - Pet Related
	pethubRoutes.Get("/users/petinfo", controllers.GetAllPets)
	pethubRoutes.Get("/users/pets/:id", controllers.GetPetByID)
	pethubRoutes.Get("/user/:id/pet", controllers.GetShelterWithPetsByID)

	// ---------------- Shelter Routes ----------------
	app.Post("/shelter/register", controllers.RegisterShelter)
	app.Post("/shelter/login", controllers.LoginShelter)
	pethubRoutes.Get("/shelters", controllers.GetAllShelters)
	pethubRoutes.Get("/shelter", controllers.GetShelterByName)
	pethubRoutes.Get("/shelter/:id", controllers.GetShelterInfoByID)
	pethubRoutes.Put("/shelter/:id/update-info", controllers.UpdateShelterDetails)
	pethubRoutes.Post("/shelter/:id/upload-media", controllers.UploadShelterMedia)
	// pethubRoutes.Post("/shelter/:id/add-pet-info", controllers.AddPetInfo)
	pethubRoutes.Get("/shelter/:id/petinfo", controllers.GetPetInfoByPetID)
	pethubRoutes.Put("/shelter/:id/update-pet-info", controllers.UpdatePetInfo)
	pethubRoutes.Put("/shelter/:id/archive-pet", controllers.SetPetStatusToArchive)
	pethubRoutes.Put("/shelter/:id/unarchive-pet", controllers.SetPetStatusToUnarchive)
	pethubRoutes.Get("/shelter/:id/petcount", controllers.CountPetsByShelter)
	pethubRoutes.Get("/filter/:id/pets/search", controllers.FetchAndSearchPets)
	pethubRoutes.Get("/shelter/archive/pets/:id/search", controllers.FetchAndSearchArchivedPets)
	pethubRoutes.Get("/shelter/:id/get/donationinfo", controllers.GetShelterDonationInfo)
	pethubRoutes.Put("/shelter/:id/update/donationinfo", controllers.UpdateShelterDonations)
	pethubRoutes.Put("/shelter/:id/change-password", controllers.ShelterChangePassword)
	pethubRoutes.Put("/shelter/:id/pet/update-priority-status", controllers.UpdatePriorityStatus)
	pethubRoutes.Get("/shelter/:id/adoption-applications", controllers.GetAdoptionApplications)
	pethubRoutes.Get("/shelter/:application_id/application-details", controllers.GetApplicationByApplicationID)
	pethubRoutes.Get("/shelterinfo/:shelter_id", controllers.GetShelterInfo)
	pethubRoutes.Post("/shelter/:shelter_id/add-pet", controllers.AddPetInfo2)

	// ---------------- General Shared Routes ----------------s
	app.Get("/allshelter", controllers.GetShelter)
	app.Get("/users/shelters/:id", controllers.GetAllSheltersByID)

	// ---------------- Forgot Password ----------------
	app.Post("/shelter/forgot-password", controllers.ShelterForgotPassword)
	app.Post("/shelter/verify-code", controllers.ShelterVerifyResetCode)
	app.Post("/shelter/reset-password", controllers.ShelterResetPassword)

	app.Post("/adopter/forgot-password", controllers.AdopterForgotPassword)
	app.Post("/adopter/verify-code", controllers.AdopterVerifyResetCode)
	app.Post("/adopter/reset-password", controllers.AdopterResetPassword)

	// ---------------- Adoption Application ----------------
	// app.Post("/adoption/application/:adopter_id/:pet_id", controllers.AdoptionApplication)
}
