// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"context"
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/spf13/cobra"
)

type LoadFlags struct {
	BadgeID string
}

type LoadCommand struct {
	cache        *clicache.Cache
	badgeService badge.BadgeService
}

func NewCmdLoad(
	cache *clicache.Cache,
	badgeService badge.BadgeService,
) *cobra.Command {
	flags := NewLoadFlags()

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load a badge configuration",
		Run: func(cmd *cobra.Command, args []string) {
			c := LoadCommand{
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

func NewLoadFlags() *LoadFlags {
	return &LoadFlags{}
}

func (f *LoadFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.BadgeID, "badge-id", "b", "", "The ID of the badge to load")
}

func (cmd *LoadCommand) Run(ctx context.Context, flags *LoadFlags) error {
	err := cmd.cache.ValidateForBadge()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	// if the badge id is not set, prompt the user for it interactively
	if flags.BadgeID == "" {
		fmt.Fprintf(os.Stdout, "Badge ID to load:\n")

		_, err := fmt.Scanln(&flags.BadgeID)
		if err != nil {
			return fmt.Errorf("error reading badge ID: %v", err)
		}
	}

	if flags.BadgeID == "" {
		return fmt.Errorf("no badge ID provided")
	}

	// check the badge id is valid
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

	if badge == nil {
		return fmt.Errorf("badge with ID %s not found", flags.BadgeID)
	}

	// save the metadata id to the cache
	cmd.cache.BadgeId = flags.BadgeID

	err = clicache.SaveCache(cmd.cache)
	if err != nil {
		return fmt.Errorf("error saving local configuration: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Loaded badge with ID: %s\n", flags.BadgeID)

	return nil
}
