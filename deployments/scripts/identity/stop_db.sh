#!/usr/bin/env bash
# Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0


set -Eeuo pipefail

# Make sure .env file exists in deploy/scripts/ directory
./deployments/scripts/identity/env_setup_node.sh

docker compose -f ./deployments/docker-compose/identity/docker-compose.couchdb.yml down
