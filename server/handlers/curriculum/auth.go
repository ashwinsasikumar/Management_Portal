package curriculum

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"server/db"
	"server/models"

	"golang.org/x/crypto/bcrypt"
)

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var loginReq models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		log.Println("Error decoding login request:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	log.Printf("Login attempt for username: %s", loginReq.Username)

	// Query user from database
	var user models.User
	query := `SELECT id, username, password_hash, full_name, email, role, is_active, created_at, updated_at, last_login 
	          FROM users WHERE username = ? AND is_active = TRUE`

	err = db.DB.QueryRow(query, loginReq.Username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.FullName, &user.Email,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)

	if err == sql.ErrNoRows {
		log.Printf("User not found or inactive: %s", loginReq.Username)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	} else if err != nil {
		log.Println("Error querying user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	log.Printf("User found: %s, verifying password", user.Username)
	log.Printf("Password hash from DB: %s", user.PasswordHash)
	log.Printf("Password provided: %s", loginReq.Password)

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		log.Printf("Password verification failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	log.Printf("Login successful for user: %s", user.Username)

	// Update last login time
	_, _ = db.DB.Exec("UPDATE users SET last_login = ? WHERE id = ?", time.Now(), user.ID)

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.LoginResponse{
		Success: true,
		Message: "Login successful",
		User:    &user,
		Token:   "dummy-token", // In production, generate a proper JWT token
	})
}
