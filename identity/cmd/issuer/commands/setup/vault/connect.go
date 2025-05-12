// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an existing vault. Options include: txt, keychain, 1Password",
}

var connectTxtCmd = &cobra.Command{
	Use:   "txt",
	Short: "Connect to a local .txt file",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to local .txt file")
	},
}

var connectKeychainCmd = &cobra.Command{
	Use:   "keychain",
	Short: "Connect to your local Keychain",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to local Keychain")
	},
}

var connect1PasswordCmd = &cobra.Command{
	Use:   "1password",
	Short: "Connect to 1Password",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Connecting to 1Password")
	},
}

func init() {
	// Add the different vault connect commands to VaultConnectCmd
	ConnectCmd.AddCommand(connectTxtCmd)
	ConnectCmd.AddCommand(connectKeychainCmd)
	ConnectCmd.AddCommand(connect1PasswordCmd)
}
