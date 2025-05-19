// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/agntcy/identity/cmd/issuer/commands/vaults"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage your vault and generate cryptographic keys",
	Long: `
The setup command is used to configure your local environment for the Identity CLI tool. With it you can:

- (connect) Manage your vault and generate cryptographic keys
- (list) List your existing vault configurations
- (show) Show details of a vault configuration
- (forget) Forget a vault configuration
`,
}

//nolint:mnd // Allow magic number for args
var vaultConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Manage your vault and generate cryptographic keys",
	Long: `
The connect command is used to manage your vault and generate cryptographic keys. With it you can:
- (txt) Connect to a local .txt file and generate cryptographic keys
- (1password) Connect to 1Password and generate cryptographic keys
`,
}

var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing vault configurations",
	Run: func(cmd *cobra.Command, args []string) {

		vaults, err := vault.ListVaultIds()
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
			fmt.Fprintf(os.Stdout, "- %s\n", vault)
		}
	},
}
var vaultShowCmd = &cobra.Command{
	Use:   "show [vault_id]",
	Short: "Show details of a vault configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultId := args[0]
		vault, err := vault.GetVault(vaultId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting vault: %v\n", err)
			return
		}
		if vault == nil {
			fmt.Fprintf(os.Stdout, "No vault found with ID: %s\n", vaultId)
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
	Use:   "forget [vault_id]",
	Short: "Forget an vault configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultId := args[0]
		err := vault.ForgetVault(vaultId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting vault: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Forgot vault with ID: %s\n", vaultId)
	},
}

func init() {

	// Add the vault types to the vault connect command
	vaultConnectCmd.AddCommand(vaults.TxtCmd)
	vaultConnectCmd.AddCommand(vaults.OnePasswordCmd)

	VaultCmd.AddCommand(vaultConnectCmd)
	VaultCmd.AddCommand(vaultListCmd)
	VaultCmd.AddCommand(vaultShowCmd)
	VaultCmd.AddCommand(vaultForgetCmd)
}
