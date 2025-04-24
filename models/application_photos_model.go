package models

import "time"

type ApplicationPhoto struct {
	ID            uint   `gorm:"primaryKey"`
	ApplicationID uint   `gorm:"not null"`
	PhotoType     string `gorm:"not null"` // ✅ Ensure this exists and is exported
	Base64Data    string `gorm:"type:text"`
	UploadedAt    time.Time
}

// TableName overrides default table name
func (ApplicationPhoto) TableName() string {
	return "application_photos"
}
