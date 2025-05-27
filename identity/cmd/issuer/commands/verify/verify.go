// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package verify

import (
	"encoding/json"
	"fmt"
	"os"

	coreV1alpha "github.com/agntcy/identity/api/server/agntcy/identity/core/v1alpha1"
	issuerVerify "github.com/agntcy/identity/internal/issuer/verify"

	"github.com/spf13/cobra"
)

var (
	badgeFilePath string
)

var VerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify an Agent or MCP Server Badge from a file",
	Run: func(cmd *cobra.Command, args []string) {

		// if the file path is not set, prompt the user for it interactively
		if badgeFilePath == "" {
			fmt.Fprintf(os.Stderr, "Full file path to the badge file: ")
			_, err := fmt.Scanln(&badgeFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file path: %v\n", err)
				return
			}
		}
		if badgeFilePath == "" {
			fmt.Fprintf(os.Stderr, "No file path provided\n")
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
		var vc coreV1alpha.VerifiableCredential
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
		fmt.Fprintf(os.Stdout, "Successfully verified badge: %s\n", *vc.Id)

	},
}

func init() {
	VerifyCmd.Flags().StringVarP(&badgeFilePath, "file", "f", "", "Path to the badge file")
}
