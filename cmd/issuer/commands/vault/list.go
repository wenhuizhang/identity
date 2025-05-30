// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"fmt"
	"os"

	vaultsrv "github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	vaultService vaultsrv.VaultService
}

func NewCmdList(vaultService vaultsrv.VaultService) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your existing vault configurations",
		Run: func(cmd *cobra.Command, args []string) {
			c := ListCommand{
				vaultService: vaultService,
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
	vaults, err := cmd.vaultService.GetAllVaults()
	if err != nil {
		return fmt.Errorf("error listing vaults: %w", err)
	}

	if len(vaults) == 0 {
		fmt.Fprintf(os.Stdout, "No vaults found.\n")
		return nil
	}

	fmt.Fprintf(os.Stdout, "Existing vaults:\n")

	for _, vault := range vaults {
		fmt.Fprintf(os.Stdout, "- %s (%s vault), id: %s\n", vault.Name, vault.Type, vault.Id)
	}

	return nil
}
