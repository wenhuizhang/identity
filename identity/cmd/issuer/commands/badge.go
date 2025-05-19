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
	Use:   "issue [vault_id] [issuer_id] [metadata_id] [badge_file_path]",
	Short: "Issue a new badge for the current metadata",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		vaultId := args[0]
		issuerId := args[1]
		metadataId := args[2]
		badgeFilePath := args[3]
		_, err := issuerBadge.IssueBadge(vaultId, issuerId, metadataId, badgeFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error issuing badge: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Creating a new badge")
	},
}

//nolint:mnd // Allow magic number for args
var badgePublishCmd = &cobra.Command{
	Use:   "publish [vault_id] [issuer_id] [metadata_id] [badge_id]",
	Short: "Publish the chosen badge",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		vaultId := args[0]
		issuerId := args[1]
		metadataId := args[2]
		badgeId := args[3]

		badge, err := issuerBadge.GetBadge(vaultId, issuerId, metadataId, badgeId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting badge: %v\n", err)
			return
		}

		_, err = issuerBadge.PublishBadge(vaultId, issuerId, metadataId, badge)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error publishing badge: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Publishing the current badge")
	},
}

//nolint:mnd // Allow magic number for args
var badgeListCmd = &cobra.Command{
	Use:   "list [vault_id] [issuer_id] [metadata_id]",
	Short: "List your existing badges for the current metadata",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		vaultId := args[0]
		issuerId := args[1]
		metadataId := args[2]

		badgeIds, err := issuerBadge.ListBadgeIds(vaultId, issuerId, metadataId)
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
	Use:   "show [vault_id] [issuer_id] [metadata_id] [badge_id]",
	Short: "Show details of the chosen badge",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		vaultId := args[0]
		issuerId := args[1]
		metadataId := args[2]
		badgeId := args[3]

		badge, err := issuerBadge.GetBadge(vaultId, issuerId, metadataId, badgeId)
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
	Use:   "forget [vault_id] [issuer_id] [metadata_id] [badge_id]",
	Short: "Forget the chosen badge",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		vaultId := args[0]
		issuerId := args[1]
		metadataId := args[2]
		badgeId := args[3]

		err := issuerBadge.ForgetBadge(vaultId, issuerId, metadataId, badgeId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting badge: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Forgot badge with ID: %s\n", badgeId)
	},
}

func init() {
	BadgeCmd.AddCommand(badgeIssueCmd)
	BadgeCmd.AddCommand(badgePublishCmd)
	BadgeCmd.AddCommand(badgeListCmd)
	BadgeCmd.AddCommand(badgeShowCmd)
	BadgeCmd.AddCommand(badgeForgetCmd)
}
