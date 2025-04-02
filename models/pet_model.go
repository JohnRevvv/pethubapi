package models

import "time"

type PetInfo struct {
	PetID           uint   `gorm:"primaryKey" json:"pet_id"`
	ShelterID       uint   `gorm:"autoIncrement:false" json:"shelter_id"`
	PetName         string `json:"pet_name"`
	PetAge          int    `json:"pet_age"`
	AgeType         string `json:"age_type"`
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
	PetID uint `gorm:"not null;constraint:OnDelete:CASCADE;" json:"pet_id"`
	PetImage1 string ` json:"pet_image1"`            // Base64-encoded image
	PetImage2 string ` json:"pet_image2"`
	PetImage3 string `json:"pet_image3"`
	PetImage4 string `json:"pet_image4"`
}

// TableName overrides default table name
func (PetMedia) TableName() string {
	return "petmedia"
}
