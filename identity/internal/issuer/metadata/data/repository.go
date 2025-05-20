// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type MetadataRepository interface {
	GenerateMetadata(
		vaultId, issuerId string, idpConfig *internalIssuerTypes.IdpConfig,
	) (*coreV1alpha.ResolverMetadata, error)
	ListMetadataIds(vaultId, issuerId string) ([]string, error)
	GetMetadata(vaultId, issuerId, metadataId string) (*coreV1alpha.ResolverMetadata, error)
	ForgetMetadata(vaultId, issuerId, metadataId string) error
}
