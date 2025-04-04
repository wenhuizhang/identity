#!/usr/bin/env bash

set -Eeuo pipefail

docker network rm identity-network || echo "Network does not exist"
