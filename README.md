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

*An initialization script (root/init.sql) is used to create these tables. A db-seeder service seeds the users table with initial data using bcrypt for password hashing more details at ###Seeding below*

## Getting Started
Just fire the ./start.sh from the project root

### Prerequisites
- Docker compose v3.8
- go1.21.4 linux/amd64
- netcat for linux distro

### Installation
1. `git clone` [https://github.com/kristofkruller/BookingApp.git]
2. `cd BookingApp`

### Ports _links for dev_
- `auth-service`: [8081](http://127.0.0.1:8081)
- `check-service`: [8082](http://127.0.0.1:8082)
- `booking-service`: [8083](http://127.0.0.1:8083)
- `db` is on default 5432

### Development Env and manual start
The project is set up for development with VS Code through WSL Debian. A launch.json file is included for debugging:
- `Run and Debug - Ctrl+Shift+D` then you can start all services separately without containerized environment.
- Run `docker-compose -f docker-compose.yml up db -d` this will set up the db as a separate container but without the other services. You should seed it with `go run ./db-seeder/main.go`

### Env
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

### Seeding
The db-seeder service runs automatically during **start.sh** and seeds the users table.
The admin user is seeded with a bcrypt-hashed password for enhanced security.

### Notes on Security
JWT secret keys and database passwords are managed via environment variables for security.
Use HTTPS in production to ensure secure communication.