// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package mcp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"github.com/agntcy/identity/internal/issuer/mcp"
	"github.com/stretchr/testify/assert"
)

// McpServerURL is the URL of the MCP server
// This is used for testing purposes
const McpServerHost = "127.0.0.1:8000"

func TestShouldDiscoverADeployedServer(t *testing.T) {
	t.Parallel()

	// Create a new MCP test server
	go createMCPTestServer(t)

	// Create a new discovery client
	discoveryClient := mcp.NewDiscoveryClient()
	mcpServer, err := discoveryClient.Discover(
		context.Background(),
		"test-server",
		fmt.Sprintf("http://%s", McpServerHost),
	)

	t.Logf("MCP Server: %+v", mcpServer)

	assert.NoError(t, err)
}

func createMCPTestServer(t *testing.T) {
	// Create SSE transport server
	transportServer := transport.NewStreamableHTTPServerTransport(McpServerHost)

	// Initialize MCP server
	mcpServer, err := server.NewServer(transportServer)
	assert.NoError(t, err)

	// Start server
	err = mcpServer.Run()
	assert.NoError(t, err)
}
