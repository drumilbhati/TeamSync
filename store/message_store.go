package store

import (
	"github.com/drumilbhati/teamsync/models"
)

func (s *Store) CreateMessage(msg *models.Message) error {
	err := s.db.QueryRow(
		`INSERT INTO messages (team_id, user_id, user_name, content)
		VALUES ($1, $2, $3, $4)
		RETURNING message_id, created_at`,
		msg.TeamID, msg.UserID, msg.UserName, msg.Content,
	).Scan(&msg.MessageID, &msg.CreatedAt)

	return err
}

func (s *Store) GetMessagesByTeamID(teamID int) ([]models.Message, error) {
	rows, err := s.db.Query(
		`SELECT message_id, team_id, user_id, user_name, content, created_at
		FROM messages
		WHERE team_id = $1
		ORDER BY created_at ASC`,
		teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.MessageID, &msg.TeamID, &msg.UserID, &msg.UserName, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
