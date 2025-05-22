// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/agntcy/identity/cmd/issuer/commands/badge"
	"github.com/agntcy/identity/cmd/issuer/commands/configuration"
	"github.com/agntcy/identity/cmd/issuer/commands/issuer"
	"github.com/agntcy/identity/cmd/issuer/commands/metadata"
	"github.com/agntcy/identity/cmd/issuer/commands/vault"
	"github.com/agntcy/identity/cmd/issuer/commands/verify"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands

//nolint:lll // Allow long lines for CLI
var rootCmd = &cobra.Command{
	Use: "identity",
	Long: `
The Identity CLI tool is a command line interface for generating, publishing and verifying identities within the Internet of Agents.

With it you can:

- (vault) Manage your vault and generate cryptographic keys
- (issuer) Register as an Issuer with an Identity Network
- (metadata) Issue and publish important metadata for your Agent and MCP Server identities
- (badge) Issue and publish badges for your Agent and MCP Server identities
- (verify) Verify the identity of other Agents and MCP Servers via their resolver metadata and badges
- (config) View the current configuration context of your Identity CLI tool
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(vault.VaultCmd)
	rootCmd.AddCommand(issuer.IssuerCmd)
	rootCmd.AddCommand(metadata.MetadataCmd)
	rootCmd.AddCommand(badge.BadgeCmd)
	rootCmd.AddCommand(verify.VerifyCmd)
	rootCmd.AddCommand(configuration.ConfigurationCmd)
}
