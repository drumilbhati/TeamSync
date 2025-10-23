package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *TeamHandler) GetTeamsByUserID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	user_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	team, err := h.store.GetTeamsByUserID(user_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := h.store.GetUserByID(team.TeamLeaderID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid team_leader_id: user does not exit", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := h.store.CreateTeam(&team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&team)
}

func (h *TeamHandler) UpdateTeamByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	team_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var team models.Team
	json.NewDecoder(r.Body).Decode(&team)

	if err := h.store.UpdateTeamByID(team_id, &team); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&team)
}

func (h *TeamHandler) DeleteTeamByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	team_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.store.DeleteTeamByID(team_id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
