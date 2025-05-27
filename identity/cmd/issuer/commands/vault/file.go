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

var FileCmd = &cobra.Command{
	Use:   "file",
	Short: "Create a local vault file to store your cryptographic keys",
	Run: func(cmd *cobra.Command, args []string) {

		// if the file path is not set, prompt the user for it interactively
		if fileCmdIn.FilePath == "" {
			fmt.Fprintf(os.Stderr, "File path to store the vault: ")
			_, err := fmt.Scanln(&fileCmdIn.FilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file path: %v\n", err)
				return
			}
		}
		if fileCmdIn.FilePath == "" {
			fmt.Fprintf(os.Stderr, "No file path provided\n")
			return
		}

		// if the vault name is not set, prompt the user for it interactively
		if fileCmdIn.VaultName == "" {
			fmt.Fprintf(os.Stderr, "Vault name: ")
			_, err := fmt.Scanln(&fileCmdIn.VaultName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault name: %v\n", err)
				return
			}
		}
		if fileCmdIn.VaultName == "" {
			fmt.Fprintf(os.Stderr, "No vault name provided\n")
			return
		}

		vaultFilesystemRepository := filesystem.NewVaultFilesystemRepository()
		vaultService := vault.NewVaultService(vaultFilesystemRepository)

		fileConfig := vaulttypes.VaultFile{
			FilePath: fileCmdIn.FilePath,
		}

		var config vaulttypes.VaultConfig = &fileConfig

		vault := vaulttypes.Vault{
			Id:     uuid.NewString(),
			Name:   fileCmdIn.VaultName,
			Type:   vaulttypes.VaultTypeFile,
			Config: config,
		}

		vaultId, err := vaultService.ConnectVault(&vault)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error configuring file vault: %v\n", err)
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
	FileCmd.Flags().StringVarP(&fileCmdIn.FilePath, "file-path", "f", "", "Path to the file")
	FileCmd.Flags().StringVarP(&fileCmdIn.VaultName, "name", "n", "", "Name of the vault")
}
