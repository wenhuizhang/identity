// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package filesystem

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	internalIssuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault/data"
	"github.com/google/uuid"
)

type vaultFilesystemRepository struct{}

func NewVaultFilesystemRepository() data.VaultRepository {
	return &vaultFilesystemRepository{}
}

// getVaultsDirectory returns the path to the vaults directory
func getVaultsDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".identity", "vaults"), nil
}

// GetVaultIdDirectory returns the path to the vault ID directory
func GetVaultIdDirectory(vaultId string) (string, error) {
	vaultsDir, err := getVaultsDirectory()
	if err != nil {
		return "", err
	}

	return filepath.Join(vaultsDir, vaultId), nil
}

// GetVaultFilePath returns the path to the vault file
func GetVaultFilePath(vaultId string) (string, error) {
	vaultIdDir, err := GetVaultIdDirectory(vaultId)
	if err != nil {
		return "", err
	}

	return filepath.Join(vaultIdDir, "vault.json"), nil
}

func saveVaultConfig() error {

	return nil
}

func (r *vaultFilesystemRepository) ConnectVault() (*internalIssuerTypes.Vault, error) {

	// Save the vault config
	if err := saveVaultConfig(); err != nil {
		return nil, err
	}

	vaultId := uuid.New().String()

	vault := internalIssuerTypes.Vault{
		Id: vaultId,
	}

	// Create idp locally in the vault directory
	vaultsDir, err := GetVaultIdDirectory(vaultId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(vaultsDir, internalIssuerConstants.DirPerm); err != nil {
		return nil, err
	}

	vaultFilePath, err := GetVaultFilePath(vaultId)
	if err != nil {
		return nil, err
	}

	// Marshal the config to JSON
	vaultData, err := json.Marshal(&vault)
	if err != nil {
		return nil, err
	}

	// Write the vault to file
	if err := os.WriteFile(vaultFilePath, vaultData, internalIssuerConstants.FilePerm); err != nil {
		return nil, err
	}

	return &vault, nil
}

func (r *vaultFilesystemRepository) ListVaultIds() ([]string, error) {
	// Get the vaults directory
	vaultsDir, err := getVaultsDirectory()
	if err != nil {
		return nil, err
	}

	// Read the vaults directory
	files, err := os.ReadDir(vaultsDir)
	if err != nil {
		return nil, err
	}

	// List the vault IDs
	var vaultIds []string

	for _, file := range files {
		if file.IsDir() {
			vaultIds = append(vaultIds, file.Name())
		}
	}

	return vaultIds, nil
}

func (r *vaultFilesystemRepository) GetVault(vaultId string) (*internalIssuerTypes.Vault, error) {
	// Get the vault file path
	vaultFilePath, err := GetVaultFilePath(vaultId)
	if err != nil {
		return nil, err
	}

	// Read the vault file
	vaultData, err := os.ReadFile(vaultFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the vault data
	var vault internalIssuerTypes.Vault
	if err := json.Unmarshal(vaultData, &vault); err != nil {
		return nil, err
	}

	return &vault, nil
}

func (r *vaultFilesystemRepository) ForgetVault(vaultId string) error {
	// Get the vault directory
	vaultDir, err := GetVaultIdDirectory(vaultId)
	if err != nil {
		return err
	}

	// Check if the vault directory exists
	if _, err := os.Stat(vaultDir); os.IsNotExist(err) {
		return errors.New("Vault does not exist")
	}

	// Remove the vault directory
	if err := os.RemoveAll(vaultDir); err != nil {
		return err
	}

	return nil
}
