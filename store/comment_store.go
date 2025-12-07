package store

import (
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
