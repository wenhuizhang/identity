// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

//nolint:mnd // Allow magic number for args
var OnePasswordCmd = &cobra.Command{
	Use:   "1password [service-account-token] [vault-id] [item-id]",
	Short: "Connect to your 1Password account",
	Long:  "Connect to your 1Password account",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		serviceAccountToken := args[0]
		vaultId := args[1]
		itemId := args[2]

		vaultFilesystemRepository := filesystem.NewVaultFilesystemRepository()
		vaultService := vault.NewVaultService(vaultFilesystemRepository)

		onePasswordConfig := internalIssuerTypes.Vault1Password{
			ServiceAccountToken: serviceAccountToken,
			VaultID:             vaultId,
			ItemID:              itemId,
		}
		var config internalIssuerTypes.VaultConfig = &onePasswordConfig

		vault := internalIssuerTypes.Vault{
			Id:     uuid.NewString(),
			Name:   vaultName,
			Type:   internalIssuerTypes.VaultType1Password,
			Config: config,
		}

		vaultId, err := vaultService.ConnectVault(&vault)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating 1Password vault: %v\n", err)
			return
		}

		cmd.Printf("Successfully created 1Password vault with ID: %s\n", vaultId)

		err = cliCache.SaveCache(
			&cliCache.Cache{
				VaultId: vaultId,
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
			return
		}

	},
}

func init() {}
