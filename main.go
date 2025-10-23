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

	u := controllers.NewUserHandler(s)
	t := controllers.NewTeamHandler(s)
	m := controllers.NewMemberHandler(s)

	// Define routes
	r.HandleFunc("/users", u.GetUsers).Methods("GET")
	r.HandleFunc("/user/{id}", u.GetUserByID).Methods("GET")
	r.HandleFunc("/user", u.CreateUser).Methods("POST")
	r.HandleFunc("/user/{id}", u.UpdateUserByID).Methods("PUT")
	r.HandleFunc("/user/{id}", u.DeleteUserByID).Methods("DELETE")

	r.HandleFunc("/team/{id}", t.GetTeamByID).Methods("GET")
	r.HandleFunc("/team", t.GetTeamsByUserID).Methods("GET").Queries("user_id", "{id}")
	r.HandleFunc("/team", t.CreateTeam).Methods("POST")
	r.HandleFunc("/team/{id}", t.UpdateTeamByID).Methods("PUT")
	r.HandleFunc("/team/{id}", t.DeleteTeamByID).Methods("DELETE")

	r.HandleFunc("/member/{id}", m.GetMemberByID).Methods("GET")
	r.HandleFunc("/member", m.GetMembersByTeamID).Methods("GET").Queries("team_id", "{id}")
	r.HandleFunc("/member", m.CreateMember).Methods("POST")
	r.HandleFunc("/member/{id}", m.UpdateMemberByID).Methods("PUT")
	r.HandleFunc("/member/{id}", m.DeleteMemberByID).Methods("DELETE")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port: %s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
