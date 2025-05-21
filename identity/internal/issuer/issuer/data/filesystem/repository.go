// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package filesystem

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	internalIssuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	"github.com/agntcy/identity/internal/issuer/issuer/data"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	vaultFilesystemRepository "github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
)

type issuerFilesystemRepository struct{}

func NewIssuerFilesystemRepository() data.IssuerRepository {
	return &issuerFilesystemRepository{}
}

// getIssuersDirectory returns the path to the issuers directory
func getIssuersDirectory(vaultId string) (string, error) {
	vaultIdDir, err := vaultFilesystemRepository.GetVaultIdDirectory(vaultId)
	if err != nil {
		return "", err
	}

	return filepath.Join(vaultIdDir, "issuers"), nil
}

// GetIssuerIdDirectory returns the path to the issuer ID directory
func GetIssuerIdDirectory(vaultId, issuerId string) (string, error) {
	issuersDir, err := getIssuersDirectory(vaultId)
	if err != nil {
		return "", err
	}

	return filepath.Join(issuersDir, issuerId), nil
}

// GetIssuerFilePath returns the path to the issuer file
func GetIssuerFilePath(vaultId, issuerId string) (string, error) {
	issuerIdDir, err := GetIssuerIdDirectory(vaultId, issuerId)
	if err != nil {
		return "", err
	}

	return filepath.Join(issuerIdDir, "issuer.json"), nil
}

func saveIssuerConfig(vaultId, identityNodeAddress string, idpConfig internalIssuerTypes.IdpConfig) error {
	// Get the issuer ID directory
	issuerIdDir, err := GetIssuerIdDirectory(vaultId, idpConfig.ClientId)
	if err != nil {
		return err
	}

	// Create the issuer ID directory if it doesn't exist
	if err := os.MkdirAll(issuerIdDir, internalIssuerConstants.DirPerm); err != nil {
		return err
	}

	// Create the issuer config
	issuerConfig := internalIssuerTypes.IssuerConfig{
		IdentityNodeConfig: &internalIssuerTypes.IdentityNodeConfig{
			IdentityNodeAddress: identityNodeAddress,
		},
		IdpConfig: &idpConfig,
	}

	// Marshal the config to JSON
	configData, err := json.Marshal(issuerConfig)
	if err != nil {
		return err
	}

	// Write the config to file
	configFilePath := filepath.Join(issuerIdDir, "idp_config.json")
	if err := os.WriteFile(configFilePath, configData, internalIssuerConstants.FilePerm); err != nil {
		return err
	}

	return nil
}

func (r *issuerFilesystemRepository) AddIssuer(
	vaultId, identityNodeAddress string, idpConfig internalIssuerTypes.IdpConfig, issuer *coreV1alpha.Issuer,
) (string, error) {

	// Create idp locally in the issuer directory
	issuersDir, err := GetIssuerIdDirectory(vaultId, idpConfig.ClientId)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(issuersDir, internalIssuerConstants.DirPerm); err != nil {
		return "", err
	}

	issuerFilePath, err := GetIssuerFilePath(vaultId, idpConfig.ClientId)
	if err != nil {
		return "", err
	}

	// Marshal the config to JSON
	issuerData, err := json.Marshal(&issuer)
	if err != nil {
		return "", err
	}

	// Write the issuer to file
	if err := os.WriteFile(issuerFilePath, issuerData, internalIssuerConstants.FilePerm); err != nil {
		return "", err
	}

	return idpConfig.ClientId, nil
}

func (r *issuerFilesystemRepository) GetAllIssuers(vaultId string) ([]*internalIssuerTypes.Issuer, error) {
	// Get the issuers directory
	issuersDir, err := getIssuersDirectory(vaultId)
	if err != nil {
		return nil, err
	}

	// Read the issuers directory
	files, err := os.ReadDir(issuersDir)
	if err != nil {
		return nil, err
	}

	// List the issuer IDs
	var issuers []*internalIssuerTypes.Issuer

	for _, file := range files {
		if file.IsDir() {
			coreIssuer, err := r.GetIssuer(vaultId, file.Name())
			if err != nil {
				return nil, err
			}
			// Append the issuer to the list
			issuers = append(issuers, &internalIssuerTypes.Issuer{
				Id:     file.Name(),
				Issuer: coreIssuer,
			})
		}
	}

	return issuers, nil
}

func (r *issuerFilesystemRepository) GetIssuer(vaultId, issuerId string) (*coreV1alpha.Issuer, error) {
	// Get the issuer file path
	issuerFilePath, err := GetIssuerFilePath(vaultId, issuerId)
	if err != nil {
		return nil, err
	}

	// Read the issuer file
	issuerData, err := os.ReadFile(issuerFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the issuer data
	var issuer coreV1alpha.Issuer
	if err := json.Unmarshal(issuerData, &issuer); err != nil {
		return nil, err
	}

	return &issuer, nil
}

func (r *issuerFilesystemRepository) RemoveIssuer(vaultId, issuerId string) error {
	// Get the issuer directory
	issuerDir, err := GetIssuerIdDirectory(vaultId, issuerId)
	if err != nil {
		return err
	}

	// Check if the issuer directory exists
	if _, err := os.Stat(issuerDir); os.IsNotExist(err) {
		return errors.New("issuer does not exist")
	}

	// Remove the issuer directory
	if err := os.RemoveAll(issuerDir); err != nil {
		return err
	}

	return nil
}
