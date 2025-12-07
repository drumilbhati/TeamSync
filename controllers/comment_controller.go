package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/drumilbhati/teamsync/models"
	"github.com/drumilbhati/teamsync/store"
	"github.com/gorilla/mux"
)

type CommentHandler struct {
	store *store.Store
}

func NewCommentHander(s *store.Store) *CommentHandler {
	return &CommentHandler{store: s}
}

func (c *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := c.store.GetTaskByTaskID(comment.TaskID)
	if err != nil {
		http.Error(w, "No task for given id found", http.StatusNotFound)
		return
	}

	_, err = c.store.GetUserByID(comment.UserID)
	if err != nil {
		http.Error(w, "No user for given id found", http.StatusNotFound)
		return
	}

	err = c.store.CreateComment(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func (c *CommentHandler) GetCommentsByTaskID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	taskId, err := strconv.Atoi(params["task_id"])
	if err != nil {
		http.Error(w, "Invalid task id", http.StatusBadRequest)
		return
	}

	_, err = c.store.GetTaskByTaskID(taskId)
	if err != nil {
		http.Error(w, "No task found for given task_id", http.StatusNotFound)
		return
	}

	comments, err := c.store.GetCommentsByTaskID(taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
