// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	issuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	"github.com/agntcy/identity/internal/issuer/metadata"
	"github.com/google/uuid"
)

type publishRequest struct {
	VC    *vctypes.EnvelopedCredential `json:"vc"`
	Proof *vctypes.Proof               `json:"proof"`
}

// getBadgeDirectory returns the path to the badges directory for a metadata
func getBadgesDirectory(issuerId, metadataId string) (string, error) {
	metadataIdDir, err := metadata.GetMetadataIdDirectory(issuerId, metadataId)
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

func IssueBadge(issuerId, metadataId, badgeValueFilePath string) (*vctypes.EnvelopedCredential, error) {
	// Read the badge value from the file
	badgeValueData, err := os.ReadFile(badgeValueFilePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the badge value
	var badgeValue string
	if err := json.Unmarshal(badgeValueData, &badgeValue); err != nil {
		return nil, err
	}

	envelopedCredential := vctypes.EnvelopedCredential{
		EnvelopeType: vctypes.CREDENTIAL_ENVELOPE_TYPE_JOSE,
		Value:        badgeValue,
	}

	// Ensure badges directory exists
	badgesDir, err := getBadgesDirectory(issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(badgesDir, issuerConstants.DirPerm); err != nil {
		return nil, err
	}

	// Create badge ID directory with a unique ID
	badgeId := uuid.New().String()

	badgesIdDir, err := GetBadgeIdDirectory(issuerId, metadataId, badgeId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(badgesIdDir, issuerConstants.DirPerm); err != nil {
		return nil, err
	}

	// Save badge to file
	badgeFilePath, err := GetBadgeFilePath(issuerId, metadataId, badgeId)
	if err != nil {
		return nil, err
	}

	badgeData, err := json.Marshal(envelopedCredential)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(badgeFilePath, badgeData, issuerConstants.FilePerm); err != nil {
		return nil, err
	}

	return &envelopedCredential, nil
}

func PublishBadge(
	issuerId,
	metadataId string,
	badge *vctypes.EnvelopedCredential,
) (*vctypes.EnvelopedCredential, error) {
	proof := vctypes.Proof{
		Type:         "RsaSignature2018",
		ProofPurpose: "assertionMethod",
		ProofValue:   "example-proof-value",
	}

	publishRequest := publishRequest{
		VC:    badge,
		Proof: &proof,
	}

	log.Default().Println("Publishing badge with request: ", publishRequest)

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

func GetBadge(issuerId, metadataId, badgeId string) (*vctypes.EnvelopedCredential, error) {
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
	var badge vctypes.EnvelopedCredential
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

	// Remove the badge directory
	if err := os.RemoveAll(badgeIdDir); err != nil {
		return err
	}

	return nil
}
