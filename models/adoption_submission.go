package models

import (
	"time"
)

type AdoptionSubmission struct {
	// Application Info
	ApplicationID       uint      `json:"application_id" gorm:"primaryKey;autoIncrement"`
	ShelterID           uint      `json:"shelter_id"`
	PetID               uint      `json:"pet_id"`
	AdopterID           uint      `json:"adopter_id"`
	AltFName            string    `json:"alt_f_name" gorm:"not null"`
	AltLName            string    `json:"alt_l_name" gorm:"not null"`
	Relationship        string    `json:"relationship" gorm:"not null"`
	AltContactNumber    string    `json:"alt_contact_number" gorm:"not null"`
	AltEmail            string    `json:"alt_email" gorm:"not null"`
	ReasonForAdoption   string    `json:"reason_for_adoption" gorm:"not null"`
	IdealPetDescription string    `json:"ideal_pet_description"`
	HousingSituation    string    `json:"housing_situation" gorm:"not null"`
	PetsAtHome          string    `json:"pets_at_home"`
	Allergies           string    `json:"allergies"`
	FamilySupport       string    `json:"family_support"`
	PastPets            string    `json:"past_pets"`
	InterviewSetting    string    `json:"interview_setting"`
	ImageID             uint      `json:"image_id"`
	Status              string    `json:"status" gorm:"type:varchar(20);default:'pending'"`
	ReasonForRejection  string    `json:"reason_for_rejection"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Shelter           ShelterInfo       `json:"shelter"`
	Adopter           AdopterInfo       `json:"adopter"`
	Pet               PetInfo           `json:"pet"`
	ScheduleInterview ScheduleInterview `gorm:"foreignKey:ApplicationID;references:ApplicationID"  json:"scheduleinterview"`
}

// TableName overrides default table name
func (AdoptionSubmission) TableName() string {
	return "adoption_submissions"
}

// For structured Valid ID info
type ApplicationPhotos struct {
	ImageID        uint   `gorm:"primaryKey"`
	AdopterIDType  string `json:"adopter_id_type"`
	AdopterValidID string `json:"adopter_valid_id"`
	AltIDType      string `json:"alt_id_type"`
	AltValidID     string `json:"alt_valid_id"`
	HomeImage1     string `json:"home_image1"`
	HomeImage2     string `json:"home_image2"`
	HomeImage3     string `json:"home_image3"`
	HomeImage4     string `json:"home_image4"`
	HomeImage5     string `json:"home_image5"`
	HomeImage6     string `json:"home_image6"`
	HomeImage7     string `json:"home_image7"`
	HomeImage8     string `json:"home_image8"`
}

// TableName overrides (both map to the same table)
func (ApplicationPhotos) TableName() string {
	return "application_photos"
}

type ScheduleInterview struct {
	InterviewID     uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ApplicationID   uint      `json:"application_id"`
	ShelterID       uint      `json:"shelter_id"`
	AdopterID       uint      `json:"adopter_id"`
	InterviewDate   time.Time `json:"interview_date"`
	InterviewTime   string    `json:"interview_time"`
	InterviewNotes  string    `json:"interview_notes"`
	InterviewStatus string    `json:"interview_status" gorm:"type:varchar(20);default:'scheduled'"`
	CreatedAt       time.Time `json:"created_at"`
}

func (ScheduleInterview) TableName() string {
	return "schedule_interview"
}
