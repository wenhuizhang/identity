#!/bin/bash
# Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0

NODE_ENV=./identity/cmd/node/.env
IDENTITY_DEPLOYMENTS_DIR=./deployments/docker-compose/identity

# This script sets up the environment for the Identity Node service.
# It checks for the existence of a .env file in the identity directory
# and creates one with default values if it doesn't exist.
if [ ! -f "$NODE_ENV" ]; then
  echo ".env File not found in the Identity directory, using defaults"

  echo "Creating .env file with defaults"
  touch "$NODE_ENV" && \
  echo "DB_HOST=identity-postgres" > "$NODE_ENV" && \
  echo "DB_PORT=5432" >> "$NODE_ENV" && \
  echo "DB_USERNAME=postgres" >> "$NODE_ENV" && \
  echo "DB_PASSWORD=postgres" >> "$NODE_ENV" && \
  echo "POSTGRES_PASSWORD=postgres" >> "$NODE_ENV" && \
  echo "POSTGRES_DB=identity" >> "$NODE_ENV"
fi

# Check if the .env file exists in the deployments/docker-compose/identity directory
# If not, create a symlink to the Identity directory .env file
if [ ! -f "$IDENTITY_DEPLOYMENTS_DIR/.env" ]; then
  echo ".env File not found, creating symlink to Identity directory .env file"
  cd "$IDENTITY_DEPLOYMENTS_DIR" && ln -s "../../../$NODE_ENV" .
fi
