// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	vaulttypes "github.com/agntcy/identity/internal/issuer/vault/types"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

//nolint:mnd // Allow magic number for args
var OnePasswordCmd = &cobra.Command{
	Use:   "1password [service-account-token] [vault-id] [item-id]",
	Short: "Connect to your 1Password account",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		serviceAccountToken := args[0]
		vaultId := args[1]
		itemId := args[2]

		vaultFilesystemRepository := filesystem.NewVaultFilesystemRepository()
		vaultService := vault.NewVaultService(vaultFilesystemRepository)

		onePasswordConfig := vaulttypes.Vault1Password{
			ServiceAccountToken: serviceAccountToken,
			VaultID:             vaultId,
			ItemID:              itemId,
		}
		var config vaulttypes.VaultConfig = &onePasswordConfig

		vault := vaulttypes.Vault{
			Id:     uuid.NewString(),
			Name:   vaultName,
			Type:   vaulttypes.VaultType1Password,
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
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}

	},
}

func init() {}
