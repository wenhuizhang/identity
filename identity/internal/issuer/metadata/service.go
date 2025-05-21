// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"log"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	nodeV1alpha "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
	issuerData "github.com/agntcy/identity/internal/issuer/issuer/data"
	"github.com/agntcy/identity/internal/issuer/metadata/data"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/google/uuid"
)

type MetadataService interface {
	GenerateMetadata(
		vaultId, issuerId string, idpConfig *internalIssuerTypes.IdpConfig,
	) (*coreV1alpha.ResolverMetadata, error)
	GetAllMetadata(vaultId, issuerId string) ([]*coreV1alpha.ResolverMetadata, error)
	GetMetadata(vaultId, issuerId, metadataId string) (*coreV1alpha.ResolverMetadata, error)
	ForgetMetadata(vaultId, issuerId, metadataId string) error
}

type metadataService struct {
	metadataRepository data.MetadataRepository
	issuerRepository   issuerData.IssuerRepository
}

func NewMetadataService(
	metadataRepository data.MetadataRepository,
	issuerRepository issuerData.IssuerRepository,
) MetadataService {
	return &metadataService{
		metadataRepository: metadataRepository,
		issuerRepository:   issuerRepository,
	}
}

func (s *metadataService) GenerateMetadata(
	vaultId, issuerId string, idpConfig *internalIssuerTypes.IdpConfig,
) (*coreV1alpha.ResolverMetadata, error) {

	// load the issuer
	issuer, err := s.issuerRepository.GetIssuer(vaultId, issuerId)
	if err != nil {
		return nil, err
	}

	proof := coreV1alpha.Proof{
		Type:         func() *string { s := "RsaSignature2018"; return &s }(),
		ProofPurpose: func() *string { s := "assertionMethod"; return &s }(),
		ProofValue:   func() *string { s := "example-proof-value"; return &s }(),
	}

	generateMetadataRequest := nodeV1alpha.GenerateRequest{
		Issuer: issuer,
		Proof:  &proof,
	}

	// Call the client to generate metadata
	log.Default().Println("Generating metadata with request: ", &generateMetadataRequest)

	resolverMetadata := coreV1alpha.ResolverMetadata{
		Id:                 func() *string { s := uuid.New().String(); return &s }(),
		VerificationMethod: nil,
		Service:            nil,
		AssertionMethod:    nil,
	}

	metadata, err := s.metadataRepository.AddMetadata(vaultId, issuerId, idpConfig, &resolverMetadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *metadataService) GetAllMetadata(vaultId, issuerId string) ([]*coreV1alpha.ResolverMetadata, error) {
	metadata, err := s.metadataRepository.GetAllMetadata(vaultId, issuerId)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *metadataService) GetMetadata(vaultId, issuerId, metadataId string) (*coreV1alpha.ResolverMetadata, error) {
	metadata, err := s.metadataRepository.GetMetadata(vaultId, issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *metadataService) ForgetMetadata(vaultId, issuerId, metadataId string) error {
	err := s.metadataRepository.RemoveMetadata(vaultId, issuerId, metadataId)
	if err != nil {
		return err
	}

	return nil
}
