// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"github.com/agntcy/identity/cmd/issuer/commands/setup"
	"github.com/spf13/cobra"
)

//nolint:lll // Allow long lines for CLI
var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup your vault, register with you identity provider, and connect to an identity network",
	Long: `
The setup command is used to configure your local environment for the Identity CLI tool. With it you can:

- (vault) Connect to a vault, generate and store cryptographic keys
- (idp) Register with an identity provider, such as DUO or Okta, to manage your Agent and MCP identities
- (network) Connect to an identity network, such as AGNTCY, in order to publish your Agent and MCP Server identities and verify those published by others
`,
}

func init() {
	SetupCmd.AddCommand(setup.VaultCmd)
	SetupCmd.AddCommand(setup.IdpCmd)
	SetupCmd.AddCommand(setup.NetworkCmd)
}
