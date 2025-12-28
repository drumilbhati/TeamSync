package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/drumilbhati/teamsync/logs"
	"github.com/drumilbhati/teamsync/middleware"
	"github.com/drumilbhati/teamsync/models"
	"github.com/drumilbhati/teamsync/store"
	"github.com/drumilbhati/teamsync/ws"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	store *store.Store
	wsHub *ws.Hub
}

func NewTaskHandler(s *store.Store, wsHub *ws.Hub) *TaskHandler {
	return &TaskHandler{store: s, wsHub: wsHub}
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
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

	msg := Message{
		Type: "TASK_CREATED",
		Data: task,
	}
	msgBytes, _ := json.Marshal(msg)

	t.wsHub.BroadcastToTeam(task.TeamID, msgBytes)

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

	logs.Log.Info("Task ID: ", task.TaskID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (t *TaskHandler) GetTasksByTeamIDWithPriority(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	team_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid request params", http.StatusBadRequest)
		return
	}

	priorityStr := r.URL.Query().Get("priority")
	priority := models.TaskPriority(priorityStr)
	if !priority.IsValid() {
		http.Error(w, "Invalid priority", http.StatusBadRequest)
		return
	}

	tasks, err := t.store.GetTasksByTeamIDWithPriority(team_id, models.TaskPriority(priority))
	if err != nil {
		http.Error(w, "No task found with given team_id and priority", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (t *TaskHandler) GetTasksByTeamIDWithStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	team_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid request params", http.StatusBadRequest)
		return
	}

	statusStr := r.URL.Query().Get("status")
	status := models.TaskStatus(statusStr)

	if !status.IsValid() {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	tasks, err := t.store.GetTasksByTeamIDWithStatus(team_id, models.TaskStatus(status))
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (t *TaskHandler) GetTasksByTeamID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	team_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid request params", http.StatusBadRequest)
		return
	}

	tasks, err := t.store.GetTasksByTeamID(team_id)
	if err != nil {
		http.Error(w, "Not task found with given id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
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

	updated_task.TeamID = task.TeamID

	msg := Message{
		Type: "TASK_UPDATED",
		Data: updated_task,
	}

	msg_bytes, _ := json.Marshal(msg)

	t.wsHub.BroadcastToTeam(updated_task.TeamID, msg_bytes)

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
	msg := Message{
		Type: "TASK_DELETED",
		Data: map[string]int{"task_id": task_id},
	}

	msg_bytes, _ := json.Marshal(msg)

	t.wsHub.BroadcastToTeam(task.TeamID, msg_bytes)

	w.WriteHeader(http.StatusNoContent)
}
