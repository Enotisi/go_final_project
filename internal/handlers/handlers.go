package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Enotisi/go_final_project/internal/actions"
	"github.com/Enotisi/go_final_project/internal/models"
)

func NextDateHandle(w http.ResponseWriter, r *http.Request) {

	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nowDate, err := time.Parse(actions.DateTemplate, now)

	if err != nil {
		http.Error(w, "Недопустимый формат даты", http.StatusBadRequest)
		return
	}

	nextDate, err := actions.NextDate(nowDate, date, repeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if _, err := w.Write([]byte(nextDate.Format(actions.DateTemplate))); err != nil {
		log.Printf("error return responce at NextDateHandle: %s\n", err.Error())
	}
}

func CreateTaskHandle(w http.ResponseWriter, r *http.Request) {

	taskData := models.Task{}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &taskData); err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	id, err := actions.CreateTask(taskData)

	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	resp := createJsonResponse("id", strconv.Itoa(id))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(resp)); err != nil {
		log.Printf("error return responce at CreateTaskHandle: %s\n", err.Error())
	}
}

func TasksListHandle(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")

	tasks, err := actions.TasksList(search)

	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(map[string][]models.Task{"tasks": tasks})
	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(resp)); err != nil {
		log.Printf("error return responce at TasksListHandle: %s\n", err.Error())
	}

}

func GetTaskHandle(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, createJsonResponse("error", "Не указан идентификатор"), http.StatusBadRequest)
		return
	}

	task, err := actions.GetTaskById(id)
	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		log.Printf("error return responce at GetTaskHandle: %s\n", err.Error())
	}
}

func UpdateTaskHandle(w http.ResponseWriter, r *http.Request) {

	taskData := models.Task{}
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &taskData); err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	err = actions.UpdateTask(taskData, true)
	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func DoneTaskHandle(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, createJsonResponse("error", "Не указан идентификатор"), http.StatusBadRequest)
		return
	}

	err := actions.DoneTask(id)
	if err != nil {
		http.Error(w, createJsonResponse("error", err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func DeleteTaskHandle(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, createJsonResponse("error", "Не указан идентификатор"), http.StatusBadRequest)
		return
	}

	err := actions.DeleteTaskById(id)
	if err != nil {
		http.Error(w, createJsonResponse("error", "Не указан идентификатор"), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func createJsonResponse(title, text string) string {
	return fmt.Sprintf(`{"%s":"%s"}`, title, text)
}
