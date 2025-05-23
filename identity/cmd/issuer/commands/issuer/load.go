// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

//nolint:lll // Allow long lines for CLI
var issuerLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load an issuer configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}

		// if the issuer id is not set, prompt the user for it interactively
		if loadIssuerId == "" {
			fmt.Fprintf(os.Stderr, "Issuer ID to load:\n")
			_, err := fmt.Scanln(&loadIssuerId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading issuer ID: %v\n", err)
				return
			}
		}
		if loadIssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer ID provided.\n")
			return
		}

		// check the issuer id is valid
		issuer, err := issuerService.GetIssuer(cache.VaultId, loadIssuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting issuer: %v\n", err)
			return
		}
		if issuer == nil {
			fmt.Fprintf(os.Stderr, "No issuer found with ID: %s\n", loadIssuerId)
			return
		}

		// save the issuer id to the cache
		cache.IssuerId = loadIssuerId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded issuer with ID: %s\n", loadIssuerId)

	},
}
