#!/bin/bash
# Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0


if [ ! -f ./identity/cmd/node/.env ]; then
  echo ".env File not found in the Identity directory, please create one"
  exit 1
fi

if [ ! -f ./deployments/docker-compose/identity/.env ]; then
  echo ".env File not found, creating symlink to Identity directory .env file"
  cd deployments/docker-compose/identity && ln -s ../../../identity/cmd/node/.env .
fi
