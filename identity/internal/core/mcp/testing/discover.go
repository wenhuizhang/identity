// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"context"

	mcptypes "github.com/agntcy/identity/internal/core/mcp/types"
	"github.com/agntcy/identity/internal/issuer/mcp"
)

type FakeDiscoveryClient struct {
}

func NewFakeDiscoveryClient() mcp.DiscoveryClient {
	return &FakeDiscoveryClient{}
}

func (d *FakeDiscoveryClient) Discover(
	_ context.Context,
	_, _ string,
) (*mcptypes.McpServer, error) {
	return &mcptypes.McpServer{}, nil
}
