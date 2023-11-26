# Example Dockerfile for booking-service
FROM golang:1.21.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first and download dependencies (for better caching)
COPY ./booking-service/go.mod ./booking-service/go.sum ./booking-service/
RUN cd /app/booking-service && go mod download && go mod verify

# Copy the local package files to the container's workspace
COPY ./libs /app/libs
COPY ./booking-service /app/booking-service

# Build your application
RUN cd /app/booking-service && CGO_ENABLED=0 GOOS=linux go build -v -o main .

# Use a smaller base image to create a final slim image
FROM alpine:latest  
COPY --from=builder /app/booking-service/main /app/main

# Run the application
CMD ["/app/main"]
