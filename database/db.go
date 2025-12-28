package database

import (
	"database/sql"
	"fmt"

	"github.com/drumilbhati/teamsync/logs"

	_ "github.com/lib/pq"
)

func Connect(host, port, user, password, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// sql.Open prepares a database connection
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping actually attempts to connect to the database
	// sql.Open is lazy - it doesn't verify the connection is valid
	// db.Ping() ensures we can actually communicate with the database
	if err = db.Ping(); err != nil {
		// If we can't ping the database return err
		return nil, err
	}

	logs.Log.Info("Successfully connected to database!")

	return db, nil
}

func Close(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
