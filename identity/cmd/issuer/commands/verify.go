// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var VerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Load and verify an Agent or MCP Server Badge",
	Long: `
The verify command is used to load and verify an Agent or MCP Server Badge. With it you can:

- (load) Load an existing badge
- (validate) Validate the loaded badge
- (forget) Forget the loaded badge
`,
}

var loadCmd = &cobra.Command{
	Use:   "load [badge]",
	Short: "Load an Agent of MCP Server Badge",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Loading Badge")
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the loaded Agent or MCP Server Badge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Validating the loaded Badge")
	},
}

var forgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the loaded Agent or MCP Server Badge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Forgetting the loaded Badge")
	},
}

func init() {
	VerifyCmd.AddCommand(loadCmd)
	VerifyCmd.AddCommand(validateCmd)
	VerifyCmd.AddCommand(forgetCmd)
}
