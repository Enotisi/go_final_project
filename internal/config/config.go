package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	BasePath   string
}

var Conf Config

func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	Conf = Config{
		ServerPort: os.Getenv("TODO_PORT"),
		BasePath:   os.Getenv("TODO_DBFILE"),
	}

	if Conf.ServerPort == "" {
		Conf.ServerPort = "7540"
	}

	return nil
}
