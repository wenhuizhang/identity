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

// GetConfigPath returns the path to the identity node config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".identity", "id_node_config.json"), nil
}

// ReadNetworkConfig reads and parses the identity node config file
func ReadNetworkConfig() (issuerTypes.IdentityNodeConfig, error) {
	var config issuerTypes.IdentityNodeConfig

	configPath, err := GetConfigPath()
	if err != nil {
		return config, err
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, err
	}

	// Read the config file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	// Unmarshal the config data
	if err := json.Unmarshal(configData, &config); err != nil {
		return config, err
	}

	return config, nil
}

func ConfigureNetwork(identityNodeConfig issuerTypes.IdentityNodeConfig) (string, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return "", err
	}

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
	_, err := ReadNetworkConfig()
	if err != nil {
		return err
	}

	return errors.New("TestNetworkConnection not implemented yet")
}

func ForgetNetworkConnection() (string, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return "", err
	}

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
