package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	BasePath   string
}

var Conf Config

func InitializationConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file. Error - %s", err.Error())
		return
	}

	Conf = Config{
		ServerPort: os.Getenv("TODO_PORT"),
		BasePath:   os.Getenv("TODO_DBFILE"),
	}
}
