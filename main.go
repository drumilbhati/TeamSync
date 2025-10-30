package main

import (
	"log"
	"net/http"
	"os"

	"github.com/drumilbhati/teamsync/controllers"
	"github.com/drumilbhati/teamsync/database"
	"github.com/drumilbhati/teamsync/middleware"
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

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	rdb, err := database.ConnectRedis()

	if err != nil {
		log.Fatal("Failed to connect ot redis", err)
	}

	s := store.NewStore(db, rdb)

	defer database.Close(db)
	defer database.CloseRedis(rdb)

	r := mux.NewRouter()

	u := controllers.NewUserHandler(s)
	t := controllers.NewTeamHandler(s)
	m := controllers.NewMemberHandler(s)
	k := controllers.NewTaskHandler(s)

	// Define routes

	// --- Public Auth Routes (changed prefix to /auth) ---
	r.HandleFunc("/auth/register", u.CreateUser).Methods("POST")
	r.HandleFunc("/auth/login", u.Login).Methods("POST")
	r.HandleFunc("/auth/verify", u.VerifyEmail).Methods("POST")

	// --- Protected API Routes ---
	// Create a subrouter that uses auth middleware
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// User routes (paths are now relative to /api)
	api.HandleFunc("/users", u.GetUsers).Methods("GET")
	api.HandleFunc("/user/{id}", u.GetUserByID).Methods("GET")
	api.HandleFunc("/user/{id}", u.UpdateUserByID).Methods("PUT")
	api.HandleFunc("/user/{id}", u.DeleteUserByID).Methods("DELETE")

	// Team routes
	api.HandleFunc("/team/{id}", t.GetTeamByID).Methods("GET")
	api.HandleFunc("/team", t.GetTeamsByUserID).Methods("GET").Queries("user_id", "{id}")
	api.HandleFunc("/team", t.CreateTeam).Methods("POST")
	api.HandleFunc("/team/{id}", t.UpdateTeamByID).Methods("PUT")
	api.HandleFunc("/team/{id}", t.DeleteTeamByID).Methods("DELETE")

	// Member routes
	api.HandleFunc("/member/{id}", m.GetMemberByID).Methods("GET")
	api.HandleFunc("/member", m.GetMembersByTeamID).Methods("GET").Queries("team_id", "{id}")
	api.HandleFunc("/member", m.CreateMember).Methods("POST")
	api.HandleFunc("/member/{id}", m.UpdateMemberByID).Methods("PUT")
	api.HandleFunc("/member/{id}", m.DeleteMemberByID).Methods("DELETE")

	// Task routes
	api.HandleFunc("/task/{id}", k.GetTaskByTaskID).Methods("GET")
	api.HandleFunc("/task", k.GetTasksByTeamID).Methods("GET").Queries("team_id", "{id}")
	api.HandleFunc("/task", k.CreateTask).Methods("POST")
	api.HandleFunc("/task/{id}", k.UpdateTaskByID).Methods("PUT")
	api.HandleFunc("/task/{id}", k.DeleteTaskByID).Methods("DELETE")

	// --- Start Server ---
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port: %s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
