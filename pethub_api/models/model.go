package models

import "gorm.io/gorm"

type (
	User struct {
		gorm.Model
		ID  int    `json:"primaryKey:id"`
		Name string `json:"name"`
		Section string `json:"section"`
	}
)

type (
	Game struct {
		ID        int    `json:"ID"`
		Name      string `json:"name"`
		Developer string `json:"developer"`
	}
)
