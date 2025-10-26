package store

/*
	APIs needed:
	GET:
	GetUsers
	GetUserByID
	GetUserByEmail

	POST:
	CreateUser

	PUT:
	UpdateUserByID

	DELETE:
	DeleteUserByID
*/

import (
	"context"
	"fmt"
	"time"

	"github.com/drumilbhati/teamsync/models"
	"github.com/drumilbhati/teamsync/utils"
	"github.com/redis/go-redis/v9"
)

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

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	var u models.User

	err := s.db.QueryRow(
		"SELECT user_id, user_name, email, password, role, created_at, updated_at, is_verified FROM users WHERE email = $1",
		email,
	).Scan(&u.UserID, &u.UserName, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt, &u.IsVerified)

	if err != nil {
		return nil, err
	}

	if !u.IsVerified {
		return nil, fmt.Errorf("user not verified")
	}

	return &u, nil
}

func (s *Store) VerifyUser(userID int) error {
	_, err := s.db.Exec(
		"UPDATE users SET is_verified = true WHERE user_id = $1", userID,
	)

	return err
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

func (s *Store) UpdateUserByID(id int, u *models.User) error {
	_, err := s.db.Exec(
		`UPDATE users
		SET user_name = $1, email = $2, password = $3, role = $4, updated_at = $5
		WHERE user_id = $6`,
		u.UserName, u.Email, u.Password, time.Now(), id,
	)
	return err
}

func (s *Store) DeleteUserByID(id int) error {
	_, err := s.db.Exec(
		"DELETE FROM users WHERE user_id = $1", id,
	)
	return err
}

func (s *Store) GetUserByEmailForAuth(email string) (*models.User, error) {
	var u models.User
	err := s.db.QueryRow(
		"SELECT user_id, user_name, email, password, role, created_at, updated_at, is_verified FROM users WHERE email = $1",
		email,
	).Scan(&u.UserID, &u.UserName, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt, &u.IsVerified)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func GenerateOTP() string {
	return utils.GenerateRandomNumber()
}

// Save the otp for a user in Redis with 10 minute time limit
func (s *Store) CreateOTP(userID int, otp string) error {
	ctx := context.Background()
	key := fmt.Sprintf("otp:%d", userID)

	err := s.rdb.Set(ctx, key, otp, 10*time.Minute).Err()
	return err
}

// check if the otp is valid
func (s *Store) GetValidOTP(userID int, otp string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("otp:%d", userID)

	storedOTP, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// key does not exist, or has expired
		return false, nil
	} else if err != nil {
		return false, err
	}

	return storedOTP == otp, nil
}

func (s *Store) DeleteOTP(userID int) error {
	ctx := context.Background()
	key := fmt.Sprintf("otp:%d", userID)
	return s.rdb.Del(ctx, key).Err()
}
