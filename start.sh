#!/bin/bash

# Function to check if a port is open
is_port_open() {
  nc -z 127.0.0.1 $1 &>/dev/null
  local result=$?
  if [ $result -eq 0 ]; then
      echo "Port $port is open."
  else
      echo "Port $port is not open."
  fi
  return $result
}

# Function to start Docker services
start_docker_services() {
    echo "Starting Docker services..."
    docker-compose up -d
    if [ $? -ne 0 ]; then
        echo "Failed to start Docker services."
        revert_go_mod_changes
        exit 1
    fi
}

# Function to update go.mod files for Docker builds
update_go_mod_for_docker() {
    echo "Updating go.mod files for Docker builds..."
    local services=("auth-service" "check-service" "booking-service" "payment-service")
    for service in "${services[@]}"; do
        local go_mod_path="./${service}/go.mod"
        if [ -f "$go_mod_path" ]; then
            # Comment out the replace directive
            sed -i.bak 's/^replace/\/\/replace/' "$go_mod_path"
            echo "Updated go.mod for $service"
        else
            echo "go.mod not found for $service"
        fi
    done
}

# Function to wait for services to be ready
wait_for_services() {
    local retries=3
    while [ $retries -gt 0 ]; do
        if is_port_open 5432 && is_port_open 8081 && is_port_open 8082 && is_port_open 8083; then
            echo "All services started successfully."
            return 0
        else
            echo "Waiting for services to start... Retries left: $retries"
            echo "Port 5432: $(is_port_open 5432)"
            echo "Port 8081: $(is_port_open 8081)"
            echo "Port 8082: $(is_port_open 8082)"
            echo "Port 8083: $(is_port_open 8083)"
            sleep 2
            ((retries--))
        fi
    done

    echo "Services did not start in time."
    return 1
}

# Function to revert go.mod files to original state
revert_go_mod_changes() {
    echo "Reverting go.mod files to original state..."
    local services=("auth-service" "check-service" "booking-service" "payment-service")
    for service in "${services[@]}"; do
        local go_mod_path="./${service}/go.mod"
        local backup_path="${go_mod_path}.bak"
        if [ -f "$backup_path" ]; then
            mv "$backup_path" "$go_mod_path"
            echo "Reverted go.mod for $service"
        else
            echo "Backup go.mod not found for $service"
        fi
    done
}

# Function to run db-seeder
run_db_seeder() {
    echo "Running db-seeder..."
    go run ./db-seeder/main.go
    if [ $? -ne 0 ]; then
        echo "Failed to run db-seeder."
        revert_go_mod_changes
        exit 1
    fi
}

# Main script execution
update_go_mod_for_docker
start_docker_services
wait_for_services
if [ $? -eq 0 ]; then
    run_db_seeder
    revert_go_mod_changes
else
    echo "Aborting due to services not starting correctly."
    docker-compose -f docker-compose.yml down -v --remove-orphans
    revert_go_mod_changes
    exit 1
fi