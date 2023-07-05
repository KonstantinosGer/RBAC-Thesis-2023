package main

import (
	"backend/config"
	"backend/database/migrate"
	"backend/models"
	"backend/routes"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	//     LOGGING
	config.InitLogging()

	//ENV
	err := godotenv.Load("./info/.env")
	if err != nil {
		log.Println("Error loading env in main")
	}

	// migrate db model to updated (only when gorm is used)
	migrate.MigrateDBGorm()

	db, _ := models.DBConnection()
	routes.SetupRoutes(db)
}
