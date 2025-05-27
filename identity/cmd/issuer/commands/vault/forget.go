// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

var vaultForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget a vault configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// if the vault id is not set, prompt the user for it interactively
		if forgetCmdIn.VaultID == "" {
			fmt.Fprintf(os.Stderr, "Vault ID to forget:\n")
			_, err := fmt.Scanln(&forgetCmdIn.VaultID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault ID: %v\n", err)
				return
			}
		}
		if forgetCmdIn.VaultID == "" {
			fmt.Fprintf(os.Stderr, "No vault ID provided.\n")
			return
		}

		err := vaultService.ForgetVault(forgetCmdIn.VaultID)
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

		fmt.Fprintf(os.Stdout, "Forgot vault with ID: %s\n", forgetCmdIn.VaultID)
	},
}
