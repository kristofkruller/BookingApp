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
important that either way wait for `Database seeding completed successfully` CLI feedback


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
If you post a request set the body type to JSON and include desired content I use [Insomnia API](https://insomnia.rest/)
### Auth
- [:8081/check] GET expects nothing
_Response text: Auth-service up_
- [:8081/login] POST expects Body Params (JSON format):
"name":"string"
"password":"string"
_Response feedback msg, http only cookie with jwt token exp. 1hr_
- [:8081/logout] POST set cookie by name to be expired
_Response feedback msg, cookie data, and empty token val_

### Check
- [:8081/check] GET expects nothing
_Response text: Check-service up_
- [:8082/room/{id}] GET expects an int for ID
_Response a JSON with all data of the selected room_
- [:8082/rooms] POST with **optional** filter params:
```
Body Params (JSON format):
price_min: DECIMAL(10,2) i.e.:50.00
price_max: DECIMAL(10,2) i.e.:50.00
_ava. I assume a formatted string from a date picker:_
availability_start: (format: YYYY-MM-DD) i.e.:"2023-01-01"
availability_end: (format: YYYY-MM-DD) i.e.:"2023-01-01"
```
it is possible to use only a partly filter, but with logical pair i.e.:
```
{
 	"availability_start": "2023-01-01",
  "availability_end": "2023-01-10"
}
_or_
{
  "price_min": 60.00,
  "price_max": 70.00
}
```
_Response a JSON list of rooms matching the filters_

### Booking
- [:8083/check] GET expects nothing
_Response text: Booking-service up_

## Details, mechanics
### Seeding
The db-seeder service runs automatically during **start.sh** and seeds the users table.
The admin user is seeded with a bcrypt-hashed password for enhanced security.
It is functioning as a go "script".

### Env
**MUST BE CREATED AT PROJECT _ROOT ._**
Exposed env content for development *env is not commited because of best practice*
DB_PASSWORD=asdf1234
POSTGRES_DB=BookingAppDb
POSTGRES_USER=admin
POSTGRES_PASSWORD=asdf1234
JWT_SECRET_KEY=Hu7ky4L1f3*
DB_USER=admin
DB_HOST=db
DB_NAME=BookingAppDb
DB_CONNECTION_SEED=postgres://admin:asdf1234@127.0.0.1/BookingAppDb?sslmode=disable

### Development Env and manual start
The project is set up for development with VS Code through WSL Debian. A launch.json file is included for debugging:
- `Run and Debug - Ctrl+Shift+D` then you can start all services separately without containerized environment.
- Run `docker-compose -f docker-compose.yml up db -d` this will set up the db as a separate container but without the other services. You should seed it with `go run ./db-seeder/main.go`

### Notes on Security
JWT secret keys and database passwords are managed via environment variables for security.
Use HTTPS in production to ensure secure communication.