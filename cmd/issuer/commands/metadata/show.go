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
	mdsrv "github.com/agntcy/identity/internal/issuer/metadata"
	"github.com/agntcy/identity/internal/pkg/cmdutil"
)

type ShowFlags struct {
	MetadataID string
}

type ShowCommand struct {
	cache           *clicache.Cache
	metadataService mdsrv.MetadataService
}

func NewCmdShow(
	cache *clicache.Cache,
	metadataService mdsrv.MetadataService,
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
		return fmt.Errorf("error validating local configuration: %w", err)
	}

	// if the metadata id is not set, prompt the user for it interactively
	err = cmdutil.ScanRequiredIfNotSet("Metadata ID", &flags.MetadataID)
	if err != nil {
		return fmt.Errorf("error reading metadata ID: %w", err)
	}

	metadata, err := cmd.metadataService.GetMetadata(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		flags.MetadataID,
	)
	if err != nil {
		return fmt.Errorf("error getting metadata: %w", err)
	}

	metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling metadata to JSON: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", string(metadataJSON))

	return nil
}
