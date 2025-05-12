// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

//nolint:lll // Allow long lines for CLI
var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage your vault and generate quantum-resistant cryptographic keys",
	Long: `
The Identity tool does not store or share any keys that are used to provide identity to your agents. The tool connects to popular password management applications or crypto vaults to handle the keys.

The keys that are generated via this tool use quantum safe algorithms, and you can find more information on these in our documentation.
`,
}

var vaultConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an existing vault. Options include: txt, keychain, 1Password",
}

var vaultConnectTxtCmd = &cobra.Command{
	Use:   "txt",
	Short: "Connect to a local .txt file",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to local .txt file")
	},
}

var vaultConnectKeychainCmd = &cobra.Command{
	Use:   "keychain",
	Short: "Connect to your local Keychain",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to local Keychain")
	},
}

var vaultConnect1PasswordCmd = &cobra.Command{
	Use:   "1password",
	Short: "Connect to 1Password",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to 1Password")
	},
}

var vaultForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the currently connected wallet",
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
	// Add the different vault connect commands to VaultConnectCmd
	vaultConnectCmd.AddCommand(vaultConnectTxtCmd)
	vaultConnectCmd.AddCommand(vaultConnectKeychainCmd)
	vaultConnectCmd.AddCommand(vaultConnect1PasswordCmd)

	VaultCmd.AddCommand(vaultConnectCmd)
	VaultCmd.AddCommand(vaultForgetCmd)
	VaultCmd.AddCommand(vaultGenerateCmd)
}
