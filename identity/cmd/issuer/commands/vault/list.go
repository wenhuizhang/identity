// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing vault configurations",
	Run: func(cmd *cobra.Command, args []string) {

		vaults, err := vaultService.GetAllVaults()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing vaults: %v\n", err)
			return
		}
		if len(vaults) == 0 {
			fmt.Fprintf(os.Stdout, "No vaults found.\n")
			return
		}
		fmt.Fprintf(os.Stdout, "Existing vaults:\n")
		for _, vault := range vaults {
			fmt.Fprintf(os.Stdout, "- %s (%s vault), id: %s\n", vault.Name, vault.Type, vault.Id)
		}
	},
}
