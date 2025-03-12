#!/usr/bin/env bash

set -Eeuo pipefail

docker network rm pyramid-network || echo "Network does not exist"
