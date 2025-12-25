package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string = "5432"
	Port       int
	Debug      bool
)

func Init() {
	// Load .env file
	err := godotenv.Load(".env", "config.env")
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Read environment variables
	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
	DBPort = os.Getenv("DB_PORT")
	Port, _ = strconv.Atoi(os.Getenv("PORT"))
	Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
}
