package db

import (
	"fmt"
	"itkdemo/internal/transport/db/model"
	"itkdemo/pkg/config"
	logger "itkdemo/pkg/log"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres" // psx driver
	"gorm.io/gorm"
)

var db *gorm.DB

func New() *gorm.DB {
	if db != nil {
		return db
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("failed to connect database:", err)
	}

	err = db.AutoMigrate(&model.Wallet{})
	if err != nil {
		logger.Log.Fatal("failed to migrate schema:", err)
	}

	return db
}
