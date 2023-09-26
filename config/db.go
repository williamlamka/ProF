package config

import (
	"log"
	"new_project/model"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func DBInit() {
	var err error
	db, err = gorm.Open(postgres.Open(os.Getenv("POSTGRES")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Bill{})
}

func DB() *gorm.DB {
	return db
}