package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kristofkruller/BookingApp/auth-service/config"
	"golang.org/x/crypto/bcrypt"
)

func hashPass(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ValidateUserPass(n, p string) (bool, error) {
	u, err := GetUserByName(n)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	return err == nil, nil
}

// LoginRequest represents the JSON structure for a login request
type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Decode the JSON body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the user's password
	isValid, err := ValidateUserPass(req.Name, req.Password)
	if err != nil {
		http.Error(w, "Error validating user pass", http.StatusUnauthorized)
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
		tokenString, err := config.GenerateJWT(u.ID)
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
