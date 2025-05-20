// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vaults

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/spf13/cobra"
)

var TxtCmd = &cobra.Command{
	Use:   "txt [output_path]",
	Short: "Connect to .txt file",
	Long:  "Connect to .txt file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultFilesystemRepository := filesystem.NewVaultFilesystemRepository()
		vaultService := vault.NewVaultService(vaultFilesystemRepository)

		txtConfig := internalIssuerTypes.VaultTxt{
			Path: args[0],
		}
		var config internalIssuerTypes.VaultConfig = &txtConfig

		vault, err := vaultService.ConnectVault(internalIssuerTypes.VaultTypeTxt, &config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to vault: %v\n", err)
			return
		}

		cmd.Printf("Successfully connected to vault: %s\n", vault.Id)

		cliCache.SaveCache(
			&cliCache.Cache{
				VaultId: vault.Id,
			},
		)
	},
}

func init() {}
