#!/bin/bash
# Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0


DOCKER_FILE=./deployments/docker/identity/Dockerfile.test
TEST_COMMAND='go test -coverprofile=/mnt/coverage.out ./...'

echo RUNNING TESTS
docker run -v "${PWD}":/mnt "$(docker build --no-cache -f ${DOCKER_FILE} -q .)" "$TEST_COMMAND"
