#!/usr/bin/env bash
# Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0

# If .env exists in the node directory, use it
# If not, create an env with defaults
./deployments/scripts/identity/env_setup_node.sh

# Check if dev option is set
compose_file="./deployments/docker-compose/identity/docker-compose.node.yml"
if [ "$1" == "true" ]; then
    echo "Running in dev mode"
    compose_file="./deployments/docker-compose/identity/docker-compose.node.dev.yml"
fi

docker compose -f "$compose_file" build --no-cache
docker compose -f "$compose_file" up -d
