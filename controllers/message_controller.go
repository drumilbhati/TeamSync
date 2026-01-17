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
	members, err := m.store.GetMembersByTeamID(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isMember := false
	for _, member := range members {
		if member.UserID == requester_id {
			isMember = true
			break
		}
	}
	// Also check if leader (though typically leader is also a member, but let's be safe if logic differs)
	if !isMember {
		team, err := m.store.GetTeamByID(teamID)
		if err == nil && team.TeamLeaderID == requester_id {
			isMember = true
		}
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
