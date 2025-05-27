// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

var vaultLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a vault configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// if the vault id is not set, prompt the user for it interactively
		if loadCmdIn.VaultID == "" {
			fmt.Fprintf(os.Stderr, "Vault ID to load:\n")
			_, err := fmt.Scanln(&loadCmdIn.VaultID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault ID: %v\n", err)
				return
			}
		}
		if loadCmdIn.VaultID == "" {
			fmt.Fprintf(os.Stderr, "No vault ID provided.\n")
			return
		}

		// check the vault id is valid
		vault, err := vaultService.GetVault(loadCmdIn.VaultID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting vault: %v\n", err)
			return
		}
		if vault == nil {
			fmt.Fprintf(os.Stdout, "No vault found with ID: %s\n", loadCmdIn.VaultID)
			return
		}

		// save the vault id to the cache
		err = cliCache.SaveCache(
			&cliCache.Cache{
				VaultId: vault.Id,
				//nolint:godox // To be fixed in the next PR
				// TODO: load the KID
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded vault with ID: %s\n", loadCmdIn.VaultID)

	},
}
