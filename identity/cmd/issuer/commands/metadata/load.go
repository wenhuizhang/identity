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

type LoadFlags struct {
	MetadataID string
}

type LoadCommand struct {
	cache           *clicache.Cache
	metadataService metadata.MetadataService
}

func NewCmdLoad(
	cache *clicache.Cache,
	metadataService metadata.MetadataService,
) *cobra.Command {
	flags := NewLoadFlags()

	cmd := &cobra.Command{
		Use:   "load [metadata_id]",
		Short: "Load a metadata configuration",
		Run: func(cmd *cobra.Command, args []string) {
			c := LoadCommand{
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

func NewLoadFlags() *LoadFlags {
	return &LoadFlags{}
}

func (f *LoadFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.MetadataID, "metadata-id", "m", "", "The ID of the metadata to load")
}

func (cmd *LoadCommand) Run(ctx context.Context, flags *LoadFlags) error {
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

	// check the metadata id is valid
	metadata, err := cmd.metadataService.GetMetadata(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		flags.MetadataID,
	)
	if err != nil {
		return fmt.Errorf("error getting metadata: %v", err)
	}

	if metadata == nil {
		return fmt.Errorf("no metadata found with ID: %s", flags.MetadataID)
	}

	// save the metadata id to the cache
	cmd.cache.MetadataId = metadata.ID

	err = clicache.SaveCache(cmd.cache)
	if err != nil {
		return fmt.Errorf("error saving local configuration: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Loaded metadata with ID: %s\n", flags.MetadataID)

	return nil
}
