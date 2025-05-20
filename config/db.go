package config

import (
	"effective-mobile-test/logger"
	"effective-mobile-test/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitDB() {
	log := logger.L()
	dsn := os.Getenv("DATABASE_URL")
	log.WithField("dsn", dsn).Info("Connecting to the database")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("Database connection failed")
	}
	if err := db.AutoMigrate(&models.Person{}); err != nil {
		log.WithError(err).Error("AutoMigrate failed")
	} else {
		log.Info("Database migrated successfully")
	}

	DB = db
}
