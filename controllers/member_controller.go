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
	var member models.Member

	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := m.store.GetUserByID(member.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No user with given user_id found", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	_, err = m.store.GetTeamByID(member.TeamID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No team with given team_id found", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := m.store.CreateMember(&member); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (m *MemberHandler) UpdateMemberByID(w http.ResponseWriter, r *http.Request) {
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

	if err := m.store.UpdateMemberByID(member_id, &member); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (m *MemberHandler) DeleteMemberByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	member_id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := m.store.DeleteMemberByID(member_id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
