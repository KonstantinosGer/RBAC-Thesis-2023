package models

type Policy struct {
	Id         int    `json:"id" db:"id"`
	PolicyType string `json:"ptype" db:"ptype"`
	Role       string `json:"role" db:"v0"`
	Data       string `json:"data" db:"v1"`
	Privilege  string `json:"privilege" db:"v2"`

	//Id         int    `json:"id" db:"id"`
	//PolicyType string `json:"ptype" db:"ptype"`
	//V0         string `json:"v0" db:"v0"`
	//V1         string `json:"v1" db:"v1"`
	//V2         string `json:"v2" db:"v2"`
}
