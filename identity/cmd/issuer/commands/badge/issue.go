// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/cmd/issuer/commands/badge/issue"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/issuer/badge/a2a"
	"github.com/agntcy/identity/internal/issuer/badge/mcp"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

func NewCmdIssue(
	cache *clicache.Cache,
	badgeService badge.BadgeService,
	vaultSrv vault.VaultService,
	a2aClient a2a.DiscoveryClient,
	mcpClient mcp.DiscoveryClient,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue badges using different data sources",
		Long: `
The issue command is used to create Badges for your Agent and MCP Server identities from various data sources.
`,
	}

	cmd.AddCommand(issue.NewCmdIssueA2A(
		cache,
		badgeService,
		vaultSrv,
		a2aClient,
	))
	cmd.AddCommand(issue.NewCmdMcp(
		cache,
		badgeService,
		vaultSrv,
		mcpClient,
	))
	cmd.AddCommand(issue.NewCmdIssueOasf(
		cache,
		badgeService,
		vaultSrv,
	))

	return cmd
}
