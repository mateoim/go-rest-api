package models

import (
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

	err = db.AutoMigrate(&Organization{})
	if err != nil {
		panic("Error setting up Organizations table")
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Error setting up Users table")
	}

	DB = db
}
