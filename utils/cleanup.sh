#!/bin/bash
# HANDLE WITH CARE!
docker-compose down
docker container prune
docker volume prune
docker network prune
echo "°CLEANUP DONE°"