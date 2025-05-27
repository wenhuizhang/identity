// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"context"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/core/keystore"
	vaulttypes "github.com/agntcy/identity/internal/issuer/vault/types"
	"github.com/spf13/cobra"
)

var keyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys in the vault",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// get the vault configuration
		vault, err := vaultService.GetVault(cache.VaultId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting vault: %v\n", err)
			return
		}

		var service keystore.KeyService

		switch vault.Type {
		case vaulttypes.VaultTypeFile:
			fileVault, ok := vault.Config.(*vaulttypes.VaultFile)
			if !ok {
				fmt.Fprintf(os.Stderr, "Error: vault config is not of type VaultFile\n")
				return
			}
			fileConfig := keystore.FileStorageConfig{
				FilePath: fileVault.FilePath,
			}
			service, err = keystore.NewKeyService(keystore.FileStorage, fileConfig)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating key service: %v\n", err)
				return
			}

		case vaulttypes.VaultTypeHashicorp:
			hashicorpVault, ok := vault.Config.(*vaulttypes.VaultHashicorp)
			if !ok {
				fmt.Fprintf(os.Stderr, "Error: vault config is not of type VaultHashicorp\n")
				return
			}
			hashicorpConfig := keystore.VaultStorageConfig{
				Address:   hashicorpVault.Address,
				Token:     hashicorpVault.Token,
				Namespace: hashicorpVault.Namespace,
			}
			service, err = keystore.NewKeyService(keystore.VaultStorage, hashicorpConfig)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating key service: %v\n", err)
				return
			}

		default:
			fmt.Fprintf(os.Stderr, "Unsupported vault type: %s\n", vault.Type)
			return
		}

		ctx := context.Background()

		keys, err := service.ListKeys(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing keys: %v\n", err)
			return
		}
		if len(keys) == 0 {
			fmt.Fprintf(os.Stderr, "No keys found in the vault\n")
			return
		}
		fmt.Fprintf(os.Stdout, "Keys in vault '%s':\n", vault.Name)
		for _, key := range keys {
			fmt.Fprintf(os.Stdout, "- %s\n", key)
		}
	},
}
