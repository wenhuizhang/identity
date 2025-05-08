#!/usr/bin/env bash
# Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0


set -Eeuo pipefail

docker network create identity-network || echo "Network already exists"
