#!/usr/bin/env bash
# Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0


set -Eeuo pipefail

# Make sure .env file exists in deploy/scripts/ directory
./deployments/scripts/backend/env_setup_backend.sh

docker compose -f ./deployments/docker-compose/backend/docker-compose.node.yml down
