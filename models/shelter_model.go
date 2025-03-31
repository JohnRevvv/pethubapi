package models

import "time"

// ShelterAccount model (linked to existing "shelteraccount" table)
type ShelterAccount struct {
	ShelterID uint   `gorm:"primaryKey" json:"shelter_id"`
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `json:"password"`
	CreatedAt time.Time
	//Info      ShelterInfo `gorm:"foreignKey:ShelterID;constraint:OnDelete:CASCADE" json:"info"`
}

// TableName overrides default table name
func (ShelterAccount) TableName() string {
	return "shelteraccount"
}

// ShelterInfo model (linked to existing "shelterinfo" table)
type ShelterInfo struct {
	ShelterID          uint   `gorm:"primaryKey;autoIncrement:false" json:"shelter_id"`
	ShelterName        string `json:"shelter_name"`
	ShelterAddress     string `json:"shelter_address"`
	ShelterLandmark    string `json:"shelter_landmark"`
	ShelterContact     string `json:"shelter_contact"`
	ShelterEmail       string `json:"shelter_email"`
	ShelterOwner       string `json:"shelter_owner"`
	ShelterDescription string `json:"shelter_description"`
	ShelterSocial      string `json:"shelter_social"`
}

// TableName overrides default table name
func (ShelterInfo) TableName() string {
	return "shelterinfo"
}

type ShelterMedia struct {
	ShelterID      uint   `gorm:"primaryKey;autoIncrement:false" json:"shelter_id"`
	ShelterProfile string `json:"shelter_profile"`
	ShelterCover   string `json: "shelter_cover"`
}

func (ShelterMedia) TableName() string {
	return "sheltermedia"
}
