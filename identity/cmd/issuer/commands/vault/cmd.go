// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/cmd/issuer/commands/vault/connect"
	"github.com/agntcy/identity/cmd/issuer/commands/vault/key"
	vaultsrv "github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

func NewCmd(
	cache *clicache.Cache,
	vaultService vaultsrv.VaultService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Manage your vaults and generate cryptographic keys",
		Long: `
The vault command is used to configure and manage your vaults.
`,
	}

	cmd.AddCommand(connect.NewCmd(vaultService))
	cmd.AddCommand(NewCmdList(vaultService))
	cmd.AddCommand(NewCmdShow(vaultService))
	cmd.AddCommand(NewCmdForget(vaultService))
	cmd.AddCommand(NewCmdLoad(vaultService))
	cmd.AddCommand(key.NewCmd(cache, vaultService))

	return cmd
}
