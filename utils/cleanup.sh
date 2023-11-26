#!/bin/bash
# HANDLE WITH CARE!
docker-compose down
docker container prune
docker volume prune
docker network prune

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
revert_go_mod_changes
echo "°CLEANUP DONE°"