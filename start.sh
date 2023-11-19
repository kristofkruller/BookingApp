#!/bin/bash

# Function to check if a port is open
is_port_open() {
    nc -z localhost $1
}

# Function to start Docker services
start_docker_services() {
    echo "Starting Docker services..."
    docker-compose up -d
    if [ $? -ne 0 ]; then
        echo "Failed to start Docker services."
        exit 1
    fi
}

# Function to wait for services to be ready
wait_for_services() {
    local retries=10
    while [ $retries -gt 0 ]; do
        if is_port_open 5432 && is_port_open 8081 && is_port_open 8082 && is_port_open 8083; then
            echo "All services started successfully."
            return 0
        else
            echo "Waiting for services to start... Retries left: $retries"
            sleep 3
            ((retries--))
        fi
    done

    echo "Services did not start in time."
    return 1
}

# Function to run db-seeder
run_db_seeder() {
    echo "Running db-seeder..."
    go run ./db-seeder/main.go
    if [ $? -ne 0 ]; then
        echo "Failed to run db-seeder."
        exit 1
    fi
}

# Main script execution
start_docker_services
wait_for_services
if [ $? -eq 0 ]; then
    run_db_seeder
else
    echo "Aborting due to services not starting correctly."
    docker-compose -f docker-compose.yml down -v --remove-orphans
    exit 1
fi