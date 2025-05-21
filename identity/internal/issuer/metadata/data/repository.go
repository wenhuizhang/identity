// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type MetadataRepository interface {
	AddMetadata(
		vaultId, issuerId string, idpConfig *internalIssuerTypes.IdpConfig, resolverMetadata *coreV1alpha.ResolverMetadata,
	) (*coreV1alpha.ResolverMetadata, error)
	GetAllMetadata(vaultId, issuerId string) ([]*coreV1alpha.ResolverMetadata, error)
	GetMetadata(vaultId, issuerId, metadataId string) (*coreV1alpha.ResolverMetadata, error)
	RemoveMetadata(vaultId, issuerId, metadataId string) error
}
