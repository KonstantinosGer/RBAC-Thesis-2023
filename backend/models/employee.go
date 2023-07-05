package models

type Employee struct {
	Id       int    `json:"id" db:"id" uri:"id" gorm:"primaryKey"`
	FullName string `json:"full_name" db:"full_name"`
	// Association (with user)
	// One employee can be associated with more than one (firebase) users
	Users []User //`gorm:"polymorphic:Owner;"`
}

// TableName Change table's name manually, instead of letting gorm do it automatically
//func (Employee) TableName() string {
//	return "employees"
//}
