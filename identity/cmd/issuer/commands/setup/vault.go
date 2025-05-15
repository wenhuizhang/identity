// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	"github.com/agntcy/identity/cmd/issuer/commands/setup/vault"
	"github.com/spf13/cobra"
)

var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage your vault and generate cryptographic keys",
	Long: `
The vault command is used to manage your vault and generate cryptographic keys. With it you can:
- (txt) Connect to a local .txt file and generate cryptographic keys
- (1password) Connect to 1Password and generate cryptographic keys
`,
}

func init() {
	VaultCmd.AddCommand(vault.TxtCmd)
	VaultCmd.AddCommand(vault.OnePasswordCmd)
}
