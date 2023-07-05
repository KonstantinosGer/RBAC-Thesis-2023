package database

import (
	"gorm.io/gorm"
)

func CloseDBConnectionGorm(db *gorm.DB) {
	dbConn, _ := db.DB()
	dbConn.Close()
}
