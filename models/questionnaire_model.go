package models

import (
	"time"
)

// Questionnaires represents the adoption questionnaire table
type Questionnaires struct {
	QuestionID               uint   `json:"question_id" gorm:"primaryKey;autoIncrement"`
	ApplicationID            uint   `json:"application_id" gorm:"not null;constraint:OnDelete:CASCADE;"` // Foreign key to applications
	PetType                  string `json:"pet_type" gorm:"not null"`                                    // Cat, Dog, Both, Not decided
	SpecificShelterAnimal    bool   `json:"specific_shelter_animal" gorm:"default:false"`                // Applying for a specific animal?
	IdealPetDescription      string `json:"ideal_pet_description,omitempty"`
	BuildingType             string `json:"building_type,omitempty"`
	Rent                     bool   `json:"rent" gorm:"default:false"`
	PetMovePlan              string `json:"pet_move_plan,omitempty"`
	HouseholdComposition     string `json:"household_composition,omitempty"`
	AllergiesToAnimals       bool   `json:"allergies_to_animals" gorm:"default:false"`
	CareResponsibility       string `json:"care_responsibility,omitempty"`
	FinancialResponsibility  string `json:"financial_responsibility,omitempty"`
	VacationCarePlan         string `json:"vacation_care_plan,omitempty"`
	AloneTime                int    `json:"alone_time,omitempty"`
	IntroductionPlan         string `json:"introduction_plan,omitempty"`
	FamilySupport            bool   `json:"family_support" gorm:"default:false"`
	FamilySupportExplanation string `json:"family_support_explanation,omitempty"`
	OtherPets                bool   `json:"other_pets" gorm:"default:false"`
	PastPets                 bool   `json:"past_pets" gorm:"default:false"`

	// Home Photos (stored as JSONB for multiple images)
	HomePhotos []string `json:"home_photos,omitempty" gorm:"type:jsonb"`

	// Valid ID (stores front and back images as JSONB)
	ValidID map[string]string `json:"valid_id,omitempty" gorm:"type:jsonb"`

	// Interview details
	ZoomInterviewDate time.Time `json:"zoom_interview_date,omitempty"`
	ZoomInterviewTime time.Time `json:"zoom_interview_time,omitempty"`

	ShelterVisit bool      `json:"shelter_visit" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName sets the database table name
func (Questionnaires) TableName() string {
	return "questionnaires"
}
