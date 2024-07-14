package server

import (
	"net/http"
	"path/filepath"

	"github.com/Enotisi/go_final_project/internal/config"
	"github.com/Enotisi/go_final_project/internal/handlers"
	"github.com/go-chi/chi"
)

func StartServer() error {

	path, err := filepath.Abs("./web")

	if err != nil {
		return err
	}

	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir(path)))
	r.Get("/api/nextdate", handlers.NextDateHandle)
	r.Get("/api/tasks", handlers.TasksListHandle)
	r.Get("/api/task", handlers.GetTaskHandle)
	r.Post("/api/task", handlers.CreateTaskHandle)
	r.Delete("/api/task", handlers.DeleteTaskHandle)
	r.Post("/api/task/done", handlers.DoneTaskHandle)
	r.Put("/api/task", handlers.UpdateTaskHandle)

	port := config.Conf.ServerPort

	if err := http.ListenAndServe(":"+port, r); err != nil {
		return err
	}

	return nil
}
