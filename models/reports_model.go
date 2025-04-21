package models

import "time"

type Report struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	ShelterID   uint      `json:"shelter_id"`
	AdopterID   uint      ` json:"adopter_id"`
	Reason      string    `gorm:"type:text;column:reason" json:"reason"`
	Description string    `gorm:"type:text;column:description" json:"description"`
	Status      string    `gorm:"type:text;column:status;default:'pending'" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// Explicitly sets the table name to 'submittedreports'
func (Report) TableName() string {
	return "submittedreports"
}
