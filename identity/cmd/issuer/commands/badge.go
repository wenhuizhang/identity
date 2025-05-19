// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"encoding/json"
	"fmt"
	"os"

	issuerBadge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/spf13/cobra"
)

var BadgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Issue and publish badges for your Agent and MCP Server identities",
	Long: `
The badge command is used to issue and publish badges for your Agent and MCP Server identities. With it you can:

- (issue) Issue a new badge for the issuer and metadata
- (publish) Publish the chosen badge
- (list) List your existing badges for the current issuer and metadata
- (show) Show details of the chosen badge
- (forget) Forget a specific badge
`,
}

//nolint:mnd // Allow magic number for args
var badgeIssueCmd = &cobra.Command{
	Use:   "issue [issuer_id] [metadata_id] [badge_file_path]",
	Short: "Issue a new badge for the current metadata",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		issuerId := args[0]
		metadataId := args[1]
		badgeFilePath := args[2]
		_, err := issuerBadge.IssueBadge(issuerId, metadataId, badgeFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error issuing badge: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Creating a new badge")
	},
}

var badgePublishCmd = &cobra.Command{
	Use:   "publish [issuer_id] [metadata_id] [badge_id]",
	Short: "Publish the chosen badge",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		issuerId := args[0]
		metadataId := args[1]
		badgeId := args[2]

		badge, err := issuerBadge.GetBadge(issuerId, metadataId, badgeId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting badge: %v\n", err)
			return
		}

		_, err = issuerBadge.PublishBadge(issuerId, metadataId, badge)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error publishing badge: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Publishing the current badge")
	},
}

//nolint:mnd // Allow magic number for args
var badgeListCmd = &cobra.Command{
	Use:   "list [issuer_id] [metadata_id]",
	Short: "List your existing badges for the current metadata",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		issuerId := args[0]
		metadataId := args[1]
		badgeIds, err := issuerBadge.ListBadgeIds(issuerId, metadataId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing badges: %v\n", err)
			return
		}
		if len(badgeIds) == 0 {
			fmt.Fprintf(os.Stdout, "%s\n", "No badges found")
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Existing badge ids:")
		for _, badgeId := range badgeIds {
			fmt.Fprintf(os.Stdout, "- %s\n", badgeId)
		}
	},
}

//nolint:mnd // Allow magic number for args
var badgeShowCmd = &cobra.Command{
	Use:   "show [issuer_id] [metadata_id] [badge_id]",
	Short: "Show details of the chosen badge",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		issuerId := args[0]
		metadataId := args[1]
		badgeId := args[3]
		badge, err := issuerBadge.GetBadge(issuerId, metadataId, badgeId)
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

//nolint:mnd // Allow magic number for args
var badgeForgetCmd = &cobra.Command{
	Use:   "forget [issuer_id] [metadata_id] [badge_id]",
	Short: "Forget the chosen badge",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		issuerId := args[0]
		metadataId := args[1]
		badgeId := args[2]
		err := issuerBadge.ForgetBadge(issuerId, metadataId, badgeId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting badge: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Forgetting the current badge")
	},
}

func init() {
	BadgeCmd.AddCommand(badgeIssueCmd)
	BadgeCmd.AddCommand(badgePublishCmd)
	BadgeCmd.AddCommand(badgeListCmd)
	BadgeCmd.AddCommand(badgeShowCmd)
	BadgeCmd.AddCommand(badgeForgetCmd)
}
