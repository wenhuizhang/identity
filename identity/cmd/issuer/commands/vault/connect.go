// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"github.com/spf13/cobra"
)

var vaultConnectCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new vault configuration and generate cryptographic keys",
	Long: `
The create command is used to create a new vault configuration and generate cryptographic keys.
`,
}
