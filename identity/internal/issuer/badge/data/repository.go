// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type BadgeRepository interface {
	AddBadge(vaultId, issuerId, metadataId string, envelopedCredential *coreV1alpha.EnvelopedCredential) (string, error)
	GetAllBadges(vaultId, issuerId, metadataId string) ([]*internalIssuerTypes.Badge, error)
	GetBadge(vaultId, issuerId, metadataId, badgeId string) (*coreV1alpha.EnvelopedCredential, error)
	RemoveBadge(vaultId, issuerId, metadataId, badgeId string) error
}
