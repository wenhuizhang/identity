// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/issuer/badge/data/filesystem"
	"github.com/agntcy/identity/internal/issuer/badge/mcp"
	issfs "github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	mdfs "github.com/agntcy/identity/internal/issuer/metadata/data/filesystem"
	"github.com/agntcy/identity/internal/pkg/nodeapi"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/spf13/cobra"
)

var (
	// setup the command flags
	issueMcpServerUrl  string
	issueMcpServerName string
)

var IssueMcpServerCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Issue a badge based on an MCP server URL",
	Run: func(cmd *cobra.Command, args []string) {

		// setup the badge service
		badgeFilesystemRepository := filesystem.NewBadgeFilesystemRepository()
		issuerRepository := issfs.NewIssuerFilesystemRepository()
		mdRepository := mdfs.NewMetadataFilesystemRepository()
		oidcAuth := oidc.NewAuthenticator()
		nodeClientPrv := nodeapi.NewNodeClientProvider()
		badgeService := badge.NewBadgeService(badgeFilesystemRepository, mdRepository, issuerRepository, oidcAuth, nodeClientPrv)

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForBadge()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the mcp server url is not set, prompt the user for it interactively
		if issueMcpServerUrl == "" {
			fmt.Fprintf(os.Stderr, "URL of the MCP server you want to sign in the badge: \n")
			_, err := fmt.Scanln(&issueMcpServerUrl)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading mcp server URL: %v\n", err)
				return
			}
		}
		if issueMcpServerUrl == "" {
			fmt.Fprintf(os.Stderr, "No MCP server URL provided\n")
			return
		}

		// if the mcp server name is not set, prompt the user for it interactively
		if issueMcpServerName == "" {
			fmt.Fprintf(os.Stderr, "Name of the MCP server you want to sign in the badge: \n")
			_, err := fmt.Scanln(&issueMcpServerName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading mcp server name: %v\n", err)
				return
			}
		}
		if issueMcpServerName == "" {
			fmt.Fprintf(os.Stderr, "No MCP server name provided\n")
			return
		}

		// Retrieve the MCP server data
		context := context.Background()
		mcpClient := mcp.NewDiscoveryClient()
		mcpServer, err := mcpClient.Discover(context, issueMcpServerName, issueMcpServerUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error discovering MCP server: %v\n", err)
			return
		}
		if mcpServer == nil {
			fmt.Fprintf(os.Stderr, "No MCP server found\n")
			return
		}

		// Marshal the MCP server to JSON
		mcpServerData, err := json.Marshal(mcpServer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling MCP server: %v\n", err)
			return
		}

		badgeContent := string(mcpServerData)

		// TODO
		badgeId, err := badgeService.IssueBadge(
			cache.VaultId,
			cache.IssuerId,
			cache.MetadataId,
			&vctypes.CredentialContent[vctypes.BadgeClaims]{
				Type:    vctypes.CREDENTIAL_CONTENT_TYPE_AGENT_BADGE,
				Content: vctypes.BadgeClaims{Badge: badgeContent},
			},
			nil,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error issuing badge: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Issued badge with ID: %s\n", badgeId)

		// Save the badge ID to the cache
		cache.BadgeId = badgeId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
	},
}

func init() {
	IssueMcpServerCmd.Flags().StringVarP(&issueMcpServerUrl, "url", "u", "", "The URL of the MCP server")
	IssueMcpServerCmd.Flags().StringVarP(&issueMcpServerName, "name", "n", "", "The name of the MCP server")
}
