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
	//app.Get("/admin/getallshelterstry", controllers.GetAllSheltersAdmintry) // Route to get all shelters by id
	//app.Put("/admin/shelters/:id/approve", controllers.ApproveShelterRegStatus)

	// =====================
	// Public Routes (No Auth Required)
	app.Post("/user/register", controllers.RegisterAdopter)
	app.Post("/user/login", controllers.LoginAdopter)
	// =====================
	// Protected Routes Group (Auth Required)// Adopter routes
	// =====================

	pethubRoutes.Get("/user", controllers.GetAllAdopters)
	pethubRoutes.Get("/user/:adopter_id", controllers.GetAdopterInfoOnly)
	pethubRoutes.Get("/users/profile/:id", controllers.GetAdopterInfoByID)
	pethubRoutes.Put("/users/:id/update-info", controllers.UpdateAdopterDetails)
	pethubRoutes.Post("/users/:id/upload-media", controllers.UploadAdopterMedia)
	pethubRoutes.Get("/adopter/:id", controllers.GetAdopterProfile)
	pethubRoutes.Post("/adopter/:adopter_id/edit", controllers.EditAdopterProfile)
	pethubRoutes.Post("/adopter/:shelter_id/:pet_id/:adopter_id/adoption", controllers.CreateAdoption)
	pethubRoutes.Get("/adopter/profile/:adopter_id", controllers.GetAdopterInfoByID)
	pethubRoutes.Get("/adopter/get/:shelter_id/other-pets", controllers.GetOtherPetsByAdopterID)

	// Adopter - Pet Related
	pethubRoutes.Get("/users/petinfo", controllers.GetAllPets)
	pethubRoutes.Get("/users/pets/:pet_id", controllers.GetPetByID)
	pethubRoutes.Get("/user/:id/pet", controllers.GetShelterWithPetsByID)
	pethubRoutes.Post("/users/status/:id", controllers.UpdatePetStatusToPending)
	pethubRoutes.Get("/users/priority/", controllers.GetPetsWithTrueStatus)
	pethubRoutes.Get("/users/allpets", controllers.GetAllPets)
	pethubRoutes.Get("/users/pets/search/all", controllers.FetchAllPets)
	pethubRoutes.Get("/applications/adopter/:adopter_id", controllers.GetApplicationByAdopterID)
	pethubRoutes.Get("/applications/pet/:pet_id", controllers.GetAdoptionApplicationsByPetID2)
	pethubRoutes.Get("/applications/status/:application_id", controllers.GetAdoptionSubmissionStatusByApplicationID)
	pethubRoutes.Post("/reports/shelter/:shelter_id/adopter/:adopter_id", controllers.SubmitReport)
	pethubRoutes.Get("/applications/allpets/:adopter_id", controllers.ShowPetsByAdopterID)
	pethubRoutes.Get("/adopter/:adopter_id/notifications", controllers.GetAdoptionNotifications)
	

	// ---------------- Shelter Routes ----------------
	app.Post("/shelter/register", controllers.RegisterShelter)
	app.Post("/shelter/login", controllers.LoginShelter)
	pethubRoutes.Get("/shelters", controllers.GetAllShelters)
	pethubRoutes.Get("/shelter", controllers.GetShelterByName)
	pethubRoutes.Get("/shelter/:id", controllers.GetShelterInfoByID)
	pethubRoutes.Put("/shelter/:id/update-info", controllers.UpdateShelterDetails)
	pethubRoutes.Post("/shelter/:id/upload-media", controllers.UploadShelterMedia)
	pethubRoutes.Get("/shelter/:shelter_id/refined", controllers.GetShelterDetailsRefined)
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
	pethubRoutes.Get("/shelter/:application_id/application-details", controllers.GetApplicationByApplicationID)
	pethubRoutes.Get("/shelterinfo/:shelter_id", controllers.GetShelterInfo)
	pethubRoutes.Post("/shelter/:shelter_id/add-pet", controllers.AddPetInfo)
	pethubRoutes.Get("/shelter/count/:pet_id/applied", controllers.CountApplicantsByPetId)

	pethubRoutes.Get("/shelter/:shelter_id/adoption", controllers.GetPetsWithAdoptionRequestsByShelter)
	pethubRoutes.Get("/shelter/:pet_id/get/applications", controllers.GetAdoptionApplicationsByPetID)
	pethubRoutes.Post("/shelter/application/:application_id/set-interview-date", controllers.SetInterviewSchedule)
	pethubRoutes.Put("/shelter/application/:application_id/interview/reject", controllers.RejectApplication)
	pethubRoutes.Get("/shelter/:shelter_id/adoption-applications", controllers.GetAdoptionSubmissionsByShelterAndStatus)
	pethubRoutes.Post("/shelter/reject-application/:application_id", controllers.RejectApplication)
	pethubRoutes.Put("/shelter/approve-application/:application_id", controllers.ApproveApplication)
	pethubRoutes.Get("/shelter/export/:shelter_id/:application_id/letter", controllers.GetInfosForDownloadLetter)

	// ---------------- General Shared Routes ----------------s
	pethubRoutes.Get("/allshelter", controllers.GetShelters)
	pethubRoutes.Get("/get/all/shelters", controllers.GetAllShelters)
	pethubRoutes.Get("/users/shelters/:id", controllers.GetAllSheltersByID)

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
