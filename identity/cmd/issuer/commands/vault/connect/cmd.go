// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package connect

import (
	vaultsrv "github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

func NewCmd(vaultService vaultsrv.VaultService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect a new vault configuration to generate and store cryptographic keys",
		Long: `
The connect command is used to connect a new vault configuration to generate and store cryptographic keys.
`,
	}

	cmd.AddCommand(NewCmdFile(vaultService))
	cmd.AddCommand(NewCmdHashicorp(vaultService))

	return cmd
}
