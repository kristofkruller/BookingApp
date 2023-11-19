package users

import (
	"database/sql"

	"github.com/kristofkruller/BookingApp/auth-service/config"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetUserByName(n string) (*config.User, error) {
	u := &config.User{}
	err := db.QueryRow("SELECT id, name, email, password FROM users WHERE name = $1", n).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}