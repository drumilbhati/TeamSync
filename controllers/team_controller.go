package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/drumilbhati/teamsync/middleware"
	"github.com/drumilbhati/teamsync/models"
	"github.com/drumilbhati/teamsync/store"
	"github.com/gorilla/mux"
)

type TeamHandler struct {
	store *store.Store
}

func NewTeamHandler(s *store.Store) *TeamHandler {
	return &TeamHandler{store: s}
}

func (h *TeamHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	team_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	team, err := h.store.GetTeamByID(team_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) GetTeamsByTeamLeaderID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	team_leader_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teams, err := h.store.GetTeamsByTeamLeaderID(team_leader_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func (h *TeamHandler) GetTeamsByUserID(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teams, err := h.store.GetTeamsByUserID(requester_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	requesterID := r.Context().Value(middleware.UserIDKey).(int)
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	team.TeamLeaderID = requesterID

	if err := h.store.CreateTeam(&team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&team)
}

func (h *TeamHandler) UpdateTeamByID(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)

	team_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var updated_team models.Team
	json.NewDecoder(r.Body).Decode(&updated_team)

	team, err := h.store.GetTeamByID(team_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if team.TeamLeaderID != requester_id {
		http.Error(w, "Unauthorized: team leader can only update this team", http.StatusForbidden)
		return
	}

	if err := h.store.UpdateTeamByID(team_id, &updated_team); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&updated_team)
}

func (h *TeamHandler) DeleteTeamByID(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)

	team_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	team, err := h.store.GetTeamByID(team_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if team.TeamLeaderID != requester_id {
		http.Error(w, "Unauthorized: team leader can only delete this team", http.StatusForbidden)
		return
	}

	if err := h.store.DeleteTeamByID(team_id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
