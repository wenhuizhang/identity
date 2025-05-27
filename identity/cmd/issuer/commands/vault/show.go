// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var vaultShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of a vault configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// if the vault id is not set, prompt the user for it interactively
		if showCmdIn.VaultID == "" {
			fmt.Fprintf(os.Stderr, "Vault ID to show:\n")
			_, err := fmt.Scanln(&showCmdIn.VaultID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault ID: %v\n", err)
				return
			}
		}
		if showCmdIn.VaultID == "" {
			fmt.Fprintf(os.Stderr, "No vault ID provided.\n")
			return
		}

		// check the vault id is valid
		vault, err := vaultService.GetVault(showCmdIn.VaultID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting vault: %v\n", err)
			return
		}
		if vault == nil {
			fmt.Fprintf(os.Stdout, "No vault found with ID: %s\n", showCmdIn.VaultID)
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
