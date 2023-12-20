module github.com/kristofkruller/BookingApp/db-seeder

go 1.21.4

require (
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.15.0
)

require github.com/joho/godotenv v1.5.1

replace github.com/kristofkruller/BookingApp/libs => ../libs
