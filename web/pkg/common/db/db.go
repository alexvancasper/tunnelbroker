package db

import (
	"log"

	"github.com/alexvancasper/TunnelBroker/web/pkg/models"
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
	db.AutoMigrate(&models.User{}, &models.Tunnel{})

	return db
}
