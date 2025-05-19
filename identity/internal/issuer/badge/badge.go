// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	internalIssuerMetadata "github.com/agntcy/identity/internal/issuer/metadata"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	nodeV1alpha "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
	internalIssuerConstants "github.com/agntcy/identity/internal/issuer/constants"
)

// getBadgeDirectory returns the path to the badges directory for a metadata
func getBadgesDirectory(issuerId, metadataId string) (string, error) {
	metadataIdDir, err := internalIssuerMetadata.GetMetadataIdDirectory(issuerId, metadataId)
	if err != nil {
		return "", err
	}

	return filepath.Join(metadataIdDir, "badges"), nil
}

// GetBadgeIdDirectory returns the path to the badge ID directory
func GetBadgeIdDirectory(issuerId, metadataId, badgeId string) (string, error) {
	badgesDir, err := getBadgesDirectory(issuerId, metadataId)
	if err != nil {
		return "", err
	}

	return filepath.Join(badgesDir, badgeId), nil
}

// GetBadgeFilePath returns the path to the badge file
func GetBadgeFilePath(issuerId, metadataId, badgeId string) (string, error) {
	badgeIdDir, err := GetBadgeIdDirectory(issuerId, metadataId, badgeId)
	if err != nil {
		return "", err
	}

	return filepath.Join(badgeIdDir, "badge.json"), nil
}

func IssueBadge(issuerId, metadataId, badgeValueFilePath string) (*coreV1alpha.EnvelopedCredential, error) {
	// Read the badge value from the file
	badgeValueData, err := os.ReadFile(badgeValueFilePath)
	if err != nil {
		return nil, err
	}

	// Convert the badge value to a string
	badgeValue := string(badgeValueData)

	envelopedCredential := coreV1alpha.EnvelopedCredential{
		EnvelopeType: coreV1alpha.CredentialEnvelopeType_CREDENTIAL_ENVELOPE_TYPE_JOSE.Enum(),
		Value:        &badgeValue,
	}

	// Ensure badges directory exists
	badgesDir, err := getBadgesDirectory(issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(badgesDir, internalIssuerConstants.DirPerm); err != nil {
		return nil, err
	}

	// Create badge ID directory with a unique ID
	badgeId := uuid.New().String()

	badgesIdDir, err := GetBadgeIdDirectory(issuerId, metadataId, badgeId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(badgesIdDir, internalIssuerConstants.DirPerm); err != nil {
		return nil, err
	}

	// Save badge to file
	badgeFilePath, err := GetBadgeFilePath(issuerId, metadataId, badgeId)
	if err != nil {
		return nil, err
	}

	badgeData, err := json.Marshal(&envelopedCredential)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(badgeFilePath, badgeData, internalIssuerConstants.FilePerm); err != nil {
		return nil, err
	}

	return &envelopedCredential, nil
}

func PublishBadge(
	issuerId,
	metadataId string,
	badge *coreV1alpha.EnvelopedCredential,
) (*coreV1alpha.EnvelopedCredential, error) {
	proof := coreV1alpha.Proof{
		Type:         func() *string { s := "RsaSignature2018"; return &s }(),
		ProofPurpose: func() *string { s := "assertionMethod"; return &s }(),
		ProofValue:   func() *string { s := "example-proof-value"; return &s }(),
	}

	publishRequest := nodeV1alpha.PublishRequest{
		Vc:    badge,
		Proof: &proof,
	}

	log.Default().Println("Publishing badge with request: ", &publishRequest)

	return badge, nil
}

func ListBadgeIds(issuerId, metadataId string) ([]string, error) {
	// Get the badges directory
	badgesDir, err := getBadgesDirectory(issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	// Create directory if it doesn't exist
	if _, err := os.Stat(badgesDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read the badges directory
	files, err := os.ReadDir(badgesDir)
	if err != nil {
		return nil, err
	}

	// List the badge IDs
	var badgeIds []string

	for _, file := range files {
		if file.IsDir() {
			badgeIds = append(badgeIds, file.Name())
		}
	}

	return badgeIds, nil
}

func GetBadge(issuerId, metadataId, badgeId string) (*coreV1alpha.EnvelopedCredential, error) {
	// Get the badge file path
	badgeFilePath, err := GetBadgeFilePath(issuerId, metadataId, badgeId)
	if err != nil {
		return nil, err
	}

	// Read the badge file
	badgeData, err := os.ReadFile(badgeFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the badge data
	var badge coreV1alpha.EnvelopedCredential
	if err := json.Unmarshal(badgeData, &badge); err != nil {
		return nil, err
	}

	return &badge, nil
}

func ForgetBadge(issuerId, metadataId, badgeId string) error {
	// Get the badge directory
	badgeIdDir, err := GetBadgeIdDirectory(issuerId, metadataId, badgeId)
	if err != nil {
		return err
	}

	// Check if the badge directory exists
	if _, err := os.Stat(badgeIdDir); os.IsNotExist(err) {
		return errors.New("Metadata does not exist")
	}

	// Remove the badge directory
	if err := os.RemoveAll(badgeIdDir); err != nil {
		return err
	}

	return nil
}
