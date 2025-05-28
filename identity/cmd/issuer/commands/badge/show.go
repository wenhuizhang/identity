// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/spf13/cobra"
)

type ShowFlags struct {
	BadgeID string
}

type ShowCommand struct {
	cache        *clicache.Cache
	badgeService badge.BadgeService
}

func NewCmdShow(
	cache *clicache.Cache,
	badgeService badge.BadgeService,
) *cobra.Command {
	flags := NewShowFlags()

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show details of the chosen badge",
		Run: func(cmd *cobra.Command, args []string) {
			c := ShowCommand{
				cache:        cache,
				badgeService: badgeService,
			}

			err := c.Run(cmd.Context(), flags)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	flags.AddFlags(cmd)

	return cmd
}

func NewShowFlags() *ShowFlags {
	return &ShowFlags{}
}

func (f *ShowFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.BadgeID, "badge-id", "b", "", "The ID of the badge to show")
}

func (cmd *ShowCommand) Run(ctx context.Context, flags *ShowFlags) error {
	err := cmd.cache.ValidateForBadge()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	// if the badge id is not set, prompt the user for it interactively
	if flags.BadgeID == "" {
		fmt.Fprintf(os.Stdout, "Badge ID to show:\n")

		_, err := fmt.Scanln(&flags.BadgeID)
		if err != nil {
			return fmt.Errorf("error reading badge ID: %v", err)
		}
	}

	if flags.BadgeID == "" {
		return fmt.Errorf("no badge ID provided")
	}

	badge, err := cmd.badgeService.GetBadge(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		cmd.cache.MetadataId,
		flags.BadgeID,
	)
	if err != nil {
		return fmt.Errorf("error getting badge: %v", err)
	}

	badgeJSON, err := json.MarshalIndent(badge, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling badge to JSON: %v", err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", string(badgeJSON))

	return nil
}
