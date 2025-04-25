package models

import (
	"time"
)

type AdoptionSubmission struct {
	// Application Info
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt           time.Time `json:"deleted_at"`
	ApplicationID       uint      `json:"application_id" gorm:"primaryKey;autoIncrement"`
	ShelterID           uint      `gorm:"not null" json:"shelter_id"`
	PetID               uint      `json:"pet_id" gorm:"not null"`
	AdopterID           uint      `json:"adopter_id" gorm:"not null"`
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
	Status              string    `json:"status" gorm:"type:varchar(20);default:'Pending'"`

	// Shelter ShelterInfo `gorm:"foreignKey:ShelterID;references:ShelterID" json:"shelter"`
	// Adopter AdopterInfo `gorm:"foreignKey:AdopterID;references:AdopterID" json:"adopter"`
	// Pet     PetInfo     `gorm:"foreignKey:PetID;references:PetID" json:"pet"`
}

// TableName overrides default table name
func (AdoptionSubmission) TableName() string {
	return "adoption_submissions"
}
