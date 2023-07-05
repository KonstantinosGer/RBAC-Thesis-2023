package models

type FrontendUser struct {
	Id       string `json:"id" db:"id"`
	FullName string `json:"full_name" db:"full_name"`
	Email    string `json:"email" db:"email"`
	Role     string `json:"role" db:"v1"`
}
