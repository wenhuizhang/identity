// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issue

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/issuer/badge/data/filesystem"
	"github.com/spf13/cobra"
)

var (
	// setup the command flags
	issueA2AWellKnown string
)

//nolint:lll // Allow long lines for CLI
var IssueA2AWellKnownCmd = &cobra.Command{
	Use:   "a2a",
	Short: "Issue a badge based on a local file",
	Run: func(cmd *cobra.Command, args []string) {

		// setup the badge service
		badgeFilesystemRepository := filesystem.NewBadgeFilesystemRepository()
		badgeService := badge.NewBadgeService(badgeFilesystemRepository)

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

		// if the mcp server url is not set, prompt the user for it interactively
		if issueA2AWellKnown == "" {
			fmt.Fprintf(os.Stderr, "Well-known URL of the A2A agent you want to sign in the badge: \n")
			_, err := fmt.Scanln(&issueA2AWellKnown)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading A2A well-known URL: %v\n", err)
				return
			}
		}
		if issueA2AWellKnown == "" {
			fmt.Fprintf(os.Stderr, "No A2A well-known URL provided\n")
			return
		}

		// Convert the badge value to a string
		badgeContent := "A2A Agent Well-known URL"

		badgeId, err := badgeService.IssueBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, badgeContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error issuing badge: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Issued badge with ID: %s\n", badgeId)

		// Save the badge ID to the cache
		cache.BadgeId = badgeId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
	},
}

func init() {
	IssueA2AWellKnownCmd.Flags().StringVarP(&issueA2AWellKnown, "url", "u", "", "The well-known URL of the A2A agent you want to sign in the badge")
}
