// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	badgesrv "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type ShowFlags struct {
	BadgeID string
}

type ShowCommand struct {
	cache        *clicache.Cache
	badgeService badgesrv.BadgeService
}

func NewCmdShow(
	cache *clicache.Cache,
	badgeService badgesrv.BadgeService,
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
		return fmt.Errorf("error validating local configuration: %w", err)
	}

	// if the badge id is not set, prompt the user for it interactively
	if flags.BadgeID == "" {
		err := cmdutil.ScanRequired("Badge ID to show", &flags.BadgeID)
		if err != nil {
			return fmt.Errorf("error reading badge ID: %w", err)
		}
	}

	badge, err := cmd.badgeService.GetBadge(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		cmd.cache.MetadataId,
		flags.BadgeID,
	)
	if err != nil {
		return fmt.Errorf("error getting badge: %w", err)
	}

	badgeJSON, err := json.MarshalIndent(badge, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling badge to JSON: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", string(badgeJSON))

	return nil
}
