package routes

import (
	"pethub_api/controllers"

	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {
	// ---------------- Admin Routes ----------------
	app.Post("/admin/register", controllers.RegisterAdmin)
	app.Post("/admin/login", controllers.LoginAdmin)
	app.Get("/admin/getallpendingrequest", controllers.GetAllPendingRequests)
	app.Get("/admin/getalladopters", controllers.GetAllAdopters)
	app.Get("/admin/getallshelters", controllers.GetAllShelters)
	app.Post("/admin/updateregstatus", controllers.UpdateRegistrationStatus)
	app.Post("/admin/updateshelterstatus", controllers.UpdateShelterStatus)
	app.Post("/admin/updateadopterstatus", controllers.UpdateAdopterStatus)

	//try
	app.Get("/admin/getallshelterstry", controllers.GetAllSheltersAdmintry) // Route to get all shelters by id
	app.Put("/admin/shelters/:id/approve", controllers.ApproveShelterRegStatus)

	// ---------------- Adopter Routes ----------------
	app.Post("/user/register", controllers.RegisterAdopter)
	app.Post("/user/login", controllers.LoginAdopter)
	app.Get("/user", controllers.GetAllAdopters)
	app.Get("/user/:id", controllers.GetAdopterInfoByID)
	app.Get("/users/profile/:id", controllers.GetAdopterInfoByID)
	app.Put("/users/:id/update-info", controllers.UpdateAdopterDetails)
	app.Post("/users/:id/upload-media", controllers.UploadAdopterMedia)
	app.Get("/adopter/:id", controllers.GetAdopterProfile)
	app.Put("/adopter/:id", controllers.EditAdopterProfile)

	// Adopter - Pet Related
	app.Get("/users/petinfo", controllers.GetAllPets)
	app.Get("/users/pets/:id", controllers.GetPetByID)
	app.Get("/user/:id/pet", controllers.GetAllPetsInfoByShelterID)

	// ---------------- Shelter Routes ----------------
	app.Post("/shelter/register", controllers.RegisterShelter)
	app.Post("/shelter/login", controllers.LoginShelter)
	app.Get("/shelters", controllers.GetAllShelters)
	app.Get("/shelter", controllers.GetShelterByName)
	app.Get("/shelter/:id", controllers.GetShelterInfoByID)
	app.Put("/shelter/:id/update-info", controllers.UpdateShelterDetails)
	app.Post("/shelter/:id/upload-media", controllers.UploadShelterMedia)
	app.Post("/shelter/:id/add-pet-info", controllers.AddPetInfo)
	app.Get("/shelter/:id/petinfo", controllers.GetPetInfoByPetID)
	app.Put("/shelter/:id/update-pet-info", controllers.UpdatePetInfo)
	app.Put("/shelter/:id/archive-pet", controllers.SetPetStatusToArchive)
	app.Put("/shelter/:id/unarchive-pet", controllers.SetPetStatusToUnarchive)
	app.Get("/shelter/:id/petcount", controllers.CountPetsByShelter)
	app.Get("/filter/:id/pets/search", controllers.FetchAndSearchPets)
	app.Get("/shelter/archive/pets/:id/search", controllers.FetchAndSearchArchivedPets)
	app.Get("/shelter/:id/get/donationinfo", controllers.GetShelterDonationInfo)
	app.Put("/shelter/:id/update/donationinfo", controllers.UpdateShelterDonations)
	app.Put("/shelter/:id/change-password", controllers.ShelterChangePassword)
	app.Put("/shelter/:pet_id/pet/update-priority-status", controllers.UpdatePriorityStatus)

	// ---------------- General Shared Routes ----------------
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
	app.Post("/adoption/application/:adopter_id/:pet_id", controllers.AdoptionApplication)
}
