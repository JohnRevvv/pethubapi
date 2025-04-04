package models

import "time"

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
	ValidID          string    `json:"valid_id" gorm:"column:valid_id;not null"`
	PreferredDate    time.Time `json:"preferred_date" gorm:"column:preferred_date;not null"` // Changed to time.Time
	PreferredTime    time.Time `json:"preferred_time" gorm:"column:preferred_time;not null"` // Changed to time.Time
	Status           string    `json:"status" gorm:"column:status;type:varchar(20);default:'Pending'"`
	CreatedAt        time.Time `json:"created_at"` // Explicit CreatedAt field
	UpdatedAt        time.Time `json:"updated_at"` // Explicit UpdatedAt field
	DeletedAt        time.Time `json:"deleted_at"` // Explicit DeletedAt field
}

// TableName overrides default table name
func (AdoptionApplication) TableName() string {
	return "applicationform"
}
