// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"encoding/json"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

//nolint:lll // Allow long lines for CLI
var issuerShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of an issuer configuration",
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
		if showIssuerId == "" {
			fmt.Fprintf(os.Stderr, "Issuer ID to show:\n")
			_, err := fmt.Scanln(&showIssuerId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading issuer ID: %v\n", err)
				return
			}
		}
		if showIssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer ID provided.\n")
			return
		}

		issuer, err := issuerService.GetIssuer(cache.VaultId, cache.KeyID, showIssuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting issuer: %v\n", err)
			return
		}
		if issuer == nil {
			fmt.Fprintf(os.Stdout, "No issuer found with ID: %s\n", showIssuerId)
			return
		}

		issuerJSON, err := json.MarshalIndent(issuer, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling metadata to JSON: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(issuerJSON))
	},
}
