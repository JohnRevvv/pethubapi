package models

type AdminAccount struct {
<<<<<<< HEAD
	AdminID  uint   `json:"admin_id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
=======
	AdminID   uint      `json:"admin_id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
>>>>>>> b4994a12ff9dad102a75a877e0b854cd06d1356e
}

// TableName overrides default table name
func (AdminAccount) TableName() string {
	return "adminaccount"
}
