package server

import (
	"net/http"

	"github.com/Enotisi/go_final_project/internal/config"
	"github.com/Enotisi/go_final_project/internal/handlers"
	"github.com/go-chi/chi"
)

func StartServer() error {

	r := chi.NewRouter()

	r.HandleFunc("/*", handlers.WebHandler)
	r.Post("/api/sign", handlers.SignHandle)

	r.Group(func(r chi.Router) {
		r.Use(handlers.MiddlewareHandle)
		r.Get("/api/nextdate", handlers.NextDateHandle)
		r.Get("/api/tasks", handlers.TasksListHandle)
		r.Get("/api/task", handlers.GetTaskHandle)
		r.Post("/api/task", handlers.CreateTaskHandle)
		r.Delete("/api/task", handlers.DeleteTaskHandle)
		r.Post("/api/task/done", handlers.DoneTaskHandle)
		r.Put("/api/task", handlers.UpdateTaskHandle)
	})

	port := config.Conf.ServerPort

	if err := http.ListenAndServe(":"+port, r); err != nil {
		return err
	}

	return nil
}
