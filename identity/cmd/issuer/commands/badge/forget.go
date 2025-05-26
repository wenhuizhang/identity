// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

var badgeForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the chosen badge",
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
		if frgtCmdIn.BadgeID == "" {
			fmt.Fprintf(os.Stderr, "Badge ID to forget:\n")
			_, err := fmt.Scanln(&frgtCmdIn.BadgeID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
				return
			}
		}
		if frgtCmdIn.BadgeID == "" {
			fmt.Fprintf(os.Stderr, "No badge ID provided.\n")
			return
		}

		err = badgeService.ForgetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, frgtCmdIn.BadgeID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting badge: %v\n", err)
			return
		}

		// If the badge was the current badge in the cache, clear the cache of badge id
		if cache.BadgeId == frgtCmdIn.BadgeID {
			cache.BadgeId = ""
			err = cliCache.SaveCache(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
				return
			}
		}

		fmt.Fprintf(os.Stdout, "Forgot badge with ID: %s\n", frgtCmdIn.BadgeID)
	},
}
