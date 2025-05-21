// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type MetadataRepository interface {
	AddMetadata(
		vaultId, issuerId string, metadata *internalIssuerTypes.Metadata,
	) (string, error)
	GetAllMetadata(vaultId, issuerId string) ([]*internalIssuerTypes.Metadata, error)
	GetMetadata(vaultId, issuerId, metadataId string) (*internalIssuerTypes.Metadata, error)
	RemoveMetadata(vaultId, issuerId, metadataId string) error
}
