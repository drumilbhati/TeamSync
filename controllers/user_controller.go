package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/drumilbhati/teamsync/logs"
	"github.com/drumilbhati/teamsync/middleware"
	"github.com/drumilbhati/teamsync/models"
	"github.com/drumilbhati/teamsync/store"
	"github.com/drumilbhati/teamsync/worker"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	store  *store.Store
	client *asynq.Client
}

func NewUserHandler(s *store.Store, c *asynq.Client) *UserHandler {
	return &UserHandler{store: s, client: c}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.store.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/type")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingUser, err := h.store.GetUserByEmailForAuth(user.Email)

	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingUser != nil && existingUser.IsVerified {
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	if existingUser != nil {
		user.UserID = existingUser.UserID
	} else {
		if err := h.store.CreateUser(&user); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	}

	otp := store.GenerateOTP()

	if err := h.store.CreateOTP(user.UserID, otp); err != nil {
		http.Error(w, "Failed to save OTP to Redis", http.StatusInternalServerError)
		return
	}

	task, err := worker.NewEmailDeliveryTask(user.Email, user.UserName, otp)
	if err != nil {
		http.Error(w, "Failed to create email task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.client.Enqueue(task); err != nil {
		http.Error(w, "Failed to enqueue email task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = ""
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.store.GetUserByEmailForAuth(req.Email)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.IsVerified {
		http.Error(w, "User already verified", http.StatusBadRequest)
		return
	}

	// check OTP from Redis
	isValid, err := h.store.GetValidOTP(user.UserID, req.OTP)
	if err != nil {
		http.Error(w, "Error checking OTP", http.StatusInternalServerError)
		return
	}

	if !isValid {
		http.Error(w, "Invalid or expired OTP", http.StatusBadRequest)
		return
	}

	// mark as verified in SQL
	if err := h.store.VerifyUser(user.UserID); err != nil {
		http.Error(w, "Failed to verify user", http.StatusInternalServerError)
		return
	}

	// delete the OTP from redis
	if err := h.store.DeleteOTP(user.UserID); err != nil {
		logs.Log.Warnf("Failed to delete OTP for user %d: %s", user.UserID, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email verified successfully. Please log in.",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// get user by email
	user, err := h.store.GetUserByEmail(loginReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email", http.StatusUnauthorized)
		} else if err.Error() == "user not verified" {
			http.Error(w, "Account not verified. Please check your email.", http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// create a token
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// sign the token with the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the token
	w.Header().Set("Content-Type", "applicaion/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func (h *UserHandler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)

	user_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user_id != requester_id {
		http.Error(w, "Unauthorized: you can edit only your profile", http.StatusForbidden)
		return
	}

	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	if err := h.store.UpdateUserByID(user_id, &user); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/type")
	json.NewEncoder(w).Encode("Updated user successfully")
}

func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	requester_id, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)

	user_id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user_id != requester_id {
		http.Error(w, "Unauthorized: you can delete only your profile", http.StatusForbidden)
		return
	}

	if err := h.store.DeleteUserByID(user_id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
