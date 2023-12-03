package reserv

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kristofkruller/BookingApp/booking-service/config"
	"github.com/kristofkruller/BookingApp/booking-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// HELPER FUNCS to convert between protobuf and internal types, parse dates, create bookings in DB, etc.
func convertBookingsToProto(bookings []config.Booking) *pb.BookingsRes {
	protoBookings := make([]*pb.Booking, 0, len(bookings))
	for _, b := range bookings {
		protoBooking := &pb.Booking{
			Id:         int32(b.ID),
			UserId:     int32(b.UserID),
			PropertyId: int32(b.PropertyID),
			RoomId:     int32(b.RoomID),
			Cost:       float32(b.Cost),
			ReservInterval: &pb.Tsrange{
				StartDate: timestamppb.New(b.ReservInterval.Lower.Time),
				EndDate:   timestamppb.New(b.ReservInterval.Upper.Time),
			},
			CreationDate: timestamppb.New(b.CreationDate),
		}
		protoBookings = append(protoBookings, protoBooking)
	}
	return &pb.BookingsRes{Bookings: protoBookings}
}

func convertProtoToBookingReq(req *pb.BookingReq) config.BookingReq {
	return config.BookingReq{
		UserID:     int(req.UserId),
		PropertyID: int(req.PropertyId),
		RoomID:     int(req.RoomId),
		Cost:       req.Cost,
		StartDate:  req.StartDate.AsTime().Format("2006-01-02"),
		EndDate:    req.EndDate.AsTime().Format("2006-01-02"),
	}
}

func parseBookingPeriod(br config.BookingReq) (time.Time, time.Time, error) {
	start, err := time.Parse("2006-01-02", br.StartDate)
	if err != nil {
		return time.Time{}, time.Time{}, status.Errorf(codes.InvalidArgument, "Invalid start date format")
	}

	end, err := time.Parse("2006-01-02", br.EndDate)
	if err != nil {
		return time.Time{}, time.Time{}, status.Errorf(codes.InvalidArgument, "Invalid end date format")
	}

	if !start.Before(end) {
		return time.Time{}, time.Time{}, status.Errorf(codes.InvalidArgument, "End date must be after start date")
	}

	return start, end, nil
}

func createBookingInDB(br config.BookingReq, db *sql.DB) error {
	_, err := db.Exec(`
		INSERT INTO reserv (userId, propertyId, roomId, cost, reserv_interval)
		VALUES ($1, $2, $3, $4, tsrange($5, $6, '[]'))`,
		br.UserID, br.PropertyID, br.RoomID, br.Cost, br.StartDate, br.EndDate,
	)
	return err
}

func constructBookingsQuery(req *pb.BookingsOfReq) (string, []interface{}, error) {
	var args []interface{}
	var conditions []string

	// Base query
	query := "SELECT id, userId, propertyId, roomId, cost, reserv_interval, creation_date FROM reserv WHERE userId = $1"
	args = append(args, req.UId)

	// Index for SQL arguments
	argIdx := 2

	// Check if the Filter field is present
	if req.Filter != nil {
		// Handle price filters
		if req.Filter.MinPrice != 0 {
			conditions = append(conditions, fmt.Sprintf("price >= $%d", argIdx))
			args = append(args, req.Filter.MinPrice)
			argIdx++
		}
		if req.Filter.MaxPrice != 0 {
			conditions = append(conditions, fmt.Sprintf("price <= $%d", argIdx))
			args = append(args, req.Filter.MaxPrice)
			argIdx++
		}

		// Handle date filters
		if req.Filter.CreationDate != "" {
			conditions = append(conditions, fmt.Sprintf("DATE(creation_date) = $%d", argIdx))
			args = append(args, req.Filter.CreationDate)
			argIdx++
		}
		if req.Filter.StartDate != "" && req.Filter.EndDate != "" {
			conditions = append(conditions, fmt.Sprintf("reserv_interval && tsrange($%d, $%d, '[]')", argIdx, argIdx+1))
			args = append(args, req.Filter.StartDate, req.Filter.EndDate)
			argIdx += 2
		}
	}

	// Append conditions to query if any
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	return query, args, nil
}

func fetchUserBookings(userID int, db *sql.DB) ([]config.Booking, error) {
	rows, err := db.Query(`
		SELECT id, userId, propertyId, roomId, cost, reserv_interval, creation_date 
		FROM reserv 
		WHERE userId = $1`, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error querying database: %v", err)
	}
	defer rows.Close()

	var bookings []config.Booking
	for rows.Next() {
		var b config.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.PropertyID, &b.RoomID, &b.Cost, &b.ReservInterval, &b.CreationDate); err != nil {
			return nil, status.Errorf(codes.Internal, "Error reading booking data: %v", err)
		}
		bookings = append(bookings, b)
	}
	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Error iterating booking rows: %v", err)
	}

	return bookings, nil
}

func processBookingRows(rows *sql.Rows) ([]config.Booking, error) {
	var bookings []config.Booking
	for rows.Next() {
		var b config.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.PropertyID, &b.RoomID, &b.Cost, &b.ReservInterval, &b.CreationDate); err != nil {
			return nil, status.Errorf(codes.Internal, "Error reading booking data: %v", err)
		}
		bookings = append(bookings, b)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Error iterating booking rows: %v", err)
	}
	return bookings, nil
}
