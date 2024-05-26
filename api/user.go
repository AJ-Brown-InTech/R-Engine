package api

import (
	"Engine/storage"
	"Engine/types"
	"encoding/json"
	"net/http"
	"time"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	
)

// RegisterAccount registers a new user account
func RegisterAccount(db *storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user types.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate user input
		if err := user.Validate(); err != nil {
			http.Error(w, "Invalid user data", http.StatusBadRequest)
			return
		}

		// Check if user or email is already registered
		userExists, emailExists, err := user.UserAndEmailExist(r.Context(), db.Db)
		if err != nil {
			http.Error(w, "Failed to check existing users", http.StatusInternalServerError)
			return
		}
		if userExists {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		if emailExists {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		// Generate salt and hash the password
		hashedPassword, salt, err := types.Hash(user.UserPassword)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.UserPassword = hashedPassword
		user.Salt = salt

		// Set default values
		user.UserID = uuid.New().String()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		// Insert user into database
		if err := user.Create(r.Context(), db.Db); err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Generate a session token
		var session types.Session
		sessionToken, err := session.Create(user.Username)
		if err != nil {
			http.Error(w, "Failed to create session token", http.StatusInternalServerError)
			return
		}

		// Set the session token in the response headers
		w.Header().Set("Authorization", sessionToken)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

// CheckUserExistence checks if a username or email exists in the database
func Exist(db *storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Username string `json:"username,omitempty"`
			Email    string `json:"email,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var user types.User
		user.Username = request.Username
		user.Email = request.Email

		response := make(map[string]bool)

		if request.Username != "" {
			usernameExists, err := user.UserExist(r.Context(), db.Db)
			if err != nil {
				http.Error(w, "Failed to check username", http.StatusInternalServerError)
				return
			}
			response["username_exists"] = usernameExists
		}

		if request.Email != "" {
			emailExists, err := user.EmailExist(r.Context(), db.db)
			if err != nil {
				http.Error(w, "Failed to check email", http.StatusInternalServerError)
				return
			}
			response["email_exists"] = emailExists
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

// PatchAccount updates user account details
func PatchAccount(db *storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate user input
		v := validator.New()
		if err := v.Var(updates["email"], "omitempty,email"); err != nil {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}
		if err := v.Var(updates["username"], "omitempty,min=8,max=15"); err != nil {
			http.Error(w, "Invalid username format", http.StatusBadRequest)
			return
		}

		// // Update user in database
		// user := &types.User{}
		// if err := user.PartialUpdate(context.Background(), db.DB, updates, userID); err != nil {
		// 	http.Error(w, "Failed to update user", http.StatusInternalServerError)
		// 	return
		// }

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(updates); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}
