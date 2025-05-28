// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/metadata"
)

type ForgetFlags struct {
	MetadataID string
}

type ForgetCommand struct {
	cache           *clicache.Cache
	metadataService metadata.MetadataService
}

func NewCmdForget(
	cache *clicache.Cache,
	metadataService metadata.MetadataService,
) *cobra.Command {
	flags := NewForgetFlags()

	cmd := &cobra.Command{
		Use:   "forget [metadata_id]",
		Short: "Forget the chosen metadata",
		Run: func(cmd *cobra.Command, args []string) {
			c := ForgetCommand{
				cache:           cache,
				metadataService: metadataService,
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
	cmd.Flags().StringVarP(&f.MetadataID, "metadata-id", "m", "", "The ID of the metadata to forget")
}

func (cmd *ForgetCommand) Run(ctx context.Context, flags *ForgetFlags) error {
	err := cmd.cache.ValidateForMetadata()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	// if the metadata id is not set, prompt the user for it interactively
	if flags.MetadataID == "" {
		fmt.Fprintf(os.Stdout, "Metadata ID: ")

		_, err := fmt.Scanln(&flags.MetadataID)
		if err != nil {
			return fmt.Errorf("error reading metadata ID: %v", err)
		}
	}

	if flags.MetadataID == "" {
		return fmt.Errorf("no metadata ID provided")
	}

	err = cmd.metadataService.ForgetMetadata(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		flags.MetadataID,
	)
	if err != nil {
		return fmt.Errorf("error forgetting metadata: %v", err)
	}

	// If the metadata was the current metadata in the cache, clear the cache of metadata, and badge IDs
	if cmd.cache.MetadataId == flags.MetadataID {
		cmd.cache.MetadataId = ""
		cmd.cache.BadgeId = ""

		err = clicache.SaveCache(cmd.cache)
		if err != nil {
			return fmt.Errorf("error saving local configuration: %v", err)
		}
	}

	fmt.Fprintf(os.Stdout, "Forgot metadata with ID: %s\n", flags.MetadataID)

	return nil
}
