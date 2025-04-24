package models

import (
	"time"
)

// Questionnaires represents the adoption questionnaire table
type Questionnaires struct {
	QuestionID                uint   `json:"question_id" gorm:"primaryKey;autoIncrement"`
	ApplicationID             uint   `json:"application_id" gorm:"column:application_id"` // Foreign key to applications
	PetType                   string `json:"pet_type" gorm:"not null"`                    // Cat, Dog, Both, Not decided
	SpecificShelterAnimal     string `json:"specific_shelter_animal" gorm:"default:'no'"` // Change to string ("yes"/"no")
	IdealPetDescription       string `json:"ideal_pet_description,omitempty" gorm:"column:ideal_pet_description"`
	BuildingType              string `json:"building_type,omitempty"`
	Rent                      string `json:"rent" gorm:"default:'no'"` // Change to string ("yes"/"no")
	PetMovePlan               string `json:"pet_move_plan,omitempty"`
	HouseholdComposition      string `json:"household_composition,omitempty"`
	AllergiesToAnimals        string `json:"allergies_to_animals" gorm:"default:'no'"` // Change to string ("yes"/"no")
	CareResponsibility        string `json:"care_responsibility,omitempty"`
	FinancialResponsibility   string `json:"financial_responsibility,omitempty"`
	VacationCarePlan          string `json:"vacation_care_plan,omitempty"`
	AloneTime                 string `json:"alone_time,omitempty"`
	IntroductionPlan          string `json:"introduction_plan,omitempty"`
	FamilySupport             string `json:"family_support" gorm:"default:'no'"` // Change to string ("yes"/"no")
	FamilySupportExplanation  string `json:"family_support_explanation,omitempty"`
	OtherPets                 string `json:"other_pets" gorm:"default:'no'"` // Change to string ("yes"/"no")
	PastPets                  string `json:"past_pets" gorm:"default:'no'"`  // Change to string ("yes"/"no")
	HomePhotos                string // store paths as comma-separated string
	ValidID                   string
	IDType                    string    `json:"id_type,omitempty"`
	PreferredInterviewSetting string    `json:"preferred_interview_setting"` // NEW FIELD
	CreatedAt                 time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt                 time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName sets the database table name
func (Questionnaires) TableName() string {
	return "questionnaires"
}
