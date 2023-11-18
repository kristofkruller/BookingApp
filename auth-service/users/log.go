package users

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

func hashPass(p string) string {
	hasher := sha256.New()
	hasher.Write([]byte(p))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ValidateUserPass(n, p string) (bool, error) {
	u, err := GetUserByName(n)
	if err != nil {
		return false, err
	}

	hashedPass := hashPass(p)
	return u.Password == hashedPass, nil
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
		http.Error(w, "Error validating user pass", http.StatusInternalServerError)
		return
	}

	if isValid {
		fmt.Fprintf(w, "Login Successful")
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Invalidate the user's session or token here (not implemented)

	fmt.Fprintf(w, "Logout Successful")
}
