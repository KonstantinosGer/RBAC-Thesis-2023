package migrate

import (
	"backend/database"
	"backend/models"
	"log"
)

func MigrateDBGorm() {

	db, err := database.ConnectToRBACGorm()
	if err != nil {
		log.Println(err)
	}
	defer database.CloseDBConnectionGorm(db)

	err = db.Migrator().AutoMigrate(
		&models.Customer{},
		&models.Employee{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)
	if err != nil {
		log.Println(err)
	}

}
