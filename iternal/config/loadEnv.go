package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Loaded .env file:", os.Getenv("DB_URL"))
}
