// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package vaults

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var TxtCmd = &cobra.Command{
	Use:   "txt",
	Short: "Connect to .txt file",
	Long:  "Connect to .txt file",
}

var txtConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an existing .txt file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Connecting to local .txt file")
	},
}

var txtForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the current .txt file by deleting the config file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Forgetting the current .txt file by deleting the config file")
	},
}

var txtGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate keys and store them in your .txt file",
	Long:  `Generate keys and store them in your .txt file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Generating keys and storing them in your .txt file")
	},
}

var txtLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a key from your .txt file",
	Long:  `Load a key from your .txt file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Loading a key from your .txt file")
	},
}

func init() {
	TxtCmd.AddCommand(txtConnectCmd)
	TxtCmd.AddCommand(txtForgetCmd)
	TxtCmd.AddCommand(txtGenerateCmd)
	TxtCmd.AddCommand(txtLoadCmd)
}
