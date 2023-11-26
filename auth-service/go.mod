module github.com/kristofkruller/BookingApp/auth-service

go 1.21.4

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/kristofkruller/BookingApp/libs v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.15.0
)

require github.com/lib/pq v1.10.9 // indirect

replace github.com/kristofkruller/BookingApp/libs => ../libs