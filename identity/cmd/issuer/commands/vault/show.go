// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	vaultsrv "github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type ShowFlags struct {
	VaultID string
}

type ShowCommand struct {
	vaultService vaultsrv.VaultService
}

func NewCmdShow(vaultService vaultsrv.VaultService) *cobra.Command {
	flags := NewShowFlags()

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show details of a vault configuration",
		Run: func(cmd *cobra.Command, args []string) {
			c := ShowCommand{
				vaultService: vaultService,
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
	cmd.Flags().StringVarP(&f.VaultID, "vault-id", "v", "", "The ID of the vault to show")
}

func (cmd *ShowCommand) Run(ctx context.Context, flags *ShowFlags) error {
	// if the vault id is not set, prompt the user for it interactively
	if flags.VaultID == "" {
		err := cmdutil.ScanRequired("Vault ID to show", &flags.VaultID)
		if err != nil {
			return fmt.Errorf("error reading vault ID: %w", err)
		}
	}

	// check the vault id is valid
	vault, err := cmd.vaultService.GetVault(flags.VaultID)
	if err != nil {
		return fmt.Errorf("error getting vault: %w", err)
	}

	if vault == nil {
		return fmt.Errorf("no vault found with ID: %s", flags.VaultID)
	}

	vaultJSON, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling metadata to JSON: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", string(vaultJSON))

	return nil
}
