// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Load and verify an Agent Badge",
}

var loadCmd = &cobra.Command{
	Use:   "load [agent_badge]",
	Short: "Load an Agent Badge",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Printf("Loading agent badge")
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the loaded Agent Badge",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Validating the loaded Agent Badge")
	},
}

var forgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the loaded Agent Badge",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Forgetting the loaded Agent Badge")
	},
}

func init() {
	VerifyCmd.AddCommand(loadCmd)
	VerifyCmd.AddCommand(validateCmd)
	VerifyCmd.AddCommand(forgetCmd)
}
