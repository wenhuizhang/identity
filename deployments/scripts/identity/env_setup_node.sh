#!/bin/bash
# Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0


if [ ! -f ./backend/cmd/node/.env ]; then
  echo ".env File not found in the Backend directory, please create one"
  exit 1
fi

if [ ! -f ./deployments/docker-compose/backend/.env.node ]; then
  echo ".env File not found, creating symlink to Backend directory .env file"
  cd deployments/docker-compose/backend && ln -s ../../../backend/cmd/node/.env .env.node
fi
