package reserv

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/kristofkruller/BookingApp/booking-service/pb"
	"github.com/kristofkruller/BookingApp/libs/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *GRPCServer) LetsBook(ctx context.Context, req *pb.BookingReq) (*pb.BookingsRes, error) {
	st := time.Now()

	// Convert protobuf request to internal BookingReq type
	br := convertProtoToBookingReq(req)

	// Validate request
	if err := br.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// TIME HANDL Bookingperiod
	start, end, err := parseBookingPeriod(br)
	if err != nil {
		return nil, err
	}

	bp := &pgtype.Tsrange{
		Lower:     pgtype.Timestamp{Time: start, Status: pgtype.Present},
		Upper:     pgtype.Timestamp{Time: end, Status: pgtype.Present},
		LowerType: pgtype.Inclusive,
		UpperType: pgtype.Exclusive,
	}

	// Check room availability
	if bp != nil && !helpers.IsRoomAvailable(br.RoomID, *bp, s.Db) {
		return nil, status.Errorf(codes.FailedPrecondition, "Room is not available for the selected dates")
	}

	// INSERT INTO DB
	if err := createBookingInDB(br, s.Db); err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating booking: %v", err)
	}

	// Fetch all bookings for user for validation
	bookings, err := fetchUserBookings(br.UserID, s.Db)
	if err != nil {
		return nil, err
	}

	log.Printf("Booking created successfully in %v", time.Since(st))

	// Convert internal bookings to protobuf response
	res := convertBookingsToProto(bookings)
	return res, nil
}
