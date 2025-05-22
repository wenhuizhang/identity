# Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0

# Docker build args
variable "IMAGE_REPO" {default = ""}
variable "IMAGE_TAG" {default = "v0.0.0"}

function "get_tag" {
  params = [tags, name]
  result = [for tag in tags: "${IMAGE_REPO}/${name}:${tag}"]
}

group "node" {
  targets = [
    "node",
  ]
}

target "_common" {
  output = [
    "type=image",
  ]
  platforms = [
    "linux/arm64",
    "linux/amd64",
  ]
}

target "node" {
  context = "../.."
  dockerfile = "./identity/Dockerfile.node"
  target = "node"
  inherits = [
    "_common"
  ]
  tags = get_tag(target.docker-metadata-action.tags, "${target.node.name}")
}
