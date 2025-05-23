// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/core/keystore"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/agntcy/identity/internal/pkg/joseutil"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	filePath  string
	vaultName string
)

var TxtCmd = &cobra.Command{
	Use:   "file",
	Short: "Create a local file with generated cryptographic keys",
	Run: func(cmd *cobra.Command, args []string) {

		// if the file path is not set, prompt the user for it interactively
		if filePath == "" {
			fmt.Fprintf(os.Stderr, "File path to store the generated keys: ")
			_, err := fmt.Scanln(&filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file path: %v\n", err)
				return
			}
		}
		if filePath == "" {
			fmt.Fprintf(os.Stderr, "No file path provided\n")
			return
		}

		// if the vault name is not set, prompt the user for it interactively
		if vaultName == "" {
			fmt.Fprintf(os.Stderr, "Vault name: ")
			_, err := fmt.Scanln(&vaultName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault name: %v\n", err)
				return
			}
		}
		if vaultName == "" {
			fmt.Fprintf(os.Stderr, "No vault name provided\n")
			return
		}

		vaultFilesystemRepository := filesystem.NewVaultFilesystemRepository()
		vaultService := vault.NewVaultService(vaultFilesystemRepository)

		fileStorageConfig := keystore.FileStorageConfig{
			FilePath: filePath,
		}

		service, err := keystore.NewKeyService(keystore.FileStorage, fileStorageConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating key service: %v\n", err)
			return
		}

		priv, err := joseutil.GenerateJWK("RS256", "sig", "test-rsa")
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

		txtConfig := internalIssuerTypes.VaultFile{
			FilePath: filePath,
		}

		var config internalIssuerTypes.VaultConfig = &txtConfig

		vault := internalIssuerTypes.Vault{
			Id:     uuid.NewString(),
			Name:   vaultName,
			Type:   internalIssuerTypes.VaultTypeFile,
			Config: config,
		}

		vaultId, err := vaultService.ConnectVault(&vault)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file vault: %v\n", err)
			return
		}

		cmd.Printf("Successfully created file vault with ID: %s\n", vaultId)

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

func init() {
	// Add flags to the command
	TxtCmd.Flags().StringVarP(&filePath, "file-path", "f", "", "Path to the file")
	TxtCmd.Flags().StringVarP(&vaultName, "name", "n", "", "Name of the vault")
}
