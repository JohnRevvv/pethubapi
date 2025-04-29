package models

type AdminAccount struct {
	AdminID  uint   `json:"admin_id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// TableName overrides default table name
func (AdminAccount) TableName() string {
	return "adminaccount"
}
