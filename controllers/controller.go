package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/drumilbhati/teamsync/database"
	"github.com/drumilbhati/teamsync/models"
	"github.com/gorilla/mux"
)

/*
GetUsers handles GET request to /users
It retrives all users from the database
w (ResponseWriter) is used to write the HTTP response
r (*Request) contains the HTTP request information
*/
func GetUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Getting all users")

	// Execute a SQL query to get all users
	// database.DB.Query() runs a query that returns multiple rows
	// It returns a *Rows object that we can iterate over, and an error
	rows, err := database.DB.Query("SELECT id, name, email, created_at FROM users")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// This is important to free the databse resources
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User

		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

/*
GetUser handles a GET request at /user/{id}
It retrives a single user by id
*/
func GetUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Getting user by id")

	// Extract the params from the request
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var u models.User

	// QueryRow executes a query that returns atmost one row
	err = database.DB.QueryRow(
		"SELECT id, name, email, created_at FROM users WHERE id = $1", id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

/*
CreateUser handles a POST request at /user
It creates a new user in the database
*/
func CreateUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Creating new user")

	var u models.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at",
		u.Name, u.Email,
	).Scan(&u.ID, &u.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// WriteHeader sets the status code to 201 -> Created (successful creation)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

/*
UpdateUser hadles PUT request at /user/{id}
It updats the user
*/
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Updating a user")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Execute the UPDATE query
	// Exec is used for queries that don't return rows (INSERT, UPDATE, DELETE)
	// Returns a Result object and an error
	// We use _ to ignore the Result
	_, err = database.DB.Exec(
		"UPDATE users SET name = $1, email = $2 WHERE id = $3",
		u.Name, u.Email, id,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

/*
DeleteUser handles a DELETE request at /user/{id}
It deletes the user from Database
*/
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Deleting a user")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete a user with given id
	_, err = database.DB.Exec(
		"DELETE FROM users WHERE id = $1", id,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// WriteHeader sends the status code
	// 204 = No Content (successful deletion, no body to return)
	// When using WriteHeader, no body should be written after
	w.WriteHeader(http.StatusNoContent)
}
