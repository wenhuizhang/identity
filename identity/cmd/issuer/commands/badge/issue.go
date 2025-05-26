// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"github.com/spf13/cobra"
)

var badgeIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue badges using different data sources",
	Long: `
The issue command is used to create Badges for your Agent and MCP Server identities from various data sources.
`,
}
