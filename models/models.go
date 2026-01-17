package models

import (
	"database/sql"
	"time"
)

type User struct {
	UserID     int          `json:"user_id"`
	UserName   string       `json:"user_name"`
	Email      string       `json:"email"`
	Password   string       `json:"password"`
	Role       string       `json:"role"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
	IsVerified bool         `json:"is_verified"`
}

type Member struct {
	MemberID  int       `json:"member_id"`
	UserID    int       `json:"user_id"`
	TeamID    int       `json:"team_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"user_name,omitempty"`
	Email     string    `json:"email,omitempty"`
}

type Team struct {
	TeamID         int       `json:"team_id"`
	TeamName       string    `json:"team_name"`
	TeamLeaderID   int       `json:"team_leader_id"`
	TeamLeaderName string    `json:"team_leader_name,omitempty"`
	Members        []Member  `json:"members"`
	CreatedAt      time.Time `json:"created_at"`
}

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusInReview   TaskStatus = "in_review"
	TaskStatusDone       TaskStatus = "done"
)

func (ts TaskStatus) IsValid() bool {
	switch ts {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusInReview, TaskStatusDone:
		return true
	}
	return false
}

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

func (tp TaskPriority) IsValid() bool {
	switch tp {
	case TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh:
		return true
	}
	return false
}

type Task struct {
	TaskID       int            `json:"task_id"`
	TeamID       int            `json:"team_id"`
	CreatorID    int            `json:"creator_id"`
	AssigneeID   sql.NullInt64  `json:"assignee_id"`
	AssigneeName string         `json:"assignee_name,omitempty"`
	Title        string         `json:"title"`
	Description  sql.NullString `json:"description"`
	Status       TaskStatus     `json:"status"`
	Priority     TaskPriority   `json:"priority"`
	DueDate      sql.NullTime   `json:"due_date"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
}

type Comment struct {
	CommentID int       `json:"comment_id"`
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	MessageID int       `json:"message_id"`
	TeamID    int       `json:"team_id"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
