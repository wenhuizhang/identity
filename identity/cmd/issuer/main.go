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
