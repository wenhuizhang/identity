# Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0

# Docker build args
variable "IMAGE_REPO" {default = ""}
variable "IMAGE_TAG" {default = "v0.0.0"}

function "get_tag" {
  params = [tags, name]
  result = [for tag in tags: "${IMAGE_REPO}/${name}:${tag}"]
}

group "default" {
  targets = [
    "node",
  ]
}

target "docker-metadata-action" {
  tags = []
}

target "_common" {
  output = [
    "type=image",
  ]
  platforms = [
    "linux/amd64",
  ]
}

target "node" {
  context = "."
  dockerfile = "./deployments/docker/identity/Dockerfile.node"
  inherits = [
    "_common",
    "docker-metadata-action",
  ]
  tags = get_tag(target.docker-metadata-action.tags, "${target.node.name}")
}
