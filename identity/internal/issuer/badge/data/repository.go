// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type BadgeRepository interface {
	AddBadge(vaultId, issuerId, metadataId string, envelopedCredential *internalIssuerTypes.Badge) (string, error)
	GetAllBadges(vaultId, issuerId, metadataId string) ([]*internalIssuerTypes.Badge, error)
	GetBadge(vaultId, issuerId, metadataId, badgeId string) (*internalIssuerTypes.Badge, error)
	RemoveBadge(vaultId, issuerId, metadataId, badgeId string) error
}
