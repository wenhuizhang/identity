// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	mcptypes "github.com/agntcy/identity/internal/core/mcp/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
)

// The discoverClient interface defines the core methods for
// discovering a deployed MCP server
type DiscoveryClient interface {
	Discover(ctx context.Context, name, url string) (*mcptypes.McpServer, error)
}

// The discoverClient struct implements the DiscoverClient interface
type discoveryClient struct {
}

// NewDiscoverClient creates a new instance of the DiscoverClient
func NewDiscoveryClient() DiscoveryClient {
	return &discoveryClient{}
}

func (d *discoveryClient) Discover(
	ctx context.Context,
	name, url string,
) (*mcptypes.McpServer, error) {
	// Create streameable http client
	// We only support streamable http client for now
	transportClient, err := transport.NewStreamableHTTPClientTransport(fmt.Sprintf("%s/mcp", url))
	if err != nil {
		return nil, errutil.Err(
			err,
			"failed to create streamable http client transport",
		)
	}

	// Initialize MCP client
	mcpClient, err := client.NewClient(
		transportClient,
	)
	if err != nil {
		return nil, errutil.Err(
			err,
			"failed to create mcp client",
		)
	}
	defer mcpClient.Close()

	// Discover MCP server
	// First the tools
	toolsList, err := mcpClient.ListTools(ctx)
	if err != nil {
		return nil, errutil.Err(
			err,
			"failed to discover mcp server",
		)
	}

	// After that the resources
	resourcesList, err := mcpClient.ListResources(ctx)
	if err != nil {
		return nil, errutil.Err(
			err,
			"failed to discover mcp server",
		)
	}

	// Parse the tools and resources
	// Get the first batch of tools
	availableTools := make([]*mcptypes.McpTool, 0)

	for _, tool := range toolsList.Tools {
		// Convert parameters to JSON string
		jsonParams, err := json.Marshal(tool.InputSchema)
		if err != nil {
			jsonParams = []byte{}
		}

		availableTools = append(availableTools, &mcptypes.McpTool{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  string(jsonParams),
		})
	}

	// Get the first batch of resources
	availableResources := make([]*mcptypes.McpResource, 0)

	for _, resource := range resourcesList.Resources {
		availableResources = append(availableResources, &mcptypes.McpResource{
			Name:        resource.Name,
			Description: resource.Description,
			URI:         resource.URI,
		})
	}

	return &mcptypes.McpServer{
		Name:      name,
		URL:       url,
		Tools:     availableTools,
		Resources: availableResources,
	}, nil
}
