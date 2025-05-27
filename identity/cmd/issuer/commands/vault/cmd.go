// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/spf13/cobra"
)

type ShowCmdInput struct {
	VaultID string
}

type ForgetCmdInput struct {
	VaultID string
}

type LoadCmdInput struct {
	VaultID string
}

var (
	// setup the vault service
	vaultFilesystemRepository = filesystem.NewVaultFilesystemRepository()
	vaultService              = vault.NewVaultService(vaultFilesystemRepository)

	// setup the vault command flags
	showCmdIn   = &ShowCmdInput{}
	forgetCmdIn = &ForgetCmdInput{}
	loadCmdIn   = &LoadCmdInput{}
)

var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage your vaults and generate cryptographic keys",
	Long: `
The vault command is used to configure and manage your vaults.
`,
}

func init() {
	// Add the vault types to the vault connect command
	vaultConnectCmd.AddCommand(TxtCmd)

	VaultCmd.AddCommand(vaultConnectCmd)

	VaultCmd.AddCommand(vaultListCmd)

	vaultShowCmd.Flags().StringVarP(&showCmdIn.VaultID, "vault-id", "v", "", "The ID of the vault to show")
	VaultCmd.AddCommand(vaultShowCmd)

	vaultForgetCmd.Flags().StringVarP(&forgetCmdIn.VaultID, "vault-id", "v", "", "The ID of the vault to forget")
	VaultCmd.AddCommand(vaultForgetCmd)

	vaultLoadCmd.Flags().StringVarP(&loadCmdIn.VaultID, "vault-id", "v", "", "The ID of the vault to load")
	VaultCmd.AddCommand(vaultLoadCmd)
}
