// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "none"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Version   : %s\n", Version)
			cmd.Printf("Build date: %s\n", BuildDate)
			cmd.Printf("Git commit: %s\n", GitCommit)
		},
	}
}
