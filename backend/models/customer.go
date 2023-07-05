package models

type Customer struct {
	Id       int    `json:"id" db:"id" uri:"id" gorm:"primaryKey"`
	FullName string `json:"full_name" db:"full_name"`
	// Association (with user)
	// One customer can be associated with more than one (firebase) users
	// (omitempty means if nothing is given from frontend for this field, omit that field and do not return it from backend)
	//When omitempty is used, it tells the encoder to omit a field if its value is the zero value for its type (0 for integers, "" for strings, false for booleans, nil for pointers, slices, maps, channels, and interfaces, etc.) during encoding. During decoding, if the JSON value for the field is missing or null, the decoder will set the field to its zero value.
	//The purpose of omitempty is to reduce the size of JSON objects when the zero value is considered meaningless or redundant. This can be particularly useful when dealing with large JSON objects or when sending data over a slow network connection.
	Users []*User `json:"users,omitempty" gorm:"many2many:customer_user;"`
}
