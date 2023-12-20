module github.com/kristofkruller/BookingApp/check-service

go 1.21.4

require (
	github.com/gorilla/mux v1.8.1
	github.com/jackc/pgtype v1.14.0
	github.com/kristofkruller/BookingApp/libs v0.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace github.com/kristofkruller/BookingApp/libs => ../libs
