package store

import (
	"database/sql"
	"time"

	"github.com/drumilbhati/teamsync/models"
)

type Store struct {
	db *sql.DB
}

// create a new store by injecting db
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUsers() ([]models.User, error) {
	// we get all the rows
	rows, err := s.db.Query(
		"SELECT user_id, user_name, email, role, created_at, updated_at FROM users",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.UserName, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		u.Password = ""
		users = append(users, u)
	}
	return users, nil
}

func (s *Store) GetUserByID(id int) (*models.User, error) {
	var u models.User

	err := s.db.QueryRow(
		"SELECT user_id, user_name, email, role, created_at, updated_at FROM users WHERE user_id = $1",
		id,
	).Scan(&u.UserID, &u.UserName, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}
	u.Password = ""
	return &u, nil
}

func (s *Store) CreateUser(u *models.User) error {
	err := s.db.QueryRow(
		`INSERT INTO users (user_name, email, password, role) 
		VALUES ($1, $2, $3, $4) 
		RETURNING user_id, created_at`,
		u.UserName, u.Email, u.Password, u.Role,
	).Scan(&u.UserID, &u.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateUser(id int, u *models.User) error {
	_, err := s.db.Exec(
		`UPDATE users
		SET user_name = $1, email = $2, password = $3, role = $4, updated_at = $5
		WHERE user_id = $6`,
		u.UserName, u.Email, u.Password, time.Now(), id,
	)
	return err
}

func (s *Store) DeleteUser(id int) error {
	_, err := s.db.Exec(
		"DELETE FROM users WHERE user_id = $1", id,
	)
	return err
}
