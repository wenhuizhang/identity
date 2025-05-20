// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
)

type BadgeRepository interface {
	IssueBadge(vaultId, issuerId, metadataId, badgeValueFilePath string) (string, error)
	PublishBadge(
		vaultId, issuerId, metadataId string, badge *coreV1alpha.EnvelopedCredential,
	) (*coreV1alpha.EnvelopedCredential, error)
	ListBadgeIds(vaultId, issuerId, metadataId string) ([]string, error)
	GetBadge(vaultId, issuerId, metadataId, badgeId string) (*coreV1alpha.EnvelopedCredential, error)
	ForgetBadge(vaultId, issuerId, metadataId, badgeId string) error
}
