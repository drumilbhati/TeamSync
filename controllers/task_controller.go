package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/drumilbhati/teamsync/middleware"
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
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if task.CreatorID != requester_id {
		http.Error(w, "Unauthorized: You can use only your user_id to create task", http.StatusForbidden)
		return
	}

	team, err := t.store.GetTeamByID(task.TeamID)
	if err != nil {
		http.Error(w, "Error getting team details", http.StatusInternalServerError)
		return
	}

	is_member_of_team := false
	if team.TeamLeaderID == requester_id {
		is_member_of_team = true
	} else {
		for _, m := range team.Members {
			if m.UserID == requester_id {
				is_member_of_team = true
				break
			}
		}
	}

	if !is_member_of_team {
		http.Error(w, "Unauthorized: You must be a member of the team to create tasks", http.StatusForbidden)
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
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)

	task_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, "Invalid request params", http.StatusBadRequest)
		return
	}

	var updated_task models.Task
	if err := json.NewDecoder(r.Body).Decode(&updated_task); err != nil {
		http.Error(w, "Invalid task details", http.StatusBadRequest)
		return
	}

	task, err := t.store.GetTaskByTaskID(task_id)
	if err != nil {
		http.Error(w, "Unable to get task", http.StatusInternalServerError)
		return
	}

	if requester_id != task.CreatorID {
		http.Error(w, "Unauthorized: Only the creator can edit this task", http.StatusForbidden)
		return
	}

	if err := t.store.UpdateTaskByID(task_id, &updated_task); err != nil {
		http.Error(w, "Error updating the task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated_task)
}

func (t *TaskHandler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)

	task_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, "Invalid request params", http.StatusBadRequest)
		return
	}

	task, err := t.store.GetTaskByTaskID(task_id)
	if err != nil {
		http.Error(w, "Cannot get task for given id", http.StatusInternalServerError)
		return
	}

	if task.CreatorID != requester_id {
		http.Error(w, "Unauthorized: only the creator can delete this task", http.StatusForbidden)
		return
	}

	if err := t.store.DeleteTaskByID(task_id); err != nil {
		http.Error(w, "Error deleting the task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
