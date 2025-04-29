package models

import (
	"time"
)

type AdoptionSubmission struct {
	// Application Info
<<<<<<< HEAD
	ApplicationID    uint   `json:"application_id" gorm:"primaryKey;autoIncrement"`
	PetID            uint   `json:"pet_id" gorm:"not null"`
	ShelterID        uint   `json:"shelter_id" gorm:"not null"`
	AdopterID        uint   `json:"adopter_id" gorm:"not null"`
	AltFName         string `json:"alt_f_name" gorm:"not null"`
	AltLName         string `json:"alt_l_name" gorm:"not null"`
	Relationship     string `json:"relationship" gorm:"not null"`
	AltContactNumber string `json:"alt_contact_number" gorm:"not null"`
	AltEmail         string `json:"alt_email" gorm:"not null"`

	// Questionnaire Info
	PetType             string `json:"pet_type" gorm:"not null"`
	ShelterAnimal       string `json:"shelter_animal"`
	IdealPetDescription string `json:"ideal_pet_description"`
	HousingSituation    string `json:"housing_situation" gorm:"not null"`
	PetsAtHome          string `json:"pets_at_home"`
	Allergies           string `json:"allergies"`
	FamilySupport       string `json:"family_support"`
	PastPets            string `json:"past_pets"`
	InterviewSetting    string `json:"interview_setting"`

	// Upload fields (just store file paths)
	ValidID    string `json:"valid_id" gorm:"not null"`     // For adopter
	AltValidID string `json:"alt_valid_id" gorm:"not null"` // For alternate contact

	// Add relationships (below existing fields but above Common Fields)
	Pet     PetInfo        `gorm:"foreignKey:PetID" json:"pet"`
	Shelter ShelterAccount `gorm:"foreignKey:ShelterID" json:"shelter"`

	// Common Fields
	Status    string    `json:"status" gorm:"type:varchar(20);default:'Pending'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time `json:"deleted_at"`
=======
	ApplicationID       uint      `json:"application_id" gorm:"primaryKey;autoIncrement"`
	ShelterID           uint      `json:"shelter_id"`
	PetID               uint      `json:"pet_id"`
	AdopterID           uint      `json:"adopter_id"`
	AltFName            string    `json:"alt_f_name" gorm:"not null"`
	AltLName            string    `json:"alt_l_name" gorm:"not null"`
	Relationship        string    `json:"relationship" gorm:"not null"`
	AltContactNumber    string    `json:"alt_contact_number" gorm:"not null"`
	AltEmail            string    `json:"alt_email" gorm:"not null"`
	PetType             string    `json:"pet_type" gorm:"not null"`
	ShelterAnimal       string    `json:"shelter_animal"`
	IdealPetDescription string    `json:"ideal_pet_description"`
	HousingSituation    string    `json:"housing_situation" gorm:"not null"`
	PetsAtHome          string    `json:"pets_at_home"`
	Allergies           string    `json:"allergies"`
	FamilySupport       string    `json:"family_support"`
	PastPets            string    `json:"past_pets"`
	InterviewSetting    string    `json:"interview_setting"`
	ValidID             string    `json:"valid_id" gorm:"not null"`     // For adopter
	AltValidID          string    `json:"alt_valid_id" gorm:"not null"` // For alternate contact
	Status              string    `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt           time.Time 

	// Shelter ShelterInfo `json:"shelter"`
	Adopter AdopterInfo `json:"adopter"`
	Pet     PetInfo     `json:"pet"`
>>>>>>> b4994a12ff9dad102a75a877e0b854cd06d1356e
}

// TableName overrides default table name
func (AdoptionSubmission) TableName() string {
	return "adoption_submissions"
}
