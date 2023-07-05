package models

type User struct {
	Id                 string `json:"id" db:"id" uri:"id" gorm:"primaryKey"`
	Email              string `json:"email" db:"email"`
	CreationTimestamp  int    `json:"creation_timestamp" db:"creation_timestamp"`
	LastLoginTimestamp int    `json:"last_login_timestamp" db:"last_login_timestamp"`
	// Association
	// User can be associated with one employee
	EmployeeID *int `json:"employee_id"` //* means it can be null

	// One (firebase) user can be associated with more than one customer
	//Users    []User `gorm:"polymorphic:Owner;"`
	Customers []*Customer `gorm:"many2many:customer_user;"`
}
