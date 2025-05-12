// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	"fmt"

	"github.com/agntcy/identity/cmd/issuer/commands/setup/vault"
	"github.com/spf13/cobra"
)

var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage your vault and generate quantum-resistant cryptographic keys",
	Long: `
The Identity tool does not store or share any keys that are used to provide identity to your agents.
The tool connects to popular password management applications or crypto vaults to handle the keys.

The vault command is used to manage your vault and generate quantum-resistant cryptographic keys. With it you can:
- (connect) Connect to a local vault, such as a .txt file or Keychain
- (generate) Generate quantum-resistant cryptographic keys and store them in your vault
- (forget) Forget the currently connected vault
`,
}

var vaultForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the currently connected vault",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println(
			"Forgetting the currently connected vault. Please connect to a new vault to continue.",
		)
	},
}

//nolint:lll // Allow long lines for CLI
var vaultGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate quantum-resistant cryptographic keys and store them in your vault",
	Long: `Generate quantum-resistant cryptographic keys and store them in your vault

In order for other users to use and verify the identity of agents you publish, you will have to also publish your public key in one of the supported Trust Anchors. You can find out more about the trust anchors and how to publish the public key in our documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println(
			"Generating quantum-resistant cryptographic keys and storing them in your vault",
		)
	},
}

func init() {
	VaultCmd.AddCommand(vault.ConnectCmd)
	VaultCmd.AddCommand(vaultForgetCmd)
	VaultCmd.AddCommand(vaultGenerateCmd)
}
