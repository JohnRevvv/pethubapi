package models

import "time"

// ShelterAccount model (linked to existing "shelteraccount" table)
type ShelterAccount struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `json:"password"`
	CreatedAt time.Time
	//Info      ShelterInfo `gorm:"foreignKey:ShelterID;constraint:OnDelete:CASCADE" json:"info"`
}




// ShelterInfo model (linked to existing "shelterinfo" table)
type ShelterInfo struct {
	ShelterID          uint   `gorm:"primaryKey;autoIncrement:false" json:"shelter_id"`
	ShelterName        string `json:"shelter_name"`
	ShelterAddress     string `json:"shelter_address"`
	ShelterContact     int    `json:"shelter_contact"`
	ShelterOwner       string `json:"shelter_owner"`
	ShelterDescription string `json:"shelter_description"`
}

// TableName overrides default table name
func (ShelterInfo) TableName() string {
	return "shelterinfo"
}
