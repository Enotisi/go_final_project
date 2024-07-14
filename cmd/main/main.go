package main

import (
	"log"

	"github.com/Enotisi/go_final_project/internal/actions"
	"github.com/Enotisi/go_final_project/internal/config"
	"github.com/Enotisi/go_final_project/internal/database"
	"github.com/Enotisi/go_final_project/internal/server"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading .env file. Error - %s", err.Error())
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error init database. Error - %s", err.Error())
	}

	actions.InitAction(db)
	err = server.StartServer()
	if err != nil {
		log.Fatalf("Error start server. Error - %s", err.Error())
	}
}
