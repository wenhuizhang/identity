// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/spf13/cobra"
)

type GenerateCmdInput struct {
	KeyID string
}

type ShowCmdInput struct {
	KeyID string
}

type LoadCmdInput struct {
	KeyID string
}

var (
	// setup the vault service
	vaultFilesystemRepository = filesystem.NewVaultFilesystemRepository()
	vaultService              = vault.NewVaultService(vaultFilesystemRepository)

	// setup the vault command flags
	generateCmdIn = &GenerateCmdInput{}
	showCmdIn     = &ShowCmdInput{}
	loadCmdIn     = &LoadCmdInput{}
)

var KeyCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage cryptographic keys for vaults",
	Long: `
The keys command is used to generate and manage cryptographic keys for in your vault.
`,
}

func init() {

	keyGenerateCmd.Flags().StringVarP(&generateCmdIn.KeyID, "key-id", "k", "", "The ID of the key to generate")
	KeyCmd.AddCommand(keyGenerateCmd)

	KeyCmd.AddCommand(keyListCmd)

	keyShowCmd.Flags().StringVarP(&showCmdIn.KeyID, "key-id", "k", "", "The ID of the key to show")
	KeyCmd.AddCommand(keyShowCmd)

	keyLoadCmd.Flags().StringVarP(&loadCmdIn.KeyID, "key-id", "k", "", "The ID of the key to load")
	KeyCmd.AddCommand(keyLoadCmd)
}
