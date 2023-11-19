package libs

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// LoginRequest represents the JSON structure for a login request
type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
