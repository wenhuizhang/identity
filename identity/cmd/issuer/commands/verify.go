// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"encoding/json"
	"fmt"
	"os"

	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	issuerVerify "github.com/agntcy/identity/internal/issuer/verify"

	"github.com/spf13/cobra"
)

var VerifyCmd = &cobra.Command{
	Use:   "verify [path_to_badge_json]",
	Short: "Verify an Agent or MCP Server Badge from a JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		badgeFilePath := args[0]
		if _, err := os.Stat(badgeFilePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "File does not exist: %s\n", badgeFilePath)
			return
		}

		// Check if the badge file exists
		if _, err := os.Stat(badgeFilePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "File does not exist: %s\n", badgeFilePath)
			return
		}

		// Read the badge file
		vcData, err := os.ReadFile(badgeFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			return
		}

		// Unmarshal the badge data
		var vc vctypes.VerifiableCredential
		if err := json.Unmarshal(vcData, &vc); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshaling badge data: %v\n", err)
			return
		}

		// Verify the badge
		_, err = issuerVerify.VerifyCredential(&vc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error verifying badge: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Successfully verified badge: %s\n", vc.ID)

	},
}
