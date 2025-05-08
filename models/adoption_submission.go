package models

import (
	"time"
)

type AdoptionSubmission struct {
	// Application Info
	ApplicationID       uint               `json:"application_id" gorm:"primaryKey;autoIncrement"`
	ShelterID           uint               `json:"shelter_id"`
	PetID               uint               `json:"pet_id"`
	AdopterID           uint               `json:"adopter_id"`
	AltFName            string             `json:"alt_f_name" gorm:"not null"`
	AltLName            string             `json:"alt_l_name" gorm:"not null"`
	Relationship        string             `json:"relationship" gorm:"not null"`
	AltContactNumber    string             `json:"alt_contact_number" gorm:"not null"`
	AltEmail            string             `json:"alt_email" gorm:"not null"`
	ReasonForAdoption   string             `json:"reason_for_adoption" gorm:"not null"`
	IdealPetDescription string             `json:"ideal_pet_description"`
	HousingSituation    string             `json:"housing_situation" gorm:"not null"`
	PetsAtHome          string             `json:"pets_at_home"`
	Allergies           string             `json:"allergies"`
	FamilySupport       string             `json:"family_support"`
	PastPets            string             `json:"past_pets"`
	InterviewSetting    string             `json:"interview_setting"`
	ValidID             string             `json:"valid_id" gorm:"not null"`     // For adopter
	AltValidID          string             `json:"alt_valid_id" gorm:"not null"` // For alternate contact
	Status              string             `json:"status" gorm:"type:varchar(20);default:'pending'"`
	ApplicationPhotos   []ApplicationPhoto `json:"application_photos" gorm:"-"`
	CreatedAt           time.Time

	// Shelter ShelterInfo `json:"shelter"`
	Adopter AdopterInfo `json:"adopter"`
	Pet     PetInfo     `json:"pet"`
}

// TableName overrides default table name
func (AdoptionSubmission) TableName() string {
	return "adoption_submissions"
}

type ApplicationPhoto struct {
	ID            uint   `gorm:"primaryKey"`
	ApplicationID uint   `gorm:"not null"`
	PhotoType     string `gorm:"not null"` // "Home Photo or PDF"
	Base64Data    string `gorm:"type:text"`
	UploadedAt    time.Time
}

// For structured Valid ID info
type ValidIDPhotos struct {
	ValidID        uint   `gorm:"primaryKey;autoIncWrement"`
	AdopterIDType  string `json:"adopter_id_type"`
	AdopterValidID string `json:"adopter_valid_id"`
	AltIDType      string `json:"alt_id_type"`
	AltValidID     string `json:"alt_valid_id"`
}

// TableName overrides (both map to the same table)
func (ApplicationPhoto) TableName() string { return "application_photos" }
func (ValidIDPhotos) TableName() string    { return "application_photos" }
