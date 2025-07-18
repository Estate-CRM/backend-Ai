package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connectdb() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DB_SSL,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	// Auto migrate models
	/* err = db.AutoMigrate(
		&model.Contact{},
	) */
	DB = db
	log.Println("Database connected successfully")
}
