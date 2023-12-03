package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kristofkruller/BookingApp/booking-service/pb"
	"github.com/kristofkruller/BookingApp/booking-service/reserv"
	"github.com/kristofkruller/BookingApp/libs/initdb"
	"google.golang.org/grpc"
)

func main() {
	port := ":8083"
	// Initialize database with retry
	db, err := initdb.InitDb()
	if err != nil {
		log.Fatalf("could not connect to database after 3 attempts %v", err)
	}

	//DB conn for rooms pkg
	reserv.SetDB(db)

	// Listen on TCP port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new GRPCServer instance from the reserv package
	grpcServer := grpc.NewServer()
	reservServer := &reserv.GRPCServer{Db: db}
	pb.RegisterBookingServiceServer(grpcServer, reservServer)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("Server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	// Stop the server
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
