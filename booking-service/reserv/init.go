package reserv

import (
	"database/sql"

	"github.com/kristofkruller/BookingApp/booking-service/pb"
)

func SetDB(sdb *sql.DB) {
	db = sdb
}

// GRPCServer struct implements the gRPC server interface
type GRPCServer struct {
	Db *sql.DB
	pb.UnimplementedBookingServiceServer
}
