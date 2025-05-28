package models

import "time"

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AdopterID uint      `json:"adopter_id"`
	PetID     uint      `json:"pet_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Category  string    `json:"category"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}
