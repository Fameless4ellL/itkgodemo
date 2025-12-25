package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DBHost     string = "db"
	DBUser     string = "postgres"
	DBPassword string = "postgres"
	DBName     string = "postgres"
	DBPort     string = "5432"
	Port       int    = 8080
	Debug      bool   = false
)

func Init() {
	// Load config.env file
	err := godotenv.Load("config.env")
	if err != nil {
		log.Println("No config.env file found, using system environment variables")
	}

	// Read environment variables
	if key, ok := os.LookupEnv("DB_HOST"); ok {
		DBHost = key
	}
	if key, ok := os.LookupEnv("DB_USER"); ok {
		DBUser = key
	}
	if key, ok := os.LookupEnv("DB_PASSWORD"); ok {
		DBPassword = key
	}
	if key, ok := os.LookupEnv("DB_NAME"); ok {
		DBName = key
	}
	if key, ok := os.LookupEnv("DB_PORT"); ok {
		DBPort = key
	}
	if key, ok := os.LookupEnv("PORT"); ok {
		Port, _ = strconv.Atoi(key)
	}
	if key, ok := os.LookupEnv("DEBUG"); ok {
		Debug, _ = strconv.ParseBool(key)
	}
}
