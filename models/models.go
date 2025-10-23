package models

import (
	"database/sql"
	"time"
)

type User struct {
	UserID    int          `json:"user_id"`
	UserName  string       `json:"user_name"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	Role      string       `json:"role"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type Member struct {
	MemberID  int       `json:"member_id"`
	UserID    int       `json:"user_id"`
	TeamID    int       `json:"team_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Team struct {
	TeamID       int       `json:"team_id"`
	TeamName     string    `json:"team_name"`
	TeamLeaderID int       `json:"team_leader_id"`
	Members      []Member  `json:"members"`
	CreatedAt    time.Time `json:"created_at"`
}
