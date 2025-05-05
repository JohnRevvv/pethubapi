package models

import (
	"time"
)

type AdoptionSubmission struct {
	// Application Info
	ApplicationID       uint   `json:"application_id" gorm:"primaryKey;autoIncrement"`
	ValidID             uint   `json:"valid_id"`
	ShelterID           uint   `json:"shelter_id"`
	PetID               uint   `json:"pet_id"`
	AdopterID           uint   `json:"adopter_id"`
	AltFName            string `json:"alt_f_name" gorm:"not null"`
	AltLName            string `json:"alt_l_name" gorm:"not null"`
	Relationship        string `json:"relationship" gorm:"not null"`
	AltContactNumber    string `json:"alt_contact_number" gorm:"not null"`
	AltEmail            string `json:"alt_email" gorm:"not null"`
	PetType             string `json:"pet_type" gorm:"not null"`
	ShelterAnimal       string `json:"shelter_animal"`
	IdealPetDescription string `json:"ideal_pet_description"`
	HousingSituation    string `json:"housing_situation" gorm:"not null"`
	PetsAtHome          string `json:"pets_at_home"`
	Allergies           string `json:"allergies"`
	FamilySupport       string `json:"family_support"`
	PastPets            string `json:"past_pets"`
	InterviewSetting    string `json:"interview_setting"`
	Status              string `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt           time.Time

	ValidIDs ValidIDPhotos `gorm:"foreignKey:ValidID;references:ValidID" json:"adoptionphotos"`
	Adopter  AdopterInfo   `json:"adopter"`
	Pet      PetInfo       `json:"pet"`
}

// TableName overrides default table name
func (AdoptionSubmission) TableName() string {
	return "adoption_submissions"
}

type ValidIDPhotos struct {
	ValidID        uint   `gorm:"primaryKey;autoIncrement"`
	AdopterIDType  string `json:"adopter_id_type"` // ✅ Ensure this exists and is exported
	AdopterValidID string `json:"adopter_valid_id"`
	AltIDType      string `json:"alt_id_type"`
	AltValidID     string `json:"alt_valid_id"`
}

// TableName overrides default table name
func (ValidIDPhotos) TableName() string {
	return "application_photos"
}
