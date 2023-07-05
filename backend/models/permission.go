package models

type Permission struct {
	Id          int    `json:"id" db:"id" gorm:"primaryKey"`
	Action      string `json:"action" db:"action"`
	Resource    string `json:"resource" db:"resource"`
	Description string `json:"description" db:"description"`
	Category    string `json:"category" db:"category"`
	CategoryNo  int    `json:"category_no" db:"category_no"`
}
