package models

import "time"

// AdopterAccount model (linked to existing "adopteraccount" table)
type AdopterAccount struct {
	AdopterID uint   `gorm:"primaryKey" json:"adopter_id"`
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `json:"password"`
	Status    string `gorm:"default:'active'" json:"status"` // Add this line
	CreatedAt time.Time
	//Info      AdopterInfo `gorm:"foreignKey:AdopterID;constraint:OnDelete:CASCADE" json:"info"`
}

// TableName overrides default table name
func (AdopterAccount) TableName() string {
	return "adopteraccount"
}

// AdopterInfo model (linked to existing "adopterinfo" table)
type AdopterInfo struct {
	AdopterID     uint   `gorm:"primaryKey;autoIncrement:false" json:"adopter_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Age           int    `json:"age"`
	Sex           string `json:"sex"`
	Address       string `json:"address"`
	ContactNumber string `json:"contact_number"`
	Email         string `gorm:"unique" json:"email"`
	Occupation    string `json:"occupation"`
	CivilStatus   string `json:"civil_status"`
	SocialMedia   string `json:"social_media"`
}

// TableName overrides default table name
func (AdopterInfo) TableName() string {
	return "adopterinfo"
}

type AdopterMedia struct {
	AdopterID      uint   `gorm:"primaryKey;autoIncrement:false" json:"adopter_id"`
	AdopterProfile string `json:"adopter_profile"`
}

func (AdoptedPet) TableName() string {
	return "adopterpets"
}

type AdoptedPet struct {
	AdoptedID uint `gorm:"column:adopted_id;primaryKey;autoIncrement" json:"adopted_id"`

	AdopterID uint `gorm:"column:adopter_id" json:"adopter_id"`

	PetID uint `gorm:"column:pet_id" json:"pet_id"`
}
