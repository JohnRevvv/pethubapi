package models

import "time"

// ShelterAccount model (linked to existing "shelteraccount" table)
type ShelterAccount struct {
	ShelterID uint   `gorm:"primaryKey" json:"shelter_id"`
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `json:"password"`
	Status    string `gorm:"default:'active'" json:"status"` // Add this line
	RegStatus string `json:"reg_status"`
	CreatedAt time.Time

	ShelterInfo ShelterInfo `gorm:"foreignKey:ShelterID" json:"shelterinfo"`
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

	ShelterMedia ShelterMedia `gorm:"foreignKey:ShelterID;references:ShelterID" json:"sheltermedia"`
}

// TableName overrides default table name
func (ShelterInfo) TableName() string {
	return "shelterinfo"
}

type ShelterMedia struct {
	ShelterID      uint   `gorm:"primaryKey;autoIncrement:false" json:"shelter_id"`
	ShelterProfile string `json:"shelter_profile"`
	ShelterCover   string `json:"shelter_cover"`
}

func (ShelterMedia) TableName() string {
	return "sheltermedia"
}

type ShelterDonations struct {
	DonationID    uint   `gorm:"primaryKey;autoIncrement:true" json:"donation_id"`
	ShelterID     uint   `json:"shelter_id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	QRImage       string `json:"qr_image"`
	CreatedAt     time.Time
}

func (ShelterDonations) TableName() string {
	return "shelterdonations"
}
