// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/metadata"
	"github.com/spf13/cobra"
)

func NewCmd(
	cache *clicache.Cache,
	metadataService metadata.MetadataService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Generate important metadata for your Agent and MCP Server identities",
		Long: `
The metadata command is used to generate important metadata for your Agent and MCP Server identities.
`,
	}

	cmd.AddCommand(NewCmdGenerate(cache, metadataService))
	cmd.AddCommand(NewCmdList(cache, metadataService))
	cmd.AddCommand(NewCmdShow(cache, metadataService))
	cmd.AddCommand(NewCmdForget(cache, metadataService))
	cmd.AddCommand(NewCmdLoad(cache, metadataService))

	return cmd
}
