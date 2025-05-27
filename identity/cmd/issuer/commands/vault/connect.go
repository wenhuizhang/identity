// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"github.com/spf13/cobra"
)

var vaultConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect a new vault configuration to generate and store cryptographic keys",
	Long: `
The connect command is used to connect a new vault configuration to generate and store cryptographic keys.
`,
}
