// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	idTypes "github.com/agntcy/identity/internal/core/id/types"
	"github.com/agntcy/identity/internal/core/issuer/types"
	issuerTypes "github.com/agntcy/identity/internal/core/issuer/types"
	vcTypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/issuer/constants"
	"github.com/agntcy/identity/internal/issuer/issuer"
	issuerTypesInternal "github.com/agntcy/identity/internal/issuer/types"
	"github.com/google/uuid"
)

type generateMetadataRequest struct {
	Issuer issuerTypes.Issuer `json:"issuer"`
	Proof  vcTypes.Proof      `json:"proof"`
}

func getMetadataDirectory(issuerId string) (string, error) {
	issuerIdDir, err := issuer.GetIssuerIdDirectory(issuerId)
	if err != nil {
		return "", err
	}
	return filepath.Join(issuerIdDir, "metadata"), nil
}

func GetMetadataIdDirectory(issuerId string, metadataId string) (string, error) {
	metadataDir, err := getMetadataDirectory(issuerId)
	if err != nil {
		return "", err
	}

	return filepath.Join(metadataDir, metadataId), nil
}

func GetMetadataFilePath(issuerId string, metadataId string) (string, error) {
	metadataIdDir, err := GetMetadataIdDirectory(issuerId, metadataId)
	if err != nil {
		return "", err
	}

	return filepath.Join(metadataIdDir, "metadata.json"), nil
}

func GenerateMetadata(issuerId string, idpConfig *issuerTypesInternal.IdpConfig) (*idTypes.ResolverMetadata, error) {
	// load the issuer from the local storage
	issuerFilePath, err := issuer.GetIssuerFilePath(issuerId)
	if err != nil {
		log.Default().Println("Error getting issuer file path: ", err)
		return nil, err
	}

	// Read the issuer file
	issuerData, err := os.ReadFile(issuerFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the issuer data
	var issuer types.Issuer
	if err := json.Unmarshal(issuerData, &issuer); err != nil {
		return nil, err
	}

	proof := vcTypes.Proof{
		Type:         "RsaSignature2018",
		ProofPurpose: "assertionMethod",
		ProofValue:   "example-proof-value",
	}

	generateMetadataRequest := generateMetadataRequest{
		Issuer: issuer,
		Proof:  proof,
	}

	// Call the client to generate metadata
	log.Default().Println("Generating metadata with request: ", generateMetadataRequest)

	resolverMetadata := &idTypes.ResolverMetadata{
		ID:                 uuid.New().String(),
		VerificationMethod: nil,
		Service:            nil,
		AssertionMethod:    nil,
	}

	// Ensure metadata directory exists
	metadataDir, err := getMetadataDirectory(issuerId)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(metadataDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Create metadata ID directory
	metadataIdDir, err := GetMetadataIdDirectory(issuerId, resolverMetadata.ID)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(metadataIdDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Save metadata to file
	metadataFilePath, err := GetMetadataFilePath(issuerId, resolverMetadata.ID)
	if err != nil {
		return nil, err
	}

	metadataData, err := json.Marshal(resolverMetadata)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(metadataFilePath, metadataData, constants.FilePerm); err != nil {
		return nil, err
	}

	return resolverMetadata, nil
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

func GetMetadata(issuerId string, metadataId string) (*idTypes.ResolverMetadata, error) {
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
	var metadata idTypes.ResolverMetadata
	if err := json.Unmarshal(metadataData, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func ForgetMetadata(issuerId string, metadataId string) error {
	// Get the metadata directory
	metadataIdDir, err := GetMetadataIdDirectory(issuerId, metadataId)
	if err != nil {
		return err
	}

	// Remove the metadata directory
	if err := os.RemoveAll(metadataIdDir); err != nil {
		return err
	}

	return nil
}
