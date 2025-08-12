package cmd

import (
	"l0/iternal/config"
	"l0/models"
	"log"
)

func Migration() {
	config.LoadEnvs()
	config.ConnectDB()
	err := config.DB.AutoMigrate(&models.Order{}, &models.Payment{}, &models.Item{}, &models.Delivery{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration completed successfully")
}
