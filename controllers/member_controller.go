package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/drumilbhati/teamsync/middleware"
	"github.com/drumilbhati/teamsync/models"
	"github.com/drumilbhati/teamsync/store"
	"github.com/gorilla/mux"
)

type MemberHandler struct {
	store *store.Store
}

func NewMemberHandler(s *store.Store) *MemberHandler {
	return &MemberHandler{store: s}
}

func (m *MemberHandler) GetMemberByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	member_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	member, err := m.store.GetMemberByID(member_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (m *MemberHandler) GetMembersByTeamID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	team_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	members, err := m.store.GetMembersByTeamID(team_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func (m *MemberHandler) CreateMember(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		TeamID int    `json:"team_id"`
		UserID int    `json:"user_id"`
		Email  string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var targetUserID int

	if req.Email != "" {
		user, err := m.store.GetUserByEmail(req.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No user with given email found", http.StatusNotFound)
			} else if err.Error() == "user not verified" {
				http.Error(w, "User exists but is not verified", http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		targetUserID = user.UserID
	} else if req.UserID != 0 {
		targetUserID = req.UserID
	} else {
		http.Error(w, "Either user_id or email must be provided", http.StatusBadRequest)
		return
	}

	// Verify target user exists (redundant if looked up by email, but safe for ID)
	if req.Email == "" {
		_, err := m.store.GetUserByID(targetUserID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No user with given user_id found", http.StatusBadRequest)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	team, err := m.store.GetTeamByID(req.TeamID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No team with given team_id found", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if team.TeamLeaderID != requester_id {
		http.Error(w, "Unauthorized: only team leader can add members", http.StatusForbidden)
		return
	}

	member := models.Member{
		TeamID: req.TeamID,
		UserID: targetUserID,
		Role:   "member", // Default role
	}

	if err := m.store.CreateMember(&member); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (m *MemberHandler) UpdateMemberByID(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)

	member_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var member models.Member

	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mem, err := m.store.GetMemberByID(member_id)
	if err != nil {
		http.Error(w, "Error getting member details", http.StatusInternalServerError)
		return
	}

	team, err := m.store.GetTeamByID(mem.TeamID)
	if err != nil {
		http.Error(w, "Error getting team details", http.StatusInternalServerError)
		return
	}

	if team.TeamLeaderID != requester_id {
		http.Error(w, "Unauthorized: only the team leader can edit member details", http.StatusForbidden)
		return
	}

	if err := m.store.UpdateMemberByID(member_id, &member); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (m *MemberHandler) DeleteMemberByID(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)

	member_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mem, err := m.store.GetMemberByID(member_id)
	if err != nil {
		http.Error(w, "Error getting member details", http.StatusInternalServerError)
		return
	}

	team, err := m.store.GetTeamByID(mem.TeamID)
	if err != nil {
		http.Error(w, "Error getting team details", http.StatusInternalServerError)
		return
	}

	if team.TeamLeaderID != requester_id && mem.MemberID != requester_id {
		http.Error(w, "Unauthorized: only the team leader can delete member details", http.StatusForbidden)
		return
	}

	if err := m.store.DeleteMemberByID(member_id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
