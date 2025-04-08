package models

import "time"

type AdoptionRequest struct {
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Address          string    `json:"address"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email"`
	Occupation       string    `json:"occupation"`
	SocialMedia      string    `json:"social_media"`
	CivilStatus      string    `json:"civil_status"`
	Sex              string    `json:"sex"`
	Birthdate        time.Time `json:"birthdate"`
	HasAdoptedBefore string    `json:"has_adopted_before"`
	IdealPet         string    `json:"ideal_pet"`
	BuildingType     string    `json:"building_type"`
	RentStatus       string    `json:"rent_status"`
	MovePlan         string    `json:"move_plan"`
	LivingWith       string    `json:"living_with"`
	Allergy          string    `json:"allergy"`
	PetCare          string    `json:"pet_care"`
	PetNeeds         string    `json:"pet_needs"`
	VacationPlan     string    `json:"vacation_plan"`
	FamilySupport    string    `json:"family_support"`
	HasOtherPets     string    `json:"has_other_pets"`
	HasPastPets      string    `json:"has_past_pets"`
	PetId            int       `json:"pet_id"`

	// Image uploads
	FrontOfHouse string `json:"front_of_house"`
	StreetPhoto  string `json:"street_photo"`
	LivingRoom   string `json:"living_room"`
	DiningArea   string `json:"dining_area"`
	Kitchen      string `json:"kitchen"`
	Bedroom      string `json:"bedroom"`
	HouseWindow  string `json:"window"`
	FrontYard    string `json:"front_yard"`
	Backyard     string `json:"backyard"`
}
