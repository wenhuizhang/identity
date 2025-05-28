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

type ListCommand struct {
	cache        *cliCache.Cache
	badgeService badge.BadgeService
}

func NewCmdList(
	cache *cliCache.Cache,
	badgeService badge.BadgeService,
) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your existing badges for the current metadata",
		Run: func(cmd *cobra.Command, args []string) {
			c := ListCommand{
				cache:        cache,
				badgeService: badgeService,
			}

			err := c.Run(cmd.Context())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
}

func (cmd *ListCommand) Run(ctx context.Context) error {
	err := cmd.cache.ValidateForBadge()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	badges, err := cmd.badgeService.GetAllBadges(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		cmd.cache.MetadataId,
	)
	if err != nil {
		return fmt.Errorf("error listing badges: %v", err)
	}

	if len(badges) == 0 {
		return fmt.Errorf("no badges found")
	}

	fmt.Fprintf(os.Stdout, "%s\n", "Existing badge ids:")

	for _, badge := range badges {
		fmt.Fprintf(os.Stdout, "- %s\n", badge.Id)
	}

	return nil
}
