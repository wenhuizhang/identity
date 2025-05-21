// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package filesystem

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	"github.com/agntcy/identity/internal/issuer/badge/data"
	internalIssuerConstants "github.com/agntcy/identity/internal/issuer/constants"
	metadataFilesystemRepository "github.com/agntcy/identity/internal/issuer/metadata/data/filesystem"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type badgeFilesystemRepository struct{}

func NewBadgeFilesystemRepository() data.BadgeRepository {
	return &badgeFilesystemRepository{}
}

// getBadgeDirectory returns the path to the badges directory for a metadata
func getBadgesDirectory(vaultId, issuerId, metadataId string) (string, error) {
	metadataIdDir, err := metadataFilesystemRepository.GetMetadataIdDirectory(vaultId, issuerId, metadataId)
	if err != nil {
		return "", err
	}

	return filepath.Join(metadataIdDir, "badges"), nil
}

// GetBadgeIdDirectory returns the path to the badge ID directory
func GetBadgeIdDirectory(vaultId, issuerId, metadataId, badgeId string) (string, error) {
	badgesDir, err := getBadgesDirectory(vaultId, issuerId, metadataId)
	if err != nil {
		return "", err
	}

	return filepath.Join(badgesDir, badgeId), nil
}

// GetBadgeFilePath returns the path to the badge file
func GetBadgeFilePath(vaultId, issuerId, metadataId, badgeId string) (string, error) {
	badgeIdDir, err := GetBadgeIdDirectory(vaultId, issuerId, metadataId, badgeId)
	if err != nil {
		return "", err
	}

	return filepath.Join(badgeIdDir, "badge.json"), nil
}

func (r *badgeFilesystemRepository) AddBadge(
	vaultId, issuerId, metadataId string, envelopedCredential *coreV1alpha.EnvelopedCredential,
) (string, error) {

	// Ensure badges directory exists
	badgesDir, err := getBadgesDirectory(vaultId, issuerId, metadataId)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(badgesDir, internalIssuerConstants.DirPerm); err != nil {
		return "", err
	}

	// Create badge ID directory with a unique ID
	badgeId := uuid.New().String()

	badgesIdDir, err := GetBadgeIdDirectory(vaultId, issuerId, metadataId, badgeId)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(badgesIdDir, internalIssuerConstants.DirPerm); err != nil {
		return "", err
	}

	// Save badge to file
	badgeFilePath, err := GetBadgeFilePath(vaultId, issuerId, metadataId, badgeId)
	if err != nil {
		return "", err
	}

	badgeData, err := json.Marshal(&envelopedCredential)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(badgeFilePath, badgeData, internalIssuerConstants.FilePerm); err != nil {
		return "", err
	}

	return badgeId, nil
}

func (r *badgeFilesystemRepository) GetAllBadges(vaultId, issuerId, metadataId string) ([]*internalIssuerTypes.Badge, error) {
	// Get the badges directory
	badgesDir, err := getBadgesDirectory(vaultId, issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	// Read the badges directory
	files, err := os.ReadDir(badgesDir)
	if err != nil {
		return nil, err
	}

	// List the badge IDs
	var badges []*internalIssuerTypes.Badge

	for _, file := range files {
		if file.IsDir() {
			badge, err := r.GetBadge(vaultId, issuerId, metadataId, file.Name())
			if err != nil {
				return nil, err
			}
			// Append the badge to the list
			badges = append(badges, &internalIssuerTypes.Badge{
				Id:    file.Name(),
				Badge: badge,
			})
		}
	}

	return badges, nil
}

func (r *badgeFilesystemRepository) GetBadge(
	vaultId, issuerId, metadataId, badgeId string,
) (*coreV1alpha.EnvelopedCredential, error) {
	// Get the badge file path
	badgeFilePath, err := GetBadgeFilePath(vaultId, issuerId, metadataId, badgeId)
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

func (r *badgeFilesystemRepository) RemoveBadge(vaultId, issuerId, metadataId, badgeId string) error {
	// Get the badge directory
	badgeIdDir, err := GetBadgeIdDirectory(vaultId, issuerId, metadataId, badgeId)
	if err != nil {
		return err
	}

	// Check if the badge directory exists
	if _, err := os.Stat(badgeIdDir); os.IsNotExist(err) {
		return errors.New("metadata does not exist")
	}

	// Remove the badge directory
	if err := os.RemoveAll(badgeIdDir); err != nil {
		return err
	}

	return nil
}
