#!/usr/bin/env bash

set -Eeuo pipefail

docker network create identity-network || echo "Network already exists"
