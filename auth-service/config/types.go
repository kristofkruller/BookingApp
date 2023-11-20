package config

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
