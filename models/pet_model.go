package models

import "time"

type PetInfo struct {
	PetID           uint   `gorm:"primaryKey" json:"pet_id"`
	ShelterID       int    `gorm:"autoIncrement:false" json:"shelter_id"`
	PetName         string `json:"pet_name"`
	PetAge          int    `json:"pet_age"`
	PetSex          string `json:"pet_sex"`
	PetDescriptions string `json:"pet_descriptions"`
	Status          string `gorm:"default:'available'" json:"status"`
	CreatedAt       time.Time
}

// TableName overrides default table name
func (PetInfo) TableName() string {
	return "petinfo"
}

type PetMedia struct {
	PetID      uint   `gorm:"primaryKey" json:"pet_id"`
	PetProfile []byte `json:"pet_profile"` // Binary data for profile image
	PetImage1  []byte `json:"pet_image1"`  // Binary data for image 1
	PetImage2  []byte `json:"pet_image2"`  // Binary data for image 2
	PetImage3  []byte `json:"pet_image3"`  // Binary data for image 3
}

// TableName overrides default table name
func (PetMedia) TableName() string {
	return "petmedia"
}
