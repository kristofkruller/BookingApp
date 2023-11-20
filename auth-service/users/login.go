package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kristofkruller/BookingApp/assets"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUserPass(n, p string) (bool, error) {
	u, err := GetUserByName(n)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	return err == nil, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req assets.LoginRequest

	// Decode the JSON body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the user's password
	isValid, err := ValidateUserPass(req.Name, req.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			// Database or other server errors
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	if isValid {
		fmt.Println("login ok, starting token creation")
		u, err := GetUserByName(req.Name)
		if err != nil {
			http.Error(w, "Error extracting user name at token creation", http.StatusInternalServerError)
			return
		}

		// JWT session
		tokenString, err := assets.GenerateJWT(u.ID)
		if err != nil {
			http.Error(w, "Error at token creation", http.StatusInternalServerError)
			return
		}

		// Set the token as an HTTP-only cookie
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

// clear the JWT cookie
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), // time in the past
		HttpOnly: true,
	})

	fmt.Fprint(w, "Logout Successful")
}
