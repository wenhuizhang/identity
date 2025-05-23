// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

var issuerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing issuer configurations",
	Long:  "List your existing issuer configurations",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}

		err = cache.ValidateForIssuer()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		issuers, err := issuerService.GetAllIssuers(cache.VaultId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing issuers: %v\n", err)
			return
		}
		if len(issuers) == 0 {
			fmt.Fprintf(os.Stdout, "No issuers found.\n")
			return
		}
		fmt.Fprintf(os.Stdout, "Existing issuers:\n")
		for _, issuer := range issuers {
			fmt.Fprintf(os.Stdout, "- %s, %s\n", issuer.ID, issuer.CommonName)
		}
	},
}
