// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	badgecmd "github.com/agntcy/identity/cmd/issuer/commands/badge"
	"github.com/agntcy/identity/cmd/issuer/commands/configuration"
	"github.com/agntcy/identity/cmd/issuer/commands/issuer"
	"github.com/agntcy/identity/cmd/issuer/commands/metadata"
	vaultcmd "github.com/agntcy/identity/cmd/issuer/commands/vault"
	"github.com/agntcy/identity/cmd/issuer/commands/verify"
	"github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/issuer/badge/a2a"
	badgefs "github.com/agntcy/identity/internal/issuer/badge/data/filesystem"
	"github.com/agntcy/identity/internal/issuer/badge/mcp"
	issuerfs "github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	mdfs "github.com/agntcy/identity/internal/issuer/metadata/data/filesystem"
	"github.com/agntcy/identity/internal/issuer/vault"
	vaultfs "github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/agntcy/identity/internal/pkg/nodeapi"
	"github.com/agntcy/identity/internal/pkg/oidc"

	"github.com/spf13/cobra"
)

func main() {
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

	// load the cache to get the vault, issuer, metadata and badge ids
	cache, err := clicache.LoadCache()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize repositories
	badgeFilesystemRepository := badgefs.NewBadgeFilesystemRepository()
	issuerRepository := issuerfs.NewIssuerFilesystemRepository()
	mdRepository := mdfs.NewMetadataFilesystemRepository()
	vaultRepository := vaultfs.NewVaultFilesystemRepository()

	// Initialize clients
	a2aClient := a2a.NewDiscoveryClient()
	mcpClient := mcp.NewDiscoveryClient()
	nodeClientPrv := nodeapi.NewNodeClientProvider()

	oidcAuth := oidc.NewAuthenticator()

	// Initialize services
	badgeService := badge.NewBadgeService(
		badgeFilesystemRepository,
		mdRepository,
		issuerRepository,
		oidcAuth,
		nodeClientPrv,
	)
	vaultSrv := vault.NewVaultService(vaultRepository)

	rootCmd.AddCommand(vaultcmd.VaultCmd)
	rootCmd.AddCommand(issuer.IssuerCmd)
	rootCmd.AddCommand(metadata.MetadataCmd)

	rootCmd.AddCommand(badgecmd.NewCmd(
		cache,
		badgeService,
		vaultSrv,
		a2aClient,
		mcpClient,
	))

	rootCmd.AddCommand(verify.VerifyCmd)
	rootCmd.AddCommand(configuration.ConfigurationCmd)

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
