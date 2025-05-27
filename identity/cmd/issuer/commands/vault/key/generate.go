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
	"github.com/agntcy/identity/internal/pkg/joseutil"
	"github.com/spf13/cobra"
)

var keyGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new cryptographic key for the vault",
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
		if generateCmdIn.KeyID == "" {
			fmt.Fprintf(os.Stderr, "Key ID: ")
			_, err := fmt.Scanln(&generateCmdIn.KeyID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading key ID: %v\n", err)
				return
			}
		}
		if generateCmdIn.KeyID == "" {
			fmt.Fprintf(os.Stderr, "No key ID provided\n")
			return
		}

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

		priv, err := joseutil.GenerateJWK("RS256", "sig", generateCmdIn.KeyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating JWK: %v\n", err)
			return
		}

		ctx := context.Background()
		err = service.SaveKey(ctx, priv.KID, priv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving key: %v\n", err)
			return
		}

		cmd.Printf("Successfully generated key with ID: %s\n", generateCmdIn.KeyID)

		cache.KeyID = generateCmdIn.KeyID
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}

	},
}
