# BookingApp
A simple accommodation booking system

## Overview
The system is divided into the following microservices:

- `auth-service`: *Manages user authentication.*
- `check-service`: *Handles hotels and rooms management.*
- `booking-service`: *Manages booking and cancellation of rooms.*

The PostgreSQL database is set up with the following tables:

- `users`: *Contains user information.*
- `properties`: *Contains details about properties like hotels and flats.*
- `rooms`: *Details about individual rooms linked to properties.*
- `reserv`: *Reservation handling*

*An initialization script (root/init.sql) is used to create these tables. A db-seeder service seeds the users table with initial data using bcrypt for password hashing more details at ###Seeding below*

## Getting Started
### Installation
1. `git clone` [https://github.com/kristofkruller/BookingApp.git]
2. `cd BookingApp`
3. start:
`a)` Just fire the ./start.sh from the project root 
**_OR_** 
`b)` docker-compose up --build and if done go run ./db-seeder/main.go
**_important_** that either way wait for both
docker compose will log to your terminal as access and error log
```
Admin user seeded successfully.
Reservations seeded successfully.
```
CLI feedbacks

### Prerequisites
- Docker compose v3.8
- go1.21.4 linux/amd64
- netcat for linux distro

### Ports _links for dev_
- `auth-service`: [8081](http://127.0.0.1:8081)
- `check-service`: [8082](http://127.0.0.1:8082)
- `booking-service`: [8083](http://127.0.0.1:8083)
- `db` is on default 5432

## API list
### provided API json at ./utils/ 
If you post a request set the body type to JSON and include desired content I use [Insomnia API](https://insomnia.rest/)
_I assume a formatted string from a date picker by date values_
For requests with filtering it is possible to use only a "partly" filter, but with logical pair i.e.:
```
{
  "availability_start": "2023-01-01",
  "availability_end": "2023-01-10"
}
_or_
{
  "price_min": 60,
  "price_max": 70
}
_or_
{
  "creation_date": "2023-11-20"
}
_or combined even_
{
  "creation_date": "2023-11-20",
  "start_date": "2023-02-01",
  "end_date": "2023-02-05"
}
etc.
```
**if you want a filter-free list then post and empty JSON object `{}` with the REQ to get the full array as RES**
float in this case below always a DECIMAL(10,2)

### Auth
- [:8081/check] GET expects nothing
_Response text: Auth-service up_
- [:8081/login] POST expects:
```
REQ Body Params (JSON format):
{
  "name":"string",
  "password":"string"
}
```
_Response feedback msg, http only cookie with jwt token exp. 1hr_
- [:8081/logout] POST invalidate by setting cookie to be expired
_Response feedback msg, cookie data, and empty token val_

### Check
- [:8081/check] GET expects nothing
_Response text: Check-service up_
- [:8082/room/{id}] GET expects an int for ID
_Response a JSON object with all data of the selected room_
- [:8082/rooms] POST with **optional** filter params:
```
REQ Body Params (JSON format):
{
  price_min: float,
  price_max: float,
  availability_start: "YYYY-MM-DD",
  availability_end: "YYYY-MM-DD"
}
```
_Response a JSON list(array of objects) of rooms matching the filters_

### Booking
- [:8083/check] GET expects nothing
_Response text: Booking-service up_
- [:8083/bookingsof/{uId}] POST with **optional** filter params:
```
REQ Body Params (JSON format):
{
  "min_price": float,
  "max_price": float,
  "creation_date": "YYYY-MM-DD",
  "start_date": "YYYY-MM-DD",
  "end_date": "YYYY-MM-DD"
}
```
_Response a JSON list(array of objects) of bookings matching the filters_
- [:8083/letsbook] POST with **mandatory** params:
```
REQ Body Params (JSON format):
{
  "userId": int,
  "propertyId": int,
  "roomId": int,
  "cost": float,
  "start_date": "YYYY-MM-DD",
  "end_date": "YYYY-MM-DD"
}
```
_Response text: Booking created successfully_
- [:8083/dontbook/{bookingId}] POST to delete a booking by id
_Response text: Booking canceled successfully_

### Payment
- [:8084/check] GET expects nothing
_Response text: Payment-service up_
- [:8084/pay/{bookingId}] POST to pay booking by id with **mandatory** params:
```
REQ Body Params (JSON format):
{
  "bookingId": int,
  "amount": float,
  "currency": string,
  "cardToken": string
}
```
_Response JSON object about success or failed_


## Details, mechanics
### Seeding
The db-seeder service runs automatically during **start.sh** and seeds the users and reserv tables.
The admin user is seeded with a bcrypt-hashed password for enhanced security.
It is functioning as a go "script".

### Env
**MUST BE CREATED AT PROJECT _ROOT ._**
_Exposed env content for development_
*.env is not commited because of best practice*
```
DB_PASSWORD=asdf1234
POSTGRES_DB=BookingAppDb
POSTGRES_USER=admin
POSTGRES_PASSWORD=asdf1234
JWT_SECRET_KEY=Hu7ky4L1f3*
DB_USER=admin
DB_HOST=db
DB_NAME=BookingAppDb
DB_CONNECTION_SEED=postgres://admin:asdf1234@127.0.0.1/BookingAppDb?sslmode=disable
```

### Development Env and manual start
The project is set up for development with VS Code through WSL Debian. A launch.json file is included for debugging:
- `Run and Debug - Ctrl+Shift+D` then you can start all services separately without containerized environment.
- Run `docker-compose -f docker-compose.yml up db` this will set up the db as a separate container but without the other services. You should seed it with `go run ./db-seeder/main.go`
**there should be a local.env for launch.json, where DB_HOST=127.0.0.1** otherwise the connection will die with timeout.
This way of starting produces a brand new fresh binary to the out folder as well optimized for Linux environments.

## Notes on possible improvements
- Helper functions, types and code for general use must be regorganized to a lib, with functionality like in every `main.go` the program exits gracefully or time handlers. 
- Error handling and logging should be ogranized to a service or lib also health checkers for db, and endpoints
- Testcases
- For large datasets, consider indexing the reserv_interval column in the reserv table.
- Queries should be transferred into postgre as a function
- Frontend should be one GUI with an nginx reverse proxy channeled to :443