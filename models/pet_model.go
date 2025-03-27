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
