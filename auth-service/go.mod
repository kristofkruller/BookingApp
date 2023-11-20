module github.com/kristofkruller/BookingApp/auth-service

go 1.21.4

require (
	github.com/gorilla/mux v1.8.1
	github.com/kristofkruller/BookingApp/assets v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.15.0
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
)

replace github.com/kristofkruller/BookingApp/assets => ../assets
