#!/usr/bin/env bash

set -Eeuo pipefail

# Make sure .env file exists in deploy/scripts/ directory
./deployments/scripts/backend/env_setup_backend.sh

docker compose -f ./deployments/docker-compose/backend/docker-compose.node.yml build --no-cache
docker compose -f ./deployments/docker-compose/backend/docker-compose.node.yml up -d
