package store

import (
	"time"

	"github.com/drumilbhati/teamsync/models"
)

func (s *Store) CreateTask(t *models.Task) error {
	err := s.db.QueryRow(
		`INSERT INTO tasks (team_id, creator_id, assignee_id, title, description, status, priority, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING task_id, created_at`,
		&t.TeamID, &t.CreatorID, &t.AssigneeID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate,
	).Scan(&t.TaskID, &t.CreatedAt)

	return err
}

func (s *Store) GetTaskByTaskID(taskID int) (*models.Task, error) {
	var t models.Task
	var assigneeName *string
	err := s.db.QueryRow(
		`SELECT t.task_id, t.team_id, t.creator_id, t.assignee_id, u.user_name, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.user_id
		WHERE t.task_id = $1`,
		taskID,
	).Scan(&t.TaskID, &t.TeamID, &t.CreatorID, &t.AssigneeID, &assigneeName, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if assigneeName != nil {
		t.AssigneeName = *assigneeName
	}
	return &t, nil
}

func (s *Store) GetTasksByTeamID(teamID int) ([]models.Task, error) {
	rows, err := s.db.Query(
		`SELECT t.task_id, t.team_id, t.creator_id, t.assignee_id, u.user_name, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.user_id
		WHERE t.team_id = $1`,
		teamID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var assigneeName *string
		if err := rows.Scan(&t.TaskID, &t.TeamID, &t.CreatorID, &t.AssigneeID, &assigneeName, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if assigneeName != nil {
			t.AssigneeName = *assigneeName
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Store) GetTasksByTeamIDWithPriority(teamID int, priority models.TaskPriority) ([]models.Task, error) {
	rows, err := s.db.Query(
		`SELECT t.task_id, t.team_id, t.creator_id, t.assignee_id, u.user_name, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.user_id
		WHERE t.team_id = $1 AND t.priority = $2`, teamID, priority,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var assigneeName *string
		if err := rows.Scan(&t.TaskID, &t.TeamID, &t.CreatorID, &t.AssigneeID, &assigneeName, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if assigneeName != nil {
			t.AssigneeName = *assigneeName
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Store) GetTasksByTeamIDWithStatus(teamID int, status models.TaskStatus) ([]models.Task, error) {
	rows, err := s.db.Query(
		`SELECT t.task_id, t.team_id, t.creator_id, t.assignee_id, u.user_name, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.user_id
		WHERE t.team_id = $1 AND t.status = $2`, teamID, status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var assigneeName *string
		if err := rows.Scan(&t.TaskID, &t.TeamID, &t.CreatorID, &t.AssigneeID, &assigneeName, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if assigneeName != nil {
			t.AssigneeName = *assigneeName
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Store) UpdateTaskByID(taskID int, t *models.Task) error {
	_, err := s.db.Exec(
		`UPDATE tasks
		SET title = $1, assignee_id = $2, description = $3, status = $4, priority = $5, due_date = $6, updated_at = $7 WHERE task_id = $8`,
		t.Title, t.AssigneeID, t.Description, t.Status, t.Priority, t.DueDate, time.Now(), taskID,
	)

	return err
}

func (s *Store) DeleteTaskByID(taskID int) error {
	_, err := s.db.Exec(
		`DELETE FROM tasks WHERE task_id = $1`,
		taskID,
	)
	return err
}
