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

func (s *Store) GetTaskByTaskID(task_id int) (*models.Task, error) {
	var t models.Task
	err := s.db.QueryRow(
		"SELECT task_id, team_id, creator_id, assignee_id, title, description, status, priority, due_date, created_at, updated_at FROM tasks WHERE task_id = $1",
		task_id,
	).Scan(&t.TaskID, &t.TeamID, &t.CreatorID, &t.AssigneeID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Store) GetTasksByTeamID(team_id int) ([]models.Task, error) {
	rows, err := s.db.Query(
		"SELECT task_id, team_id, creator_id, assignee_id, title, description, status, priority, due_date, created_at, updated_at FROM tasks WHERE team_id = $1",
		team_id,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.TaskID, &t.TeamID, &t.CreatorID, &t.AssigneeID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Store) GetTaskByTeamIDWithPriority(team_id int, priority models.TaskPriority) ([]models.Task, error) {
	rows, err := s.db.Query(
		"SELECT task_id, team_id, creator_id, assignee_id, title, description, status, priority, due_date, created_at, updated_at FROM tasks WHERE team_id = $1 AND priority = $2", team_id, priority,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.TaskID, &t.TeamID, &t.CreatorID, &t.AssigneeID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Store) UpdateTaskByID(task_id int, t *models.Task) error {
	_, err := s.db.Exec(
		`UPDATE tasks
		SET assignee_id = $1, description = $2, status = $3, priority = $4, due_date = $5, updated_at = $6 WHERE task_id = $7`,
		t.AssigneeID, t.Description, t.Status, t.Priority, t.DueDate, time.Now(), task_id,
	)

	return err
}

func (s *Store) DeleteTaskByID(task_id int) error {
	_, err := s.db.Exec(
		`DELETE FROM tasks WHERE task_id = $1`,
		task_id,
	)
	return err
}
