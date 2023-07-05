package config

import (
	"log"
	//"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

//LOGGER
func InitLogging() {
	if os.Getenv("APP_ENV") == "prod" {
		if _, err := os.Stat("./logs/"); os.IsNotExist(err) {
			os.MkdirAll("./logs", 0700)
		}
		dt := time.Now()
		file, err := os.OpenFile("./logs/logs."+dt.Format("2006-01-02")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
		log.Println("Set log to file")
	} else {
		log.SetOutput(os.Stdout)
		log.Println("Set log to stdout")
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Port() (string, error) {
	var envs map[string]string
	envs, err := godotenv.Read("./info/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := envs["PORT"]
	return ":" + port, nil
}

func DB() string {
	dbuser := os.Getenv("DBUSER")
	dbpass := os.Getenv("DBPASS")
	dbhost := os.Getenv("DBHOST")
	dbport := os.Getenv("DBPORT")
	dbname := os.Getenv("DBNAME")

	return dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname
}

func ENV(key string) string {
	return os.Getenv(key)
}
