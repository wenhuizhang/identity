// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	"github.com/agntcy/identity/internal/issuer/metadata/data"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type MetadataService interface {
	GenerateMetadata(
		vaultId, issuerId string, idpConfig *internalIssuerTypes.IdpConfig,
	) (*coreV1alpha.ResolverMetadata, error)
	ListMetadataIds(vaultId, issuerId string) ([]string, error)
	GetMetadata(vaultId, issuerId, metadataId string) (*coreV1alpha.ResolverMetadata, error)
	ForgetMetadata(vaultId, issuerId, metadataId string) error
}

type metadataService struct {
	metadataRepository data.MetadataRepository
}

func NewMetadataService(
	metadataRepository data.MetadataRepository,
) MetadataService {
	return &metadataService{
		metadataRepository: metadataRepository,
	}
}

func (s *metadataService) GenerateMetadata(
	vaultId, issuerId string, idpConfig *internalIssuerTypes.IdpConfig,
) (*coreV1alpha.ResolverMetadata, error) {
	metadata, err := s.metadataRepository.GenerateMetadata(vaultId, issuerId, idpConfig)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *metadataService) ListMetadataIds(vaultId, issuerId string) ([]string, error) {
	metadataIds, err := s.metadataRepository.ListMetadataIds(vaultId, issuerId)
	if err != nil {
		return nil, err
	}

	return metadataIds, nil
}

func (s *metadataService) GetMetadata(vaultId, issuerId, metadataId string) (*coreV1alpha.ResolverMetadata, error) {
	metadata, err := s.metadataRepository.GetMetadata(vaultId, issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *metadataService) ForgetMetadata(vaultId, issuerId, metadataId string) error {
	err := s.metadataRepository.ForgetMetadata(vaultId, issuerId, metadataId)
	if err != nil {
		return err
	}

	return nil
}
