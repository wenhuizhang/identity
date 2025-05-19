// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/agntcy/identity/internal/issuer/issuer"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	nodeV1alpha "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
	internalIssuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

func getMetadataDirectory(issuerId string) (string, error) {
	issuerIdDir, err := issuer.GetIssuerIdDirectory(issuerId)
	if err != nil {
		return "", err
	}

	return filepath.Join(issuerIdDir, "metadata"), nil
}

func GetMetadataIdDirectory(issuerId, metadataId string) (string, error) {
	metadataDir, err := getMetadataDirectory(issuerId)
	if err != nil {
		return "", err
	}

	return filepath.Join(metadataDir, metadataId), nil
}

func GetMetadataFilePath(issuerId, metadataId string) (string, error) {
	metadataIdDir, err := GetMetadataIdDirectory(issuerId, metadataId)
	if err != nil {
		return "", err
	}

	return filepath.Join(metadataIdDir, "metadata.json"), nil
}

// SaveMetadata creates the necessary directories and saves metadata to file
func saveMetadata(issuerId string, resolverMetadata *coreV1alpha.ResolverMetadata) error {
	// Ensure metadata directory exists
	metadataDir, err := getMetadataDirectory(issuerId)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(metadataDir, internalIssuerConstants.DirPerm); err != nil {
		return err
	}

	// Create metadata ID directory
	metadataIdDir, err := GetMetadataIdDirectory(issuerId, *resolverMetadata.Id)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(metadataIdDir, internalIssuerConstants.DirPerm); err != nil {
		return err
	}

	// Save metadata to file
	metadataFilePath, err := GetMetadataFilePath(issuerId, *resolverMetadata.Id)
	if err != nil {
		return err
	}

	metadataData, err := json.Marshal(resolverMetadata)
	if err != nil {
		return err
	}

	return os.WriteFile(metadataFilePath, metadataData, internalIssuerConstants.FilePerm)
}

func GenerateMetadata(issuerId string, idpConfig *internalIssuerTypes.IdpConfig) (*coreV1alpha.ResolverMetadata, error) {
	// load the issuer from the local storage
	issuerFilePath, err := issuer.GetIssuerFilePath(issuerId)
	if err != nil {
		return nil, err
	}

	issuerData, err := os.ReadFile(issuerFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the issuer data
	var issuer coreV1alpha.Issuer
	if err := json.Unmarshal(issuerData, &issuer); err != nil {
		return nil, err
	}

	proof := coreV1alpha.Proof{
		Type:         func() *string { s := "RsaSignature2018"; return &s }(),
		ProofPurpose: func() *string { s := "assertionMethod"; return &s }(),
		ProofValue:   func() *string { s := "example-proof-value"; return &s }(),
	}

	generateMetadataRequest := nodeV1alpha.GenerateRequest{
		Issuer: &issuer,
		Proof:  &proof,
	}

	// Call the client to generate metadata
	log.Default().Println("Generating metadata with request: ", &generateMetadataRequest)

	resolverMetadata := coreV1alpha.ResolverMetadata{
		Id:                 func() *string { s := uuid.New().String(); return &s }(),
		VerificationMethod: nil,
		Service:            nil,
		AssertionMethod:    nil,
	}

	// Save the metadata to disk
	if err := saveMetadata(issuerId, &resolverMetadata); err != nil {
		return nil, err
	}

	return &resolverMetadata, nil
}

func ListMetadataIds(issuerId string) ([]string, error) {
	// Get the metadata directory
	metadataDir, err := getMetadataDirectory(issuerId)
	if err != nil {
		return nil, err
	}

	// Create directory if it doesn't exist
	if _, err := os.Stat(metadataDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read the metadata directory
	files, err := os.ReadDir(metadataDir)
	if err != nil {
		return nil, err
	}

	// List the metadata IDs
	var metadataIds []string

	for _, file := range files {
		if file.IsDir() {
			metadataIds = append(metadataIds, file.Name())
		}
	}

	return metadataIds, nil
}

func GetMetadata(issuerId, metadataId string) (*coreV1alpha.ResolverMetadata, error) {
	// Get the metadata file path
	metadataFilePath, err := GetMetadataFilePath(issuerId, metadataId)
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

func ForgetMetadata(issuerId, metadataId string) error {
	// Get the metadata directory
	metadataIdDir, err := GetMetadataIdDirectory(issuerId, metadataId)
	if err != nil {
		return err
	}

	// Check if the metadata directory exists
	if _, err := os.Stat(metadataIdDir); os.IsNotExist(err) {
		return errors.New("Metadata does not exist")
	}

	// Remove the metadata directory
	if err := os.RemoveAll(metadataIdDir); err != nil {
		return err
	}

	return nil
}
