// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package filesystem

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/agntcy/identity/internal/issuer/metadata/data"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	internalIssuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	issuerFilesystemRepository "github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type metadataFilesystemRepository struct{}

func NewMetadataFilesystemRepository() data.MetadataRepository {
	return &metadataFilesystemRepository{}
}

func getMetadataDirectory(vaultId, issuerId string) (string, error) {
	issuerIdDir, err := issuerFilesystemRepository.GetIssuerIdDirectory(vaultId, issuerId)
	if err != nil {
		return "", err
	}

	return filepath.Join(issuerIdDir, "metadata"), nil
}

func GetMetadataIdDirectory(vaultId, issuerId, metadataId string) (string, error) {
	metadataDir, err := getMetadataDirectory(vaultId, issuerId)
	if err != nil {
		return "", err
	}

	return filepath.Join(metadataDir, metadataId), nil
}

func GetMetadataFilePath(vaultId, issuerId, metadataId string) (string, error) {
	metadataIdDir, err := GetMetadataIdDirectory(vaultId, issuerId, metadataId)
	if err != nil {
		return "", err
	}

	return filepath.Join(metadataIdDir, "metadata.json"), nil
}

func (r *metadataFilesystemRepository) AddMetadata(
	vaultId, issuerId string, idpConfig *internalIssuerTypes.IdpConfig, resolverMetadata *coreV1alpha.ResolverMetadata,
) (*coreV1alpha.ResolverMetadata, error) {

	metadataDir, err := getMetadataDirectory(vaultId, issuerId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(metadataDir, internalIssuerConstants.DirPerm); err != nil {
		return nil, err
	}

	// Create metadata ID directory
	metadataIdDir, err := GetMetadataIdDirectory(vaultId, issuerId, *resolverMetadata.Id)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(metadataIdDir, internalIssuerConstants.DirPerm); err != nil {
		return nil, err
	}

	// Save metadata to file
	metadataFilePath, err := GetMetadataFilePath(vaultId, issuerId, *resolverMetadata.Id)
	if err != nil {
		return nil, err
	}

	metadataData, err := json.Marshal(resolverMetadata)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(metadataFilePath, metadataData, internalIssuerConstants.FilePerm)
	if err != nil {
		return nil, err
	}

	// Save the metadata config
	metadataConfig := internalIssuerTypes.MetadataConfig{
		IdpConfig: idpConfig,
	}

	metadataConfigFilePath := filepath.Join(metadataIdDir, "idp_config.json")
	metadataConfigData, err := json.Marshal(metadataConfig)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(metadataConfigFilePath, metadataConfigData, internalIssuerConstants.FilePerm)
	if err != nil {
		return nil, err
	}

	return resolverMetadata, nil
}

func (r *metadataFilesystemRepository) GetAllMetadata(vaultId, issuerId string) ([]*coreV1alpha.ResolverMetadata, error) {
	// Get the metadata directory
	metadataDir, err := getMetadataDirectory(vaultId, issuerId)
	if err != nil {
		return nil, err
	}

	// Read the metadata directory
	files, err := os.ReadDir(metadataDir)
	if err != nil {
		return nil, err
	}

	// List the metadata IDs
	var allMetadata []*coreV1alpha.ResolverMetadata

	for _, file := range files {
		if file.IsDir() {
			metadata, err := r.GetMetadata(vaultId, issuerId, file.Name())
			if err != nil {
				return nil, err
			}

			// Append the metadata to the list
			allMetadata = append(allMetadata, metadata)
		}
	}

	return allMetadata, nil
}

func (r *metadataFilesystemRepository) GetMetadata(
	vaultId, issuerId, metadataId string,
) (*coreV1alpha.ResolverMetadata, error) {
	// Get the metadata file path
	metadataFilePath, err := GetMetadataFilePath(vaultId, issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	// Read the metadata file
	metadataData, err := os.ReadFile(metadataFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the metadata data
	var metadata coreV1alpha.ResolverMetadata
	if err := json.Unmarshal(metadataData, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (r *metadataFilesystemRepository) RemoveMetadata(vaultId, issuerId, metadataId string) error {
	// Get the metadata directory
	metadataIdDir, err := GetMetadataIdDirectory(vaultId, issuerId, metadataId)
	if err != nil {
		return err
	}

	// Check if the metadata directory exists
	if _, err := os.Stat(metadataIdDir); os.IsNotExist(err) {
		return errors.New("metadata does not exist")
	}

	// Remove the metadata directory
	if err := os.RemoveAll(metadataIdDir); err != nil {
		return err
	}

	return nil
}
