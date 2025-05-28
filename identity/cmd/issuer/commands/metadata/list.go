// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	mdsrv "github.com/agntcy/identity/internal/issuer/metadata"
)

type ListCommand struct {
	cache           *clicache.Cache
	metadataService mdsrv.MetadataService
}

func NewCmdList(
	cache *clicache.Cache,
	metadataService mdsrv.MetadataService,
) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your existing metadata",
		Run: func(cmd *cobra.Command, args []string) {
			c := ListCommand{
				cache:           cache,
				metadataService: metadataService,
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
	err := cmd.cache.ValidateForMetadata()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %w", err)
	}

	allMetadata, err := cmd.metadataService.GetAllMetadata(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
	)
	if err != nil {
		return fmt.Errorf("error listing metadata: %w", err)
	}

	if len(allMetadata) == 0 {
		fmt.Fprintf(os.Stdout, "%s\n", "No metadata found")
	}

	fmt.Fprintf(os.Stdout, "%s\n", "Existing metadata ids:")

	for _, metadata := range allMetadata {
		fmt.Fprintf(os.Stdout, "- %s\n", metadata.ID)
	}

	return nil
}
