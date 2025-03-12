#!/usr/bin/env bash

set -Eeuo pipefail

docker network create pyramid-network || echo "Network already exists"
