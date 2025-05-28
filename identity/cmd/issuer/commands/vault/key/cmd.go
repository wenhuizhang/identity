// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package key

import (
	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

func NewCmd(
	cache *clicache.Cache,
	vaultService vault.VaultService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Manage cryptographic keys for vaults",
		Long: `
The keys command is used to generate and manage cryptographic keys for in your vault.
`,
	}

	cmd.AddCommand(NewCmdGenerate(cache, vaultService))
	cmd.AddCommand(NewCmdList(cache, vaultService))
	cmd.AddCommand(NewCmdShow(cache, vaultService))
	cmd.AddCommand(NewCmdLoad(cache, vaultService))

	return cmd
}
