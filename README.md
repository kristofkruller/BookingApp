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

*An initialization script (root/init.sql) is used to create these tables and seed them with initial data.*

## Getting Started

### Prerequisites

- Docker compose v3.8
- go1.21.4 linux/amd64

### Installation
1. `git clone` [https://github.com/kristofkruller/BookingApp.git]
2. `cd BookingApp`
3. `docker-compose up`

### Ports
- `auth-service`: [8081](http://127.0.0.1:8081)
- `check-service`: [8082](http://127.0.0.1:8082)
- `booking-service`: [8083](http://127.0.0.1:8083)
