#!/usr/bin/env bash

set -Eeuo pipefail

docker compose -f ./deployments/docker-compose/docs/docker-compose.yml build --no-cache
docker compose -f ./deployments/docker-compose/docs/docker-compose.yml up -d
