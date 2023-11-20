FROM golang:1.21 as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

# Build the Go app.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app .
# Start a new stage from scratch for a smaller final image.
FROM alpine:latest
# Copy the pre-built binary file from the previous stage.
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

CMD ["app"]