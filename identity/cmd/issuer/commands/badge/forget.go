// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"context"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/spf13/cobra"
)

type ForgetFlags struct {
	BadgeID string
}

type ForgetCommand struct {
	cache        *cliCache.Cache
	badgeService badge.BadgeService
}

func NewCmdForget(
	cache *cliCache.Cache,
	badgeService badge.BadgeService,
) *cobra.Command {
	flags := NewForgetFlags()

	cmd := &cobra.Command{
		Use:   "forget",
		Short: "Forget the chosen badge",
		Run: func(cmd *cobra.Command, args []string) {
			c := ForgetCommand{
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

func NewForgetFlags() *ForgetFlags {
	return &ForgetFlags{}
}

func (f *ForgetFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.BadgeID, "badge-id", "b", "", "The ID of the badge to forget")
}

func (cmd *ForgetCommand) Run(ctx context.Context, flags *ForgetFlags) error {
	err := cmd.cache.ValidateForBadge()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	// if the badge id is not set, prompt the user for it interactively
	if flags.BadgeID == "" {
		fmt.Fprintf(os.Stdout, "Badge ID to forget:\n")

		_, err := fmt.Scanln(&flags.BadgeID)
		if err != nil {
			return fmt.Errorf("error reading badge ID: %v", err)
		}
	}

	if flags.BadgeID == "" {
		return fmt.Errorf("no badge ID provided")
	}

	err = cmd.badgeService.ForgetBadge(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		cmd.cache.MetadataId,
		flags.BadgeID,
	)
	if err != nil {
		return fmt.Errorf("error forgetting badge: %v", err)
	}

	// If the badge was the current badge in the cache, clear the cache of badge id
	if cmd.cache.BadgeId == flags.BadgeID {
		cmd.cache.BadgeId = ""
		err = cliCache.SaveCache(cmd.cache)
		if err != nil {
			return fmt.Errorf("error saving local configuration: %v", err)
		}
	}

	fmt.Fprintf(os.Stdout, "Forgot badge with ID: %s\n", flags.BadgeID)

	return nil
}
