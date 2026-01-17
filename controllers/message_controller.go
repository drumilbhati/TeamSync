package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/drumilbhati/teamsync/middleware"
	"github.com/drumilbhati/teamsync/store"
)

type MessageHandler struct {
	store *store.Store
}

func NewMessageHandler(s *store.Store) *MessageHandler {
	return &MessageHandler{store: s}
}

func (m *MessageHandler) GetMessagesByTeamID(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teamIDStr := r.URL.Query().Get("team_id")
	if teamIDStr == "" {
		http.Error(w, "team_id is required", http.StatusBadRequest)
		return
	}

	teamID, err := strconv.Atoi(teamIDStr)
	if err != nil {
		http.Error(w, "Invalid team_id", http.StatusBadRequest)
		return
	}

	// Verify membership
	isMember, err := m.store.IsTeamMember(requester_id, teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isMember {
		http.Error(w, "Unauthorized: you are not a member of this team", http.StatusForbidden)
		return
	}

	messages, err := m.store.GetMessagesByTeamID(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
