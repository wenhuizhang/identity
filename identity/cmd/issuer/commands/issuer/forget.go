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
var issuerForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget an issuer configuration",
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
		if forgetIssuerId == "" {
			fmt.Fprintf(os.Stderr, "Issuer ID to forget:\n")
			_, err := fmt.Scanln(&forgetIssuerId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading issuer ID: %v\n", err)
				return
			}
		}
		if forgetIssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer ID provided.\n")
			return
		}

		err = issuerService.ForgetIssuer(cache.VaultId, forgetIssuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting issuer: %v\n", err)
			return
		}

		// If the issuer was the current issuer in the cache, clear the cache of issuer, metadata, and badge IDs
		if cache.IssuerId == forgetIssuerId {
			cache.IssuerId = ""
			cache.MetadataId = ""
			cache.BadgeId = ""
			err = cliCache.SaveCache(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
				return
			}
		}

		fmt.Fprintf(os.Stdout, "Forgot issuer with ID: %s\n", forgetIssuerId)
	},
}
