// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/issuer/badge/a2a"
	"github.com/agntcy/identity/internal/issuer/badge/mcp"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

func NewCmd(
	cache *clicache.Cache,
	badgeService badge.BadgeService,
	vaultSrv vault.VaultService,
	a2aClient a2a.DiscoveryClient,
	mcpClient mcp.DiscoveryClient,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "badge",
		Short: "Issue and publish badges for your Agent and MCP Server identities",
		Long: `
The badge command is used to issue and publish badges for your Agent and MCP Server identities.
`,
	}

	cmd.AddCommand(NewCmdIssue(cache, badgeService, vaultSrv, a2aClient, mcpClient))
	cmd.AddCommand(NewCmdPublish(cache, badgeService))
	cmd.AddCommand(NewCmdList(cache, badgeService))
	cmd.AddCommand(NewCmdShow(cache, badgeService))
	cmd.AddCommand(NewCmdLoad(cache, badgeService))

	return cmd
}
