package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drumilbhati/teamsync/controllers"
	"github.com/drumilbhati/teamsync/database"
	"github.com/drumilbhati/teamsync/middleware"
	"github.com/drumilbhati/teamsync/store"
	"github.com/drumilbhati/teamsync/worker"
	"github.com/drumilbhati/teamsync/ws"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

// Upgrader is used to upgrade HTTP connection to a websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(hub *ws.Hub, s *store.Store, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)

	if !ok {
		fmt.Println("User not authenticated")
		return
	}

	teams, err := s.GetTeamsByUserID(userID)
	if err != nil {
		fmt.Println("Error fetching teams")
		return
	}
	var teamIDs []int
	for _, team := range teams {
		teamIDs = append(teamIDs, team.TeamID)
	}

	hub.AddUser(conn, teamIDs)

	defer hub.RemoveUser(conn, teamIDs)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

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

	redisAddr := os.Getenv("REDIS_ADDR")
	redisOpt := asynq.RedisClientOpt{Addr: redisAddr}

	// Create client for producers/controllers
	client := asynq.NewClient(redisOpt)
	defer client.Close()

	// Create and start sever for consumer/worker
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
		},
	)

	muxServer := asynq.NewServeMux()
	muxServer.HandleFunc(worker.TypeEmailDelivery, worker.HandleEmailDeliveryTask)

	// Run worker in background
	go func() {
		if err := srv.Run(muxServer); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	s := store.NewStore(db, rdb)

	defer database.Close(db)
	defer database.CloseRedis(rdb)

	r := mux.NewRouter()
	wsHub := ws.NewHub()

	u := controllers.NewUserHandler(s, client)
	t := controllers.NewTeamHandler(s)
	m := controllers.NewMemberHandler(s)
	k := controllers.NewTaskHandler(s, wsHub)
	c := controllers.NewCommentHandler(s)

	// Define routes
	// Websocket route

	// --- Public Auth Routes (changed prefix to /auth) ---
	r.HandleFunc("/auth/register", u.CreateUser).Methods("POST")
	r.HandleFunc("/auth/login", u.Login).Methods("POST")
	r.HandleFunc("/auth/verify", u.VerifyEmail).Methods("POST")

	// --- Protected API Routes ---
	// Create a subrouter that uses auth middleware
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Websocket routes
	api.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(wsHub, s, w, r)
	})

	// User routes
	api.HandleFunc("/users", u.GetUsers).Methods("GET")
	api.HandleFunc("/user/{id}", u.GetUserByID).Methods("GET")
	api.HandleFunc("/user/{id}", u.UpdateUserByID).Methods("PUT")
	api.HandleFunc("/user/{id}", u.DeleteUserByID).Methods("DELETE")

	// Team routes
	api.HandleFunc("/team/{id}", t.GetTeamByID).Methods("GET")
	api.HandleFunc("/team", t.GetTeamsByUserID).Methods("GET")
	api.HandleFunc("/team", t.GetTeamsByTeamLeaderID).Methods("GET").Queries("team_leader_id", "{id}")
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
	api.HandleFunc("/task", k.GetTasksByTeamIDWithPriority).Methods("GET").Queries("team_id", "{id}").Queries("priority", "{priority}")
	api.HandleFunc("/task", k.GetTasksByTeamIDWithStatus).Methods("GET").Queries("team_id", "{id}").Queries("status", "{status}")
	api.HandleFunc("/task", k.CreateTask).Methods("POST")
	api.HandleFunc("/task/{id}", k.UpdateTaskByID).Methods("PUT")
	api.HandleFunc("/task/{id}", k.DeleteTaskByID).Methods("DELETE")

	// Comment routes
	api.HandleFunc("/comment", c.CreateComment).Methods("POST")
	api.HandleFunc("/comment/{task_id}", c.GetCommentsByTaskID).Methods("GET")
	api.HandleFunc("/comment/{id}", c.UpdateCommentByID).Methods("PUT")
	api.HandleFunc("/comment/{id}", c.DeleteCommentByID).Methods("DELETE")

	// --- Start Server ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: middleware.CORSMiddleware(r),
	}

	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run HTTP server in a goroutine
	go func() {
		log.Printf("Server starting on port: %s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Block until a signal is received
	<-stop
	log.Println("Shutting down server...")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Shutdown Asynq worker
	srv.Shutdown()

	log.Println("Server gracefully stopped")
}
