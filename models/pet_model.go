package models

import "time"

type PetInfo struct {
	PetID           uint      `gorm:"primaryKey" json:"pet_id"`
	ShelterID       uint      `gorm:"autoIncrement:false" json:"shelter_id"`
	PetType         string    `json:"pet_type"`
	PetName         string    `json:"pet_name"`
	PetAge          int       `json:"pet_age"`
	AgeType         string    `json:"age_type"`
	PetSex          string    `json:"pet_sex"`
	PetDescriptions string    `json:"pet_descriptions"`
	Status          string    `gorm:"default:'available'" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	PetSize         string    `json:"pet_size"`
	PriorityStatus  bool      `json:"priority_status"`
	PetMedia        PetMedia  `gorm:"foreignKey:PetID;references:PetID" json:"petmedia"`
}

func (PetInfo) TableName() string {
	return "petinfo"
}

type PetMedia struct {
	PetID     uint   `gorm:"not null" json:"pet_id"`
	PetImage1 string `json:"pet_image1"` // Base64-encoded image
}

func (PetMedia) TableName() string {
	return "petmedia"
}
