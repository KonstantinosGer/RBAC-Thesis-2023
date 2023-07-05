package models

type Role struct {
	Role        string `json:"role" db:"role" gorm:"primaryKey" form:"role"`
	Description string `json:"description" db:"description"`
}
