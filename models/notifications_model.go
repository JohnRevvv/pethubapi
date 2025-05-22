package models

import "time"

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AdopterID uint      `json:"adopter_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`   // e.g., "adoption_status"
	Status    string    `json:"status"` // e.g., "pending", "approved"
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}
