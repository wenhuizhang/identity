// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package setup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	issuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// GetIdpConfigPath returns the path to the IDP config file
func GetIdpConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".identity", "idp_config.json"), nil
}

// ReadIdpConfig reads and parses the IDP config file
func ReadIdpConfig() (issuerTypes.IdpConfig, error) {
	var config issuerTypes.IdpConfig

	configPath, err := GetIdpConfigPath()
	if err != nil {
		return config, err
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
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

func ConfigureIdp(idpConfig issuerTypes.IdpConfig) (string, error) {
	configPath, err := GetIdpConfigPath()
	if err != nil {
		return "", err
	}

	// Create directories if they don't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, issuerConstants.DirPerm); err != nil {
		return "", err
	}

	// Marshal the config to JSON
	configData, err := json.Marshal(idpConfig)
	if err != nil {
		return "", err
	}

	// Write the config to file
	if err := os.WriteFile(configPath, configData, issuerConstants.FilePerm); err != nil {
		return "", err
	}

	return configPath, nil
}

func TestIdpConnection() (*oauth2.Token, error) {
	config, err := ReadIdpConfig()
	if err != nil {
		return nil, err
	}

	// Test the connection to the Identity Provider
	ctx := context.Background()

	// Discover OIDC provider config
	provider, err := oidc.NewProvider(ctx, config.IssuerUrl)
	if err != nil {
		return nil, err
	}

	// Set up the OAuth2 client credentials config
	conf := clientcredentials.Config{
		ClientID:     config.ClientId,
		ClientSecret: config.ClientSecret,
		TokenURL:     provider.Endpoint().TokenURL,
		Scopes:       []string{},
	}

	// Retrieve a token
	token, err := conf.Token(ctx)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func ForgetIdpConnection() (string, error) {
	configPath, err := GetIdpConfigPath()
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

	fmt.Fprintf(os.Stdout, "\nSuccessfully deleted Identity Provider configuration at %s\n\n", configPath)

	return configPath, nil
}
