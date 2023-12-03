module github.com/kristofkruller/BookingApp/booking-service

go 1.21.4

require (
	github.com/gorilla/mux v1.8.1
	github.com/jackc/pgtype v1.14.0
	github.com/kristofkruller/BookingApp/libs v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
)

replace github.com/kristofkruller/BookingApp/libs => ../libs
