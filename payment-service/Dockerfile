# Example Dockerfile for payment-service
FROM golang:1.21.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first and download dependencies (for better caching)
COPY ./payment-service/go.mod ./payment-service/go.sum ./payment-service/
RUN cd /app/payment-service && go mod download && go mod verify

# Copy the local package files to the container's workspace
COPY ./libs /app/libs
COPY ./payment-service /app/payment-service

# Build your application
RUN cd /app/payment-service && CGO_ENABLED=0 GOOS=linux go build -v -o main .

# Use a smaller base image to create a final slim image
FROM alpine:latest  
COPY --from=builder /app/payment-service/main /app/main

# Run the application
CMD ["/app/main"]
