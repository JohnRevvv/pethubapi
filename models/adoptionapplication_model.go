package models

import (
	"time"
)

// AdoptionApplication model (linked to existing "applicationform" table)
type AdoptionApplication struct {
	ApplicationID    uint      `json:"application_id" gorm:"primaryKey;autoIncrement;column:application_id"`
	PetID            uint      `json:"pet_id" gorm:"column:pet_id;not null"`
	AdopterID        uint      `json:"adopter_id" gorm:"column:adopter_id;not null"`
	AltFName         string    `json:"alt_f_name" gorm:"column:alt_f_name;not null"`
	AltLName         string    `json:"alt_l_name" gorm:"column:alt_l_name;not null"`
	Relationship     string    `json:"relationship" gorm:"column:relationship;not null"`
	AltContactNumber string    `json:"alt_contact_number" gorm:"column:alt_contact_number;not null"`
	AltEmail         string    `json:"alt_email" gorm:"column:alt_email;not null"`
	HouseFile        string    `json:"housefile" gorm:"column:housefile;not null"` // Path or URL to house photos
	ValidID          string    `json:"alt_valid_id" gorm:"column:alt_valid_id;not null"`
	PreferredDate    string    `json:"preferred_date" gorm:"column:preferred_date;not null"` // Changed to time.Time
	PreferredTime    string    `json:"preferred_time" gorm:"column:preferred_time;not null"` // Changed to time.Time
	Status           string    `json:"status" gorm:"column:status;type:varchar(20);default:'Pending'"`
	CreatedAt        time.Time `json:"created_at"` // Explicit CreatedAt field
	UpdatedAt        time.Time `json:"updated_at"` // Explicit UpdatedAt field
	DeletedAt        time.Time `json:"deleted_at"` // Explicit DeletedAt field
	Questionnaires  Questionnaires `gorm:"foreignKey:ApplicationID;references:ApplicationID" json:"questionnaires"`
}

// TableName overrides default table name
func (AdoptionApplication) TableName() string {
	return "applicationform"
}

// Questionnaires represents the adoption questionnaire table
type Questionnaires struct {
	QuestionID                uint   `json:"question_id" gorm:"primaryKey;autoIncrement"`
	ApplicationID             uint   `json:"application_id" gorm:"column:application_id"` // Foreign key to applications
	PetType                   string `json:"pet_type" gorm:"not null"`             // Cat, Dog, Both, Not decided
	SpecificShelterAnimal     string `json:"specific_shelter_animal" gorm:"default:'no'" `// Change to string ("yes"/"no")
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
	FamilySupport             string `json:"family_support" gorm:"default:'no'"`// Change to string ("yes"/"no")
	FamilySupportExplanation  string `json:"family_support_explanation,omitempty"`
	OtherPets                 string `json:"other_pets" gorm:"default:'no'" `// Change to string ("yes"/"no")
	PastPets                  string `json:"past_pets" gorm:"default:'no'"`  // Change to string ("yes"/"no")
	HomePhotos                string // store paths as comma-separated string
	ValidID                   string `json:"valid_id"` // Change to string ("yes"/"no")
	IDType                    string    `json:"id_type,omitempty"`
	PreferredInterviewSetting string    `json:"preferred_interview_setting" `// NEW FIELD
	CreatedAt                 time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt                 time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName sets the database table name
func (Questionnaires) TableName() string {
	return "questionnaires"
}
