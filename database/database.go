package database

import (
	"fmt"
	"log"
	"os"

	"example.com/go-fiber-api/cmd/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/Regina",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info), //!NOTE: Enable this to log performance, you have to import logger for this as well.
	})

	if err != nil {
		log.Fatal("Failed to connect to db. \n", err)
		os.Exit(2)
	}

	log.Println("Connected to db")
	//db.Logger = logger.Default.LogMode(logger.Info)

	// Enable UUID extension
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatal("Failed to create extension \"uuid-ossp\". \n", err)
		os.Exit(2)
	}

	//building models
	log.Println("Running Migrations")
	db.AutoMigrate(&models.Item{}, &models.User{})

	DB = Dbinstance{
		Db: db,
	}
}
