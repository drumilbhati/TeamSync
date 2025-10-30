package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/drumilbhati/teamsync/models"
	"github.com/drumilbhati/teamsync/store"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	store *store.Store
}

func NewTaskHandler(s *store.Store) *TaskHandler {
	return &TaskHandler{store: s}
}

func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := t.store.CreateTask(&task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (t *TaskHandler) GetTaskByTaskID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	task_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, "Invalid request params", http.StatusBadRequest)
		return
	}

	task, err := t.store.GetTaskByTaskID(task_id)
	if err != nil {
		http.Error(w, "Not task found with given id", http.StatusNotFound)
		return
	}

	fmt.Println("Task ID: ", task.TaskID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (t *TaskHandler) GetTasksByTeamID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	team_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, "Invalid request params", http.StatusBadRequest)
		return
	}

	task, err := t.store.GetTasksByTeamID(team_id)

	if err != nil {
		http.Error(w, "Not task found with given id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (t *TaskHandler) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	task_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, "Invalid request params", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task details", http.StatusBadRequest)
		return
	}

	if err := t.store.UpdateTaskByID(task_id, &task); err != nil {
		http.Error(w, "Error updating the task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (t *TaskHandler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	task_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, "Invalid request params", http.StatusBadRequest)
		return
	}

	if err := t.store.DeleteTaskByID(task_id); err != nil {
		http.Error(w, "Error deleting the task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
