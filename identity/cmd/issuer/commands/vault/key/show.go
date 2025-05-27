// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/core/keystore"
	vaulttypes "github.com/agntcy/identity/internal/issuer/vault/types"
	"github.com/spf13/cobra"
)

var keyShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of a specific key in the vault",
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
		if showCmdIn.KeyID == "" {
			fmt.Fprintf(os.Stderr, "Key ID: ")
			_, err := fmt.Scanln(&showCmdIn.KeyID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading key ID: %v\n", err)
				return
			}
		}
		if showCmdIn.KeyID == "" {
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

		publicKey, err := service.RetrievePubKey(ctx, showCmdIn.KeyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving public key: %v\n", err)
			return
		}

		if publicKey == nil {
			fmt.Fprintf(os.Stderr, "No public key found for Key ID: %s\n", showCmdIn.KeyID)
			return
		}

		// convert the public key to a string representation
		publicKeyStr, err := json.MarshalIndent(publicKey, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling public key: %v\n", err)
			return
		}

		privateKey, err := service.RetrievePrivKey(ctx, showCmdIn.KeyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving private key: %v\n", err)
			return
		}

		// convert the private key to a string representation
		if privateKey == nil {
			fmt.Fprintf(os.Stderr, "No private key found for Key ID: %s\n", showCmdIn.KeyID)
			return
		}

		privateKeyStr, err := json.MarshalIndent(privateKey, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling private key: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nKey ID: %s\n", showCmdIn.KeyID)
		fmt.Fprintf(os.Stdout, "\nPublic Key: %s\n", publicKeyStr)
		fmt.Fprintf(os.Stdout, "\nPrivate Key: %s\n", privateKeyStr)
	},
}
