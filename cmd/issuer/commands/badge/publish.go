// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"context"
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	badgesrv "github.com/agntcy/identity/internal/issuer/badge"
	issuersrv "github.com/agntcy/identity/internal/issuer/issuer"
	"github.com/agntcy/identity/internal/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type PublishFlags struct {
	BadgeID string
}

type PublishCommand struct {
	cache         *clicache.Cache
	badgeService  badgesrv.BadgeService
	issuerService issuersrv.IssuerService
}

func NewCmdPublish(
	cache *clicache.Cache,
	badgeService badgesrv.BadgeService,
	issuerService issuersrv.IssuerService,
) *cobra.Command {
	flags := NewPublishFlags()

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish the chosen badge",
		Run: func(cmd *cobra.Command, args []string) {
			c := PublishCommand{
				cache:         cache,
				badgeService:  badgeService,
				issuerService: issuerService,
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

func NewPublishFlags() *PublishFlags {
	return &PublishFlags{}
}

func (f *PublishFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.BadgeID, "badge-id", "b", "", "The ID of the badge to publish")
}

func (cmd *PublishCommand) Run(ctx context.Context, flags *PublishFlags) error {
	err := cmd.cache.ValidateForBadge()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %w", err)
	}

	// if the badge id is not set, prompt the user for it interactively
	// if there is a badge id in the cache, use it as the default when prompting
	if cmd.cache.BadgeId != "" {
		err = cmdutil.ScanWithDefaultIfNotSet("Badge ID to publish", cmd.cache.BadgeId, &flags.BadgeID)
	} else {
		err = cmdutil.ScanRequiredIfNotSet("Badge ID to publish", &flags.BadgeID)
	}

	if err != nil {
		return fmt.Errorf("error reading badge ID: %w", err)
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

	_, err = cmd.badgeService.PublishBadge(
		ctx,
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		cmd.cache.MetadataId,
		badge,
	)
	if err != nil {
		return fmt.Errorf("error publishing badge: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", "Published the badge\n")

	issuer, err := cmd.issuerService.GetIssuer(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
	)
	if err != nil {
		return fmt.Errorf("error getting issuer: %w", err)
	}

	fmt.Fprintf(os.Stdout,
		"You can access the published badges for your metadata as a Well-Known at "+
			"%s/v1alpha1/vc/%s/.well-known/vcs.json\n\n",
		issuer.IdentityNodeURL,
		cmd.cache.MetadataId)

	fmt.Fprintf(os.Stdout,
		"To download the badges for verification, you can use the following command:\n"+
			"curl -o vcs.json %s/v1alpha1/vc/%s/.well-known/vcs.json\n\n",
		issuer.IdentityNodeURL,
		cmd.cache.MetadataId,
	)

	return nil
}
