package main

import (
	"github.com/Enotisi/go_final_project/internal/actions"
	"github.com/Enotisi/go_final_project/internal/config"
	"github.com/Enotisi/go_final_project/internal/database"
	"github.com/Enotisi/go_final_project/internal/server"
)

func main() {
	config.InitializationConfig()
	db := database.MustBeInitDB()
	actions.InitializationAction(db)
	server.MustBeStartServer()
}
