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

var keyLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a key from the vault",
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

		// if the key id is not set, prompt the user for it interactively
		if loadCmdIn.KeyID == "" {
			fmt.Fprintf(os.Stderr, "Key ID: ")
			_, err := fmt.Scanln(&loadCmdIn.KeyID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading key ID: %v\n", err)
				return
			}
		}
		if loadCmdIn.KeyID == "" {
			fmt.Fprintf(os.Stderr, "No key ID provided\n")
			return
		}

		// get the vault configuration
		vault, err := vaultService.GetVault(cache.VaultId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting vault: %v\n", err)
			return
		}

		var service keystore.KeyService

		if vault.Type == vaulttypes.VaultTypeFile {

			fileStorageConfig := keystore.FileStorageConfig{
				FilePath: vault.Config.(*vaulttypes.VaultFile).FilePath,
			}

			service, err = keystore.NewKeyService(keystore.FileStorage, fileStorageConfig)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating key service: %v\n", err)
				return
			}
		}

		ctx := context.Background()

		// check if the key exists
		_, err = service.RetrievePrivKey(ctx, loadCmdIn.KeyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving private key: %v\n", err)
			return
		}
		_, err = service.RetrievePubKey(ctx, loadCmdIn.KeyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving public key: %v\n", err)
			return
		}

		// save the key id to the cache
		cache.KeyID = loadCmdIn.KeyID
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded Key with ID: %s\n", loadCmdIn.KeyID)

	},
}
