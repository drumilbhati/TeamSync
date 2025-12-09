package store

import (
	"database/sql"

	"github.com/drumilbhati/teamsync/models"
)

func (s *Store) CreateComment(c *models.Comment) error {
	err := s.db.QueryRow(
		`INSERT INTO comments (task_id, user_id, user_name, content)
		VALUES ($1, $2, $3, $4)
		RETURNING comment_id, created_at`,
		c.TaskID, c.UserID, c.UserName, c.Content,
	).Scan(&c.CommentID, &c.CreatedAt)

	return err
}

func (s *Store) GetCommentsByTaskID(taskID int) ([]models.Comment, error) {
	query := `
		SELECT c.comment_id, c.task_id, c.user_id, u.user_name, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.user_id
		WHERE c.task_id = $1
		ORDER BY c.created_at ASC`

	rows, err := s.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []models.Comment{}

	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.CommentID, &c.TaskID, &c.UserID, &c.UserName, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (s *Store) GetCommentbyID(comment_id int) (models.Comment, error) {
	var c models.Comment
	err := s.db.QueryRow(
		"SELECT comment_id, task_id, user_id, user_name, content, created_at FROM comments WHERE comment_id = $1",
		comment_id,
	).Scan(&c.CommentID, &c.TaskID, &c.UserID, &c.UserName, &c.Content, &c.CreatedAt)
	return c, err
}

func (s *Store) UpdateCommentByID(comment_id int, c *models.Comment) error {
	res, err := s.db.Exec(
		`UPDATE comments
		SET content = $1
		WHERE comment_id = $2`,
		c.Content, comment_id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) DeleteCommentByID(comment_id int) error {
	res, err := s.db.Exec(
		"DELETE FROM comments WHERE comment_id = $1",
		comment_id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
