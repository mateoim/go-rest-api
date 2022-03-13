package config

import (
	"go-rest-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	dsn := "host=ec2-54-194-147-61.eu-west-1.compute.amazonaws.com " +
		"user=ftmicgdagkwunl " +
		"password= " +
		"dbname=dee746dp5f4pme " +
		"port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = db.AutoMigrate(&models.Organization{})
	if err != nil {
		panic("Error setting up Organizations table")
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("Error setting up Users table")
	}

	err = db.AutoMigrate(&models.Event{})
	if err != nil {
		panic("Error setting up Events table")
	}

	err = db.AutoMigrate(&models.Meeting{})
	if err != nil {
		panic("Error setting up Meetings table")
	}

	err = db.AutoMigrate(&models.Invitation{})
	if err != nil {
		panic("Error setting up Invitations table")
	}

	DB = db
}
