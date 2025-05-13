// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package vault

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/1password/onepassword-sdk-go"
	"github.com/agntcy/identity/cmd/issuer/keys"
	"github.com/spf13/cobra"
)

// Constants for file and directory permissions
const (
	filePerm = 0o600
	dirPerm  = 0o700
)

var OnePasswordCmd = &cobra.Command{
	Use:   "1password",
	Short: "Connect to 1Password",
	Long:  "Connect to 1Password",
}

var onePasswordConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to your 1Password account. ",
	Run: func(cmd *cobra.Command, args []string) {

		// request the user to enter their 1Password service account token
		var serviceAccountToken string
		fmt.Fprintf(os.Stdout, "%s\n", "\nTo connect to your 1Password account you will need to have created")
		fmt.Fprintf(os.Stdout, "%s\n", "a service account token with write permissions to at least one vault:")
		fmt.Fprintf(os.Stdout, "%s\n", "https://my.1password.com/developer-tools/infrastructure-secrets/serviceaccount/")
		fmt.Fprintf(os.Stdout, "%s\n", "\nEnter your 1Password service account token: ")
		_, err := fmt.Scanln(&serviceAccountToken)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error reading service account token: %v\n\n", err)
			return
		}

		// Validate the service account token
		if serviceAccountToken == "" {
			fmt.Fprintf(os.Stdout, "%s\n", "Service account token cannot be empty.")
			return
		}

		// Test the connection to 1Password
		client, err := onepassword.NewClient(
			context.TODO(),
			onepassword.WithServiceAccountToken(serviceAccountToken),
			onepassword.WithIntegrationInfo("Agntcy Identity 1Password Integration", "v0.0.1"),
		)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error connecting to 1Password: %v\n\n", err)
			return
		}

		// Check if the service account token is valid and has access to at least one vault
		vaults, err := client.Vaults().List(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error listing vaults: %v\n\n", err)
			return
		}
		if len(vaults) == 0 {
			fmt.Fprintf(os.Stdout, "No vaults found for the provided service account token.\n\n")
			return
		}

		// List all vaults
		fmt.Fprintf(os.Stdout, "%s\n", "\nFound vaults:")
		for _, vault := range vaults {
			fmt.Fprintf(os.Stdout, "- Name: %s, Vault ID: %s\n", vault.Title, vault.ID)
		}

		// if the user has multiple vaults, ask them to select one
		fmt.Fprintf(os.Stdout, "%s\n", "\nPlease select a vault by entering its ID:")
		var selectedVaultID string
		_, err = fmt.Scanln(&selectedVaultID)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error reading vault ID: %v\n\n", err)
			return
		}
		if selectedVaultID == "" {
			fmt.Fprintf(os.Stdout, "%s\n", "Vault ID cannot be empty.")
			return
		}

		// Validate the selected vault ID
		vaultFound := false
		for _, vault := range vaults {
			if vault.ID == selectedVaultID {
				vaultFound = true
				break
			}
		}
		if !vaultFound {
			fmt.Fprintf(os.Stdout, "Vault with ID %s not found.\n\n", selectedVaultID)
			return
		}

		// Create a config struct to hold the service account token and selected vault ID
		config := struct {
			ServiceAccountToken string `json:"serviceAccountToken"`
			VaultID             string `json:"vaultId"`
		}{
			ServiceAccountToken: serviceAccountToken,
			VaultID:             selectedVaultID,
		}

		// Save to a config file in the user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error getting home directory: %v\n\n", err)
			return
		}
		configPath := filepath.Join(homeDir, ".identity", "1password_config.json")

		// Create directories if they don't exist
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, dirPerm); err != nil {
			fmt.Fprintf(os.Stdout, "Error creating config directory: %v\n\n", err)
			return
		}

		// Marshal the config to JSON
		configData, err := json.Marshal(config)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error marshaling config: %v\n\n", err)
			return
		}

		// Write the config to file
		if err := os.WriteFile(configPath, configData, filePerm); err != nil {
			fmt.Fprintf(os.Stdout, "Error writing config file: %v\n\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nSuccessfully connected to 1Password vault.\n")
		fmt.Fprintf(os.Stdout, "\nConfiguration saved to %s\n\n", configPath)
	},
}

var onePasswordForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the current 1Password account and vault by deleting the config file",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stdout, "\nError getting home directory: %v\n", err)
			return
		}
		configPath := filepath.Join(homeDir, ".identity", "1password_config.json")

		// Check if the config file exists
		if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stdout, "\nNo 1Password configuration found at %s\n\n", configPath)
			return
		}

		// Delete the config file
		if err := os.Remove(configPath); err != nil {
			fmt.Fprintf(os.Stdout, "\nError deleting config file: %v\n\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nSuccessfully deleted 1Password configuration at %s\n\n", configPath)
	},
}

var onePasswordGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate keys and store them in your vault",
	Long:  `Generate keys and store them in your vault`,
	Run: func(cmd *cobra.Command, args []string) {

		// Generate public and private keys
		// Load the 1Password configuration
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error getting home directory: %v\n", err)
			return
		}
		configPath := filepath.Join(homeDir, ".identity", "1password_config.json")

		// Check if the config file exists
		configData, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s\n", "1Password configuration not found. Please run 'connect' command first.\n")
			return
		}

		// Parse the config
		var config struct {
			ServiceAccountToken string `json:"serviceAccountToken"`
			VaultID             string `json:"vaultId"`
		}
		if err := json.Unmarshal(configData, &config); err != nil {
			fmt.Fprintf(os.Stdout, "Error parsing config: %v\n", err)
			return
		}

		// Create the 1Password client
		client, err := onepassword.NewClient(
			context.Background(),
			onepassword.WithServiceAccountToken(config.ServiceAccountToken),
			onepassword.WithIntegrationInfo("Agntcy Identity 1Password Integration", "v0.0.1"),
		)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error connecting to 1Password: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "%s\n", "\nChecking if the vault already contains AGNTCY Identity keys...")

		// Define the item title
		itemTitle := "AGNTCY Identity"

		// Check if the item already exists
		items, err := client.Items().List(context.Background(), config.VaultID)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error listing items: %v\n", err)
			return
		}

		for _, item := range items {
			if item.Title == itemTitle {
				// Updated error messages:
				fmt.Fprintf(os.Stdout, "\nItem with the name 'AGNTCY Identity' already exists in the vault.\n")
				fmt.Fprintf(os.Stdout, "Name: %s, Item ID: %s\n\n", item.Title, item.ID)
				fmt.Fprintf(os.Stdout, "You can use the 'load' command with the item ID to use existing keys in your vault.\n\n")
				return
			}
		}

		fmt.Fprintf(os.Stdout, "%s\n", "\nGenerating and storing keys...")

		// Generate ECDSA key pair
		privateKeyPEM, publicKeyPEM := keys.GenerateECDSAKeyPair()
		if privateKeyPEM == nil {
			fmt.Fprintf(os.Stderr, "Error generating private key.")
			return
		}
		if publicKeyPEM == nil {
			fmt.Fprintf(os.Stderr, "Error generating public key.")
			return
		}

		// Define the item parameters
		itemParams := onepassword.ItemCreateParams{
			Title:    itemTitle,
			Category: onepassword.ItemCategorySecureNote,
			VaultID:  config.VaultID,
			Fields: []onepassword.ItemField{
				{
					ID:        "private_key",
					Title:     "Private Key",
					Value:     string(privateKeyPEM),
					FieldType: onepassword.ItemFieldTypeConcealed,
				},
				{
					ID:        "public_key",
					Title:     "Public Key",
					Value:     string(publicKeyPEM),
					FieldType: onepassword.ItemFieldTypeText,
				},
			},
			Tags: []string{"identity", "agent-keys"},
		}

		// Creates a new item based on the structure definition
		item, err := client.Items().Create(context.Background(), itemParams)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error storing keys in 1Password: %v\n", err)
			return
		}

		// Save the Item ID to a file to the existing config file
		configFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, filePerm)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error opening config file: %v\n", err)
			return
		}
		defer configFile.Close()
		// Read the existing config
		var existingConfig struct {
			ServiceAccountToken string `json:"serviceAccountToken"`
			VaultID             string `json:"vaultId"`
			ItemID              string `json:"itemId"`
		}
		if err := json.Unmarshal(configData, &existingConfig); err != nil {
			fmt.Fprintf(os.Stdout, "Error parsing existing config: %v\n", err)
			return
		}
		// Update the ItemID in the config
		existingConfig.ItemID = item.ID
		// Marshal the updated config to JSON
		updatedConfigData, err := json.Marshal(existingConfig)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error marshaling updated config: %v\n", err)
			return
		}
		// Write the updated config to file
		if _, err := configFile.WriteAt(updatedConfigData, 0); err != nil {
			fmt.Fprintf(os.Stdout, "Error writing updated config file: %v\n", err)
			return
		}
		if err := configFile.Truncate(int64(len(updatedConfigData))); err != nil {
			fmt.Fprintf(os.Stdout, "Error truncating config file: %v\n", err)
			return
		}
		if err := configFile.Close(); err != nil {
			fmt.Fprintf(os.Stdout, "Error closing config file: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "%s\n", "\nSuccessfully generated and stored keys in 1Password vault.")
		fmt.Fprintf(os.Stdout, "Name: %s, Item ID: %s\n\n", item.Title, item.ID)
		fmt.Fprintf(os.Stdout, "Configuration saved to %s\n\n", configPath)

	},
}

var onePasswordLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a key from your 1Password vault",
	Long:  `Load a key from your 1Password vault`,
	Run: func(cmd *cobra.Command, args []string) {

		// Load the 1Password configuration
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error getting home directory: %v\n", err)
			return
		}
		configPath := filepath.Join(homeDir, ".identity", "1password_config.json")

		// Check if the config file exists
		configData, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s\n", "1Password configuration not found. Please run 'connect' command first.\n")
			return
		}

		// Parse the config
		var config struct {
			ServiceAccountToken string `json:"serviceAccountToken"`
			VaultID             string `json:"vaultId"`
		}
		if err := json.Unmarshal(configData, &config); err != nil {
			fmt.Fprintf(os.Stdout, "Error parsing config: %v\n", err)
			return
		}

		client, err := onepassword.NewClient(
			context.Background(),
			onepassword.WithServiceAccountToken(config.ServiceAccountToken),
			onepassword.WithIntegrationInfo("Agntcy Identity 1Password Integration", "v0.0.1"),
		)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error connecting to 1Password: %v\n", err)
			return
		}

		// Get the available items in the vault
		items, err := client.Items().List(context.Background(), config.VaultID)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error listing items: %v\n", err)
			return
		}

		if len(items) == 0 {
			fmt.Fprintf(os.Stdout, "No items found in the vault.\n\n")
			return
		}

		// List all items in the vault
		fmt.Fprintf(os.Stdout, "%s\n", "\nAvailable items in the vault:")
		for _, item := range items {
			fmt.Fprintf(os.Stdout, "- Name: %s, Item ID: %s\n", item.Title, item.ID)
		}

		fmt.Fprintf(os.Stdout, "%s\n", "\nPlease enter the ID of the item you want to load:")
		var selectedItemID string
		_, err = fmt.Scanln(&selectedItemID)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error reading item ID: %v\n", err)
			return
		}
		if selectedItemID == "" {
			fmt.Fprintf(os.Stdout, "%s\n", "Item ID cannot be empty.")
			return
		}

		// Validate the selected item ID
		itemFound := false
		for _, item := range items {
			if item.ID == selectedItemID {
				itemFound = true
				break
			}
		}
		if !itemFound {
			fmt.Fprintf(os.Stdout, "Item with ID %s not found.\n\n", selectedItemID)
			return
		}

		// Update the existing config with the selected item ID
		configFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, filePerm)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error opening config file: %v\n", err)
			return
		}
		defer configFile.Close()
		// Read the existing config
		var existingConfig struct {
			ServiceAccountToken string `json:"serviceAccountToken"`
			VaultID             string `json:"vaultId"`
			ItemID              string `json:"itemId"`
		}
		if err := json.Unmarshal(configData, &existingConfig); err != nil {
			fmt.Fprintf(os.Stdout, "Error parsing existing config: %v\n", err)
			return
		}
		// Update the ItemID in the config
		existingConfig.ItemID = selectedItemID
		// Marshal the updated config to JSON
		updatedConfigData, err := json.Marshal(existingConfig)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error marshaling updated config: %v\n", err)
			return
		}
		// Write the updated config to file
		if _, err := configFile.WriteAt(updatedConfigData, 0); err != nil {
			fmt.Fprintf(os.Stdout, "Error writing updated config file: %v\n", err)
			return
		}
		if err := configFile.Truncate(int64(len(updatedConfigData))); err != nil {
			fmt.Fprintf(os.Stdout, "Error truncating config file: %v\n", err)
			return
		}
		if err := configFile.Close(); err != nil {
			fmt.Fprintf(os.Stdout, "Error closing config file: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nSuccessfully loaded keys from 1Password vault.\n")
		fmt.Fprintf(os.Stdout, "\nConfiguration saved to %s\n\n", configPath)

	},
}

func init() {
	OnePasswordCmd.AddCommand(onePasswordConnectCmd)
	OnePasswordCmd.AddCommand(onePasswordForgetCmd)
	OnePasswordCmd.AddCommand(onePasswordGenerateCmd)
	OnePasswordCmd.AddCommand(onePasswordLoadCmd)
}
