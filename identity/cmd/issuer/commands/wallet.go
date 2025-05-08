// Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

//nolint:lll // Allow long lines for CLI
var WalletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Manage your wallet and generate quantum-resistant cryptographic keys",
	Long: `
The Identity tool does not store or share any keys that are used to provide identity to your agents. The tool connects to popular password management applications or crypto wallets to handle the keys.

The keys that are generated via this tool use quantum safe algorithms, and you can find more information on these in our documentation.
`,
}

//nolint:lll // Allow long lines for CLI
var walletConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an existing wallet. Options include: 1Password, Bitwarden, Dropbox, ProtonPass, Dashlane, Zoho Vault, Keeper",
}

var walletConnect1PasswordCmd = &cobra.Command{
	Use:   "1password",
	Short: "Connect to 1Password",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to 1Password")
	},
}

var walletConnectBitwardenCmd = &cobra.Command{
	Use:   "bitwarden",
	Short: "Connect to Bitwarden",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to Bitwarden")
	},
}

var walletConnectDropboxCmd = &cobra.Command{
	Use:   "dropbox",
	Short: "Connect to Dropbox",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to Dropbox")
	},
}

var walletConnectProtonPassCmd = &cobra.Command{
	Use:   "protonpass",
	Short: "Connect to ProtonPass",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to ProtonPass")
	},
}

var walletConnectDashlaneCmd = &cobra.Command{
	Use:   "dashlane",
	Short: "Connect to Dashlane",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to Dashlane")
	},
}

var walletConnectZohoVaultCmd = &cobra.Command{
	Use:   "zohovault",
	Short: "Connect to Zoho Vault",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to Zoho Vault")
	},
}

var walletConnectKeeperCmd = &cobra.Command{
	Use:   "keeper",
	Short: "Connect to Keeper",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to Keeper")
	},
}

var walletForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the currently connected wallet",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println(
			"Forgetting the currently connected wallet. Please connect to a new wallet to continue.",
		)
	},
}

//nolint:lll // Allow long lines for CLI
var walletGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate quantum-resistant cryptographic keys and store them in your wallet",
	Long: `Generate quantum-resistant cryptographic keys and store them in your wallet

In order for other users to use and verify the identity of agents you publish, you will have to also publish your public key in one of the supported Trust Anchors. You can find out more about the trust anchors and how to publish the public key in our documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println(
			"Generating quantum-resistant cryptographic keys and storing them in your wallet",
		)
	},
}

func init() {
	// Add the different wallet connect commands to WalletConnectCmd
	walletConnectCmd.AddCommand(walletConnect1PasswordCmd)
	walletConnectCmd.AddCommand(walletConnectBitwardenCmd)
	walletConnectCmd.AddCommand(walletConnectDropboxCmd)
	walletConnectCmd.AddCommand(walletConnectProtonPassCmd)
	walletConnectCmd.AddCommand(walletConnectDashlaneCmd)
	walletConnectCmd.AddCommand(walletConnectZohoVaultCmd)
	walletConnectCmd.AddCommand(walletConnectKeeperCmd)

	WalletCmd.AddCommand(walletConnectCmd)
	WalletCmd.AddCommand(walletForgetCmd)
	WalletCmd.AddCommand(walletGenerateCmd)
}
