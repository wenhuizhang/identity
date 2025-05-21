# Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.24.1-alpine3.20 AS builder

ARG BIN="/usr/local/bin"
ARG BUF_VERSION="1.28.1"

WORKDIR /identity

COPY ./scripts/node/docker/entrypoint.sh .
COPY ./scripts/node/docker/templates ./templates

RUN apk update && apk add curl git

RUN curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m)" \
    -o "${BIN}/buf" && \
    chmod +x "${BIN}/buf"

RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN mkdir go-swagger
RUN git clone https://github.com/go-swagger/go-swagger "./go-swagger"
WORKDIR /identity/go-swagger
RUN go install ./cmd/swagger

WORKDIR /identity

RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
