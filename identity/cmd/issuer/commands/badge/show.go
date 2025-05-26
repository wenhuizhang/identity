// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"encoding/json"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

var badgeShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of the chosen badge",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForBadge()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the badge id is not set, prompt the user for it interactively
		if showCmdIn.BadgeID == "" {
			fmt.Fprintf(os.Stderr, "Badge ID to show:\n")
			_, err := fmt.Scanln(&showCmdIn.BadgeID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
				return
			}
		}
		if showCmdIn.BadgeID == "" {
			fmt.Fprintf(os.Stderr, "No badge ID provided.\n")
			return
		}

		badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, showCmdIn.BadgeID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting badge: %v\n", err)
			return
		}
		badgeJSON, err := json.MarshalIndent(badge, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling badge to JSON: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(badgeJSON))
	},
}
