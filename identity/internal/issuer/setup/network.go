// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package setup

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	issuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

func ConfigureNetwork(identityNodeConfig issuerTypes.IdentityNodeConfig) (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(homeDir, ".identity", "id_node_config.json")

	// Create directories if they don't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, issuerConstants.DirPerm); err != nil {
		return "", err
	}

	// Marshal the config to JSON
	configData, err := json.Marshal(identityNodeConfig)
	if err != nil {
		return "", err
	}

	// Write the config to file
	if err := os.WriteFile(configPath, configData, issuerConstants.FilePerm); err != nil {
		return "", err
	}

	return configPath, nil
}

func TestNetworkConnection() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(homeDir, ".identity", "id_node_config.json")
	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}
	// Read the config file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	// Unmarshal the config data
	var config issuerTypes.IdentityNodeConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return err
	}

	// TODO: Implement the connection test logic

	return errors.New("TestNetworkConnection not implemented yet")

}

func ForgetNetworkConnection() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(homeDir, ".identity", "id_node_config.json")

	// Check if the config file exists
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	// Delete the config file
	if err := os.Remove(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}
