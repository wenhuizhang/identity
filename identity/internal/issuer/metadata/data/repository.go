// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	"github.com/agntcy/identity/internal/issuer/metadata/types"
)

type MetadataRepository interface {
	AddMetadata(
		vaultId, issuerId string, metadata *types.Metadata,
	) (string, error)
	GetAllMetadata(vaultId, issuerId string) ([]*types.Metadata, error)
	GetMetadata(vaultId, issuerId, metadataId string) (*types.Metadata, error)
	RemoveMetadata(vaultId, issuerId, metadataId string) error
}
