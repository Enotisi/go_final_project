package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	BasePath      string
	Password      string
	WebPath       string
	TokenLifeTime int
}

var Conf Config

func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	path, err := filepath.Abs("./web")

	if err != nil {
		return err
	}

	Conf = Config{
		ServerPort: os.Getenv("TODO_PORT"),
		BasePath:   os.Getenv("TODO_DBFILE"),
		Password:   os.Getenv("TODO_PASSWORD"),
		WebPath:    path,
	}

	if Conf.ServerPort == "" {
		Conf.ServerPort = "7540"
	}

	return nil
}
