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
	issueFilePath string
)

var IssueFileCmd = &cobra.Command{
	Use:   "file",
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

		// if the file path is not set, prompt the user for it interactively
		if issueFilePath == "" {
			fmt.Fprintf(os.Stderr, "Full file path to the data you want to sign in the badge: \n")
			_, err := fmt.Scanln(&issueFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file path: %v\n", err)
				return
			}
		}
		if issueFilePath == "" {
			fmt.Fprintf(os.Stderr, "No file path provided\n")
			return
		}

		badgeContentData, err := os.ReadFile(issueFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			return
		}

		// Convert the badge value to a string
		badgeContent := string(badgeContentData)

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

//nolint:lll // Allow long lines for CLI
func init() {
	IssueFileCmd.Flags().StringVarP(&issueFilePath, "file-path", "f", "", "The file path to the data you want to sign in the badge")
}
