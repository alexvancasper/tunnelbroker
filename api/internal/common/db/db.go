package db

import (
	"log"

	"github.com/alexvancasper/TunnelBroker/web/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	DB = db
	err = db.AutoMigrate(&models.User{}, &models.Tunnel{})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
