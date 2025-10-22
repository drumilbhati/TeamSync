package main

import (
	"log"
	"net/http"
	"os"

	"github.com/drumilbhati/teamsync/controllers"
	"github.com/drumilbhati/teamsync/database"
	"github.com/drumilbhati/teamsync/store"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := database.Connect(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	s := store.NewStore(db)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer database.Close(db)

	r := mux.NewRouter()

	h := controllers.NewUserHandler(s)

	// Define routes
	r.HandleFunc("/users", h.GetUsers).Methods("GET")
	r.HandleFunc("/user/{id}", h.GetUserByID).Methods("GET")
	r.HandleFunc("/user", h.CreateUser).Methods("POST")
	r.HandleFunc("/user/{id}", h.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", h.DeleteUser).Methods("DELETE")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port: %s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
