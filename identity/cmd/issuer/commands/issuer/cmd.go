// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	issuer "github.com/agntcy/identity/internal/issuer/issuer"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

func NewCmd(
	cache *clicache.Cache,
	issuerService issuer.IssuerService,
	vaultSrv vault.VaultService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issuer",
		Short: "Register as an Issuer and manage issuer configurations",
		Long: `
The issuer command is used to register as an Issuer and manage issuer configurations.
`,
	}

	cmd.AddCommand(NewCmdRegister(cache, issuerService, vaultSrv))
	cmd.AddCommand(NewCmdList(cache, issuerService))
	cmd.AddCommand(NewCmdShow(cache, issuerService))
	cmd.AddCommand(NewCmdForget(cache, issuerService))
	cmd.AddCommand(NewCmdLoad(cache, issuerService))

	return cmd
}
