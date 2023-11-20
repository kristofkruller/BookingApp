package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kristofkruller/BookingApp/auth-service/config"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetUserByName(n string) (*config.User, error) {
	if db == nil {
		return nil, errors.New("database connection is not initialized")
	}

	u := &config.User{}

	// QRY
	err := db.QueryRow(`
			SELECT id, name, email, password 
			FROM users 
			WHERE name = $1
		`, n).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req config.LoginRequest

	// DECODE
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// VALIDATE
	isValid, err := validateUserPass(req.Name, req.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			// Database or other
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	// VALID
	if isValid {
		fmt.Println("login ok, starting token creation")
		u, err := GetUserByName(req.Name)
		if err != nil {
			http.Error(w, "Error extracting user name at token creation", http.StatusInternalServerError)
			return
		}

		// JWT session
		tokenString, err := config.GenerateJWT(u.ID)
		if err != nil {
			http.Error(w, "Error at token creation", http.StatusInternalServerError)
			return
		}

		// RESPONSE Set the token as an HTTP-only cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  time.Now().Add(1 * time.Hour), // Set the same expiration as your JWT token
			HttpOnly: true,
		})
		fmt.Fprint(w, "Login Successful")
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// CLEAR the JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), // time in the past
		HttpOnly: true,
	})

	fmt.Fprint(w, "Logout Successful")
}

// HELPERS
func validateUserPass(n, p string) (bool, error) {
	u, err := GetUserByName(n)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	return err == nil, nil
}
