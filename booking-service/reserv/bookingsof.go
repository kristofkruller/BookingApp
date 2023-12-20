package reserv

import (
	"context"
	"log"

	"github.com/kristofkruller/BookingApp/booking-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) BookingsOf(ctx context.Context, req *pb.BookingsOfReq) (*pb.BookingsRes, error) {
	// Validate user ID
	if req.UId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID")
	}

	// Construct SQL query with filters
	q, args, err := constructBookingsQuery(req)
	if err != nil {
		return nil, err
	}

	// Execute the query
	rows, err := s.Db.QueryContext(ctx, q, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, status.Errorf(codes.Internal, "Error executing query: %v", err)
	}
	defer rows.Close()

	// Process rows and build response
	bookings, err := processBookingRows(rows)
	if err != nil {
		return nil, err
	}

	// Convert internal bookings to protobuf response
	res := convertBookingsToProto(bookings)
	return res, nil
}
