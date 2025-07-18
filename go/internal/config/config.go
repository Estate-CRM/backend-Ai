package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	DB_SSL        string
	Port          string
	Access_token  string
	Refresh_token string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := Config{
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBName:        os.Getenv("DB_NAME"),
		DB_SSL:        os.Getenv("DB_SSL"),
		Port:          os.Getenv("PORT"),
		Access_token:  os.Getenv("ACCESS_TOKEN_SECRET"),
		Refresh_token: os.Getenv("REFRESH_TOKEN_SECRET"),
	}
	return cfg, nil
}
