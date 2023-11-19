FROM golang:1.21.4-alpine

WORKDIR /usr/src/app

COPY . .

RUN go build -o seeder

CMD ["./seeder"]
