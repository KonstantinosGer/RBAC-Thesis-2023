package models

type CasbinRule struct {
	Id    int    `json:"id" db:"id" gorm:"primaryKey"`
	Ptype string `json:"ptype" db:"ptype"`
	V0    string `json:"v0" db:"v0"`
	V1    string `json:"v1" db:"v1"`
	V2    string `json:"v2" db:"v2"`
	V3    string `json:"v3" db:"v3"`
	V4    string `json:"v4" db:"v4"`
	V5    string `json:"v5" db:"v5"`
}

//TableName Change table's name manually, instead of letting gorm do it automatically
func (CasbinRule) TableName() string {
	return "casbin_rule"
}
