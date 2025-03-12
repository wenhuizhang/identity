#!/usr/bin/env bash

set -Eeuo pipefail

docker compose -f ./deployments/docker-compose/docs/docker-compose.yml down
