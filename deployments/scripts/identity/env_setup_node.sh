#!/bin/bash
# Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0


# This script sets up the environment for the Identity Node service.
# It checks for the existence of a .env file in the identity directory
# and creates one with default values if it doesn't exist.
if [ ! -f ./identity/cmd/node/.env ]; then
  echo ".env File not found in the Identity directory, using defaults"

  echo "Creating .env file with defaults"
  cd ../../../identity/cmd/node && \
  touch .env && \
  echo "DB_HOST=identity-postgres" > .env && \
  echo "DB_PORT=5432" >> .env && \
  echo "DB_USER=postgres" >> .env && \
  echo "DB_PASSWORD=postgres" >> .env
fi

# Check if the .env file exists in the deployments/docker-compose/identity directory
# If not, create a symlink to the Identity directory .env file
if [ ! -f ./deployments/docker-compose/identity/.env ]; then
  echo ".env File not found, creating symlink to Identity directory .env file"
  cd deployments/docker-compose/identity && ln -s ../../../identity/cmd/node/.env .
fi
