package config

import (
	"os"
	"path/filepath"
	"strconv"

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
		ServerPort:    os.Getenv("TODO_PORT"),
		BasePath:      os.Getenv("TODO_DBFILE"),
		Password:      os.Getenv("TODO_PASSWORD"),
		WebPath:       path,
		TokenLifeTime: 0,
	}

	if Conf.ServerPort == "" {
		Conf.ServerPort = "7540"
	}

	lifetime, err := strconv.Atoi(os.Getenv("TODO_TOKEN_LIFETIME"))
	if err != nil {
		Conf.TokenLifeTime = 8
	} else {
		Conf.TokenLifeTime = lifetime
	}

	return nil
}
