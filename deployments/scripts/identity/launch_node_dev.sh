#!/usr/bin/env bash
# Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0


set -Eeuo pipefail

# If .env exists in the node directory, use it
# If not, create an env with defaults
./deployments/scripts/identity/env_setup_node.sh

docker compose -f ./deployments/docker-compose/identity/docker-compose.node.dev.yml build --no-cache
docker compose -f ./deployments/docker-compose/identity/docker-compose.node.dev.yml up -d
