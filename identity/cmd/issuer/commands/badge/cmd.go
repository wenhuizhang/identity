// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"github.com/agntcy/identity/cmd/issuer/commands/badge/issue"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/issuer/badge/data/filesystem"
	"github.com/spf13/cobra"
)

type PublishCmdInput struct {
	BadgeID string
}

type ShowCmdInput struct {
	BadgeID string
}

type ForgetCmdInput struct {
	BadgeID string
}

type LoadCmdInput struct {
	BadgeID string
}

var (
	// setup the badge service
	badgeFilesystemRepository = filesystem.NewBadgeFilesystemRepository()
	badgeService              = badge.NewBadgeService(badgeFilesystemRepository)

	// setup the command flags
	pubCmdIn  = &PublishCmdInput{}
	showCmdIn = &ShowCmdInput{}
	frgtCmdIn = &ForgetCmdInput{}
	loadCmdIn = &LoadCmdInput{}
)

var BadgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Issue and publish badges for your Agent and MCP Server identities",
	Long: `
The badge command is used to issue and publish badges for your Agent and MCP Server identities.
`,
}

func init() {
	badgeIssueCmd.AddCommand(issue.IssueFileCmd)
	badgeIssueCmd.AddCommand(issue.IssueOasfCmd)
	badgeIssueCmd.AddCommand(issue.IssueMcpServerCmd)
	badgeIssueCmd.AddCommand(issue.IssueA2AWellKnownCmd)
	BadgeCmd.AddCommand(badgeIssueCmd)

	badgePublishCmd.Flags().StringVarP(&pubCmdIn.BadgeID, "badge-id", "b", "", "The ID of the badge to publish")
	BadgeCmd.AddCommand(badgePublishCmd)

	BadgeCmd.AddCommand(badgeListCmd)

	badgeShowCmd.Flags().StringVarP(&showCmdIn.BadgeID, "badge-id", "b", "", "The ID of the badge to show")
	BadgeCmd.AddCommand(badgeShowCmd)

	badgeForgetCmd.Flags().StringVarP(&frgtCmdIn.BadgeID, "badge-id", "b", "", "The ID of the badge to forget")
	BadgeCmd.AddCommand(badgeForgetCmd)

	badgeLoadCmd.Flags().StringVarP(&loadCmdIn.BadgeID, "badge-id", "b", "", "The ID of the badge to load")
	BadgeCmd.AddCommand(badgeLoadCmd)
}
