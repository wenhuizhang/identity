// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/metadata"
)

type ShowFlags struct {
	MetadataID string
}

type ShowCommand struct {
	cache           *clicache.Cache
	metadataService metadata.MetadataService
}

func NewCmdShow(
	cache *clicache.Cache,
	metadataService metadata.MetadataService,
) *cobra.Command {
	flags := NewShowFlags()

	cmd := &cobra.Command{
		Use:   "show [metadata_id]",
		Short: "Show the chosen metadata",
		Run: func(cmd *cobra.Command, args []string) {
			c := ShowCommand{
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

func NewShowFlags() *ShowFlags {
	return &ShowFlags{}
}

func (f *ShowFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.MetadataID, "metadata-id", "m", "", "The ID of the metadata to show")
}

func (cmd *ShowCommand) Run(ctx context.Context, flags *ShowFlags) error {
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

	metadata, err := cmd.metadataService.GetMetadata(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		flags.MetadataID,
	)
	if err != nil {
		return fmt.Errorf("error getting metadata: %v", err)
	}

	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling metadata to JSON: %v", err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", string(metadataJSON))

	return nil
}
