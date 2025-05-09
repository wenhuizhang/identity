#!/bin/bash -e
# Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0

# Syntax build-docker.sh [-i|--image imagename]

PROJECT=identity
DOCKER_IMAGE=${PROJECT}:latest
DOCKER_FILE=Dockerfile

while [[ $# -gt 0 ]]; do
  key="${1}"

  case ${key} in
  -i | --image)
    DOCKER_IMAGE="${2}"
    shift
    shift
    ;;
  -h | --help)
    less README.md
    exit 0
    ;;
  *) # unknown
    echo Unknown Parameter "$1"
    exit 4
    ;;
  esac
done

case $DOCKER_IMAGE in
$PROJECT-node*)
  DOCKER_FILE=./deployments/docker/identity/node/Dockerfile.node
  ;;

esac

echo BUILDING DOCKER "${DOCKER_IMAGE}"
docker build -t "${DOCKER_IMAGE}" --build-arg ARTIFACTORY_TOKEN="${ARTIFACTORY_TOKEN}" -f "${DOCKER_FILE}" .
