package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drumilbhati/teamsync/controllers"
	"github.com/drumilbhati/teamsync/database"
	"github.com/drumilbhati/teamsync/logs"
	"github.com/drumilbhati/teamsync/middleware"
	"github.com/drumilbhati/teamsync/models"
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
		logs.Log.Info("User not authenticated")
		return
	}

	user, err := s.GetUserByID(userID)
	if err != nil {
		logs.Log.Errorf("Error fetching user: %v", err)
		return
	}

	teams, err := s.GetTeamsByUserID(userID)
	if err != nil {
		logs.Log.Error("Error fetching teams")
		return
	}
	var teamIDs []int
	for _, team := range teams {
		teamIDs = append(teamIDs, team.TeamID)
	}

	hub.AddUser(conn, teamIDs)

	defer hub.RemoveUser(conn, teamIDs)

	type Message struct {
		TeamID   int    `json:"team_id"`
		Content  string `json:"content"`
		UserID   int    `json:"user_id"`
		UserName string `json:"user_name"`
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			logs.Log.Errorf("Error unmarshalling message: %v", err)
			continue
		}

		// Verify the user is actually part of the team they are trying to message
		isMember := false
		for _, tid := range teamIDs {
			if tid == msg.TeamID {
				isMember = true
				break
			}
		}

		if isMember {
			msg.UserID = userID
			msg.UserName = user.UserName

			// Save to database
			dbMsg := models.Message{
				TeamID:   msg.TeamID,
				UserID:   msg.UserID,
				UserName: msg.UserName,
				Content:  msg.Content,
			}
			if err := s.CreateMessage(&dbMsg); err != nil {
				logs.Log.Errorf("Error saving message: %v", err)
			}

			updatedMessage, err := json.Marshal(msg)
			if err != nil {
				logs.Log.Errorf("Error marshalling message: %v", err)
				continue
			}
			hub.BroadcastToTeam(msg.TeamID, updatedMessage)
		}
	}
}

func main() {
	logs.InitLogger()
	defer logs.Sync()

	godotenv.Load()

	db, err := database.Connect(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		logs.Log.Fatal("Failed to connect to database: ", err)
	}

	rdb, err := database.ConnectRedis()
	if err != nil {
		logs.Log.Fatal("Failed to connect to redis: ", err)
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
			logs.Log.Fatalf("could not run server: %v", err)
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
	msgCtrl := controllers.NewMessageHandler(s)

	// Define routes
	// --- Public Auth Routes (changed prefix to /auth) ---
	r.HandleFunc("/auth/register", u.CreateUser).Methods("POST")
	r.HandleFunc("/auth/login", u.Login).Methods("POST")
	r.HandleFunc("/auth/verify", u.VerifyEmail).Methods("POST")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

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
	api.HandleFunc("/users/{id}", u.GetUserByID).Methods("GET")
	api.HandleFunc("/users/{id}", u.UpdateUserByID).Methods("PUT")
	api.HandleFunc("/users/{id}", u.DeleteUserByID).Methods("DELETE")

	// Team routes
	api.HandleFunc("/teams/{id}", t.GetTeamByID).Methods("GET")
	api.HandleFunc("/teams", t.GetTeamsByUserID).Methods("GET")
	api.HandleFunc("/teams", t.GetTeamsByTeamLeaderID).Methods("GET").Queries("team_leader_id", "{id}")
	api.HandleFunc("/teams", t.CreateTeam).Methods("POST")
	api.HandleFunc("/teams/{id}", t.UpdateTeamByID).Methods("PUT")
	api.HandleFunc("/teams/{id}", t.DeleteTeamByID).Methods("DELETE")

	// Member routes
	api.HandleFunc("/members/{id}", m.GetMemberByID).Methods("GET")
	api.HandleFunc("/members", m.GetMembersByTeamID).Methods("GET").Queries("team_id", "{id}")
	api.HandleFunc("/members", m.CreateMember).Methods("POST")
	api.HandleFunc("/members/{id}", m.UpdateMemberByID).Methods("PUT")
	api.HandleFunc("/members/{id}", m.DeleteMemberByID).Methods("DELETE")

	// Task routes
	api.HandleFunc("/tasks/{id}", k.GetTaskByTaskID).Methods("GET")
	api.HandleFunc("/tasks", k.GetTasksByTeamID).Methods("GET").Queries("team_id", "{id}")
	api.HandleFunc("/tasks", k.GetTasksByTeamIDWithPriority).Methods("GET").Queries("team_id", "{id}").Queries("priority", "{priority}")
	api.HandleFunc("/tasks", k.GetTasksByTeamIDWithStatus).Methods("GET").Queries("team_id", "{id}").Queries("status", "{status}")
	api.HandleFunc("/tasks", k.CreateTask).Methods("POST")
	api.HandleFunc("/tasks/{id}", k.UpdateTaskByID).Methods("PUT")
	api.HandleFunc("/tasks/{id}", k.DeleteTaskByID).Methods("DELETE")

	// Comment routes
	api.HandleFunc("/comments", c.CreateComment).Methods("POST")
	api.HandleFunc("/comments/{task_id}", c.GetCommentsByTaskID).Methods("GET")
	api.HandleFunc("/comments/{id}", c.UpdateCommentByID).Methods("PUT")
	api.HandleFunc("/comments/{id}", c.DeleteCommentByID).Methods("DELETE")

	// Message routes
	api.HandleFunc("/messages", msgCtrl.GetMessagesByTeamID).Methods("GET").Queries("team_id", "{id}")

	// Copilot route
	api.HandleFunc("/enhance", controllers.Describe).Methods("POST")
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
		logs.Log.Infof("Server starting on port: %s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logs.Log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Block until a signal is received
	<-stop
	logs.Log.Info("Shutting down server...")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		logs.Log.Errorf("HTTP server shutdown error: %v", err)
	}

	// Shutdown Asynq worker
	srv.Shutdown()

	logs.Log.Info("Server gracefully stopped")
}
