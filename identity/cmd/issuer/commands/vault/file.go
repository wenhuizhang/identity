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

var (
	filePath  string
	vaultName string
)

var TxtCmd = &cobra.Command{
	Use:   "file",
	Short: "Create a local vault file to store your cryptographic keys",
	Run: func(cmd *cobra.Command, args []string) {

		// if the file path is not set, prompt the user for it interactively
		if filePath == "" {
			fmt.Fprintf(os.Stderr, "File path to store the vault: ")
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

		txtConfig := vaulttypes.VaultFile{
			FilePath: filePath,
		}

		var config vaulttypes.VaultConfig = &txtConfig

		vault := vaulttypes.Vault{
			Id:     uuid.NewString(),
			Name:   vaultName,
			Type:   vaulttypes.VaultTypeFile,
			Config: config,
		}

		vaultId, err := vaultService.ConnectVault(&vault)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file vault: %v\n", err)
			return
		}

		cmd.Printf("Successfully configured file vault with ID: %s\n", vaultId)

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
