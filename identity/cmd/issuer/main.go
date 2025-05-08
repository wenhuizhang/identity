// Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/agntcy/identity/cmd/issuer/commands"

	"github.com/spf13/cobra"
)

// We use go generate to copy the web assets into the correct location for the local web server
// This needs to be rerun whenever the web assets are updated
//go:generate cp -r ../../../ui/dist ./web

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "identity",
	Long: `
  _____                           _____ _____
 |  __ \                         |_   _|  __ \
 | |__) |   _ _ __ __ _ _ __ ___   | | | |  | |
 |  ___/ | | | '__/ _' | '_ ' _ \  | | | |  | |
 | |   | |_| | | | (_| | | | | | |_| |_| |__| |
 |_|    \__, |_|  \__,_|_| |_| |_|\___/|_____/
        __/ /
       |___/

The Identity CLI tool is a command line interface for generating and publishing
identities within the Internet of Agents.

With it you can:
- Connect to a local wallet, generate and store quantum-resistant cryptographic keys
- Connect to a network, publish your identity and interact with other agents
- Create and manage your agent identities, including agent artefacts and agent artefact versions
- Verify the identity of other agents via their agent passport
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(commands.WalletCmd)
	rootCmd.AddCommand(commands.NetworkCmd)
	rootCmd.AddCommand(commands.AgentCmd)
	rootCmd.AddCommand(commands.VerifyCmd)
	rootCmd.AddCommand(commands.WebCmd)
}
