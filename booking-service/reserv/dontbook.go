package reserv

import (
	"context"
	"log"

	"github.com/kristofkruller/BookingApp/booking-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) DontBook(ctx context.Context, req *pb.DontReq) (*pb.DontRes, error) {
	// Log the incoming request
	log.Printf("Canceling booking with ID: %d", req.BookingId)

	// Execute the SQL query to delete the booking
	res, err := s.Db.Exec("DELETE FROM reserv WHERE id = $1", req.BookingId)
	if err != nil {
		// Handle any error that occurred during query execution
		log.Printf("Error canceling booking: %v", err)
		return nil, status.Errorf(codes.Internal, "Error canceling booking: %v", err)
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return nil, status.Errorf(codes.Internal, "Error during cancellation: %v", err)
	}

	if rowsAffected == 0 {
		// No rows affected means no booking was found with the given ID
		return nil, status.Errorf(codes.NotFound, "No booking found with the given ID")
	}

	// Return a successful response
	return &pb.DontRes{Cancelled: "Booking canceled successfully"}, nil
}
