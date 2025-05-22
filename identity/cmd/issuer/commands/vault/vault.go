// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"encoding/json"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/spf13/cobra"
)

var (
	// setup the vault service
	vaultFilesystemRepository = filesystem.NewVaultFilesystemRepository()
	vaultService              = vault.NewVaultService(vaultFilesystemRepository)

	// setup the vault command flags
	showVaultId   string
	forgetVaultId string
	loadVaultId   string
)

var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage your vaults and generate cryptographic keys",
	Long: `
The vault command is used to configure and manage your vaults. You can use it to:

- (create) Create new vault configurations and generate cryptographic keys
- (list) List your existing vault configurations
- (show) Show details of a vault configuration
- (load) Load a vault configuration
- (forget) Forget a vault configuration
`,
}

var vaultConnectCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vault configuration and generate cryptographic keys",
	Long: `
The create command is used to create a new vault configuration and generate cryptographic keys. You can use:
- (file) Create a local file with generated cryptographic keys
`,
}

var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing vault configurations",
	Run: func(cmd *cobra.Command, args []string) {

		vaults, err := vaultService.GetAllVaults()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing vaults: %v\n", err)
			return
		}
		if len(vaults) == 0 {
			fmt.Fprintf(os.Stdout, "No vaults found.\n")
			return
		}
		fmt.Fprintf(os.Stdout, "Existing vaults:\n")
		for _, vault := range vaults {
			fmt.Fprintf(os.Stdout, "- %s (%s vault), id: %s\n", vault.Name, vault.Type, vault.Id)
		}
	},
}

var vaultShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of a vault configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// if the vault id is not set, prompt the user for it interactively
		if showVaultId == "" {
			fmt.Fprintf(os.Stderr, "Vault ID to show:\n")
			_, err := fmt.Scanln(&showVaultId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault ID: %v\n", err)
				return
			}
		}
		if showVaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault ID provided.\n")
			return
		}

		// check the vault id is valid
		vault, err := vaultService.GetVault(showVaultId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting vault: %v\n", err)
			return
		}
		if vault == nil {
			fmt.Fprintf(os.Stdout, "No vault found with ID: %s\n", showVaultId)
			return
		}

		vaultJSON, err := json.MarshalIndent(vault, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling metadata to JSON: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(vaultJSON))
	},
}

var vaultForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget an vault configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// if the vault id is not set, prompt the user for it interactively
		if forgetVaultId == "" {
			fmt.Fprintf(os.Stderr, "Vault ID to forget:\n")
			_, err := fmt.Scanln(&forgetVaultId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault ID: %v\n", err)
				return
			}
		}
		if forgetVaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault ID provided.\n")
			return
		}

		err := vaultService.ForgetVault(forgetVaultId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting vault: %v\n", err)
			return
		}

		// Remove the cache
		err = cliCache.ClearCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error removing local configuration: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Forgot vault with ID: %s\n", forgetVaultId)
	},
}

var vaultLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a vault configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// if the vault id is not set, prompt the user for it interactively
		if loadVaultId == "" {
			fmt.Fprintf(os.Stderr, "Vault ID to load:\n")
			_, err := fmt.Scanln(&loadVaultId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault ID: %v\n", err)
				return
			}
		}
		if loadVaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault ID provided.\n")
			return
		}

		// check the vault id is valid
		vault, err := vaultService.GetVault(loadVaultId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting vault: %v\n", err)
			return
		}
		if vault == nil {
			fmt.Fprintf(os.Stdout, "No vault found with ID: %s\n", loadVaultId)
			return
		}

		// save the vault id to the cache
		err = cliCache.SaveCache(
			&cliCache.Cache{
				VaultId: vault.Id,
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded vault with ID: %s\n", loadVaultId)

	},
}

func init() {
	// Add the vault types to the vault connect command
	vaultConnectCmd.AddCommand(TxtCmd)

	VaultCmd.AddCommand(vaultConnectCmd)

	VaultCmd.AddCommand(vaultListCmd)

	vaultShowCmd.Flags().StringVarP(&showVaultId, "vault-id", "v", "", "The ID of the vault to show")
	VaultCmd.AddCommand(vaultShowCmd)

	vaultForgetCmd.Flags().StringVarP(&forgetVaultId, "vault-id", "v", "", "The ID of the vault to forget")
	VaultCmd.AddCommand(vaultForgetCmd)

	vaultLoadCmd.Flags().StringVarP(&loadVaultId, "vault-id", "v", "", "The ID of the vault to load")
	VaultCmd.AddCommand(vaultLoadCmd)
}
