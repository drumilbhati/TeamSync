package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB is a global variable that holds our database connection pool
// By making it a package-level variable, we can access it from anywhere
// sql.DB represents a pool of database connections, not a single connection
var DB *sql.DB

func Connect(host, port, user, password, dbname string) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// sql.Open prepares a database connection
	var err error
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		return err
	}

	// Ping actually attempts to connect to the database
	// sql.Open is lazy - it doesn't verify the connection is valid
	// DB.Ping() ensures we can actually communicate with the database
	if err = DB.Ping(); err != nil {
		// If we can't ping the database return err
		return err
	}

	log.Println("Successfully connected to database!")

	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
