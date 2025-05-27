// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"

	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	issuerData "github.com/agntcy/identity/internal/issuer/issuer/data"
	"github.com/agntcy/identity/internal/issuer/metadata/data"
	"github.com/agntcy/identity/internal/issuer/metadata/types"
	idptypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/pkg/nodeapi"
	"github.com/agntcy/identity/internal/pkg/oidc"
)

type MetadataService interface {
	GenerateMetadata(
		ctx context.Context,
		vaultId, keyId, issuerId string,
		idpConfig *idptypes.IdpConfig,
	) (string, error)
	GetAllMetadata(vaultId, keyId, issuerId string) ([]*types.Metadata, error)
	GetMetadata(vaultId, keyId, issuerId, metadataId string) (*types.Metadata, error)
	ForgetMetadata(vaultId, keyId, issuerId, metadataId string) error
}

type metadataService struct {
	metadataRepository data.MetadataRepository
	issuerRepository   issuerData.IssuerRepository
	auth               oidc.Authenticator
	nodeClientPrv      nodeapi.ClientProvider
}

func NewMetadataService(
	metadataRepository data.MetadataRepository,
	issuerRepository issuerData.IssuerRepository,
	auth oidc.Authenticator,
	nodeClientPrv nodeapi.ClientProvider,
) MetadataService {
	return &metadataService{
		metadataRepository: metadataRepository,
		issuerRepository:   issuerRepository,
		auth:               auth,
		nodeClientPrv:      nodeClientPrv,
	}
}

func (s *metadataService) GenerateMetadata(
	ctx context.Context,
	vaultId, keyId, issuerId string,
	idpConfig *idptypes.IdpConfig,
) (string, error) {
	issuer, err := s.issuerRepository.GetIssuer(vaultId, keyId, issuerId)
	if err != nil {
		return "", err
	}

	token, err := s.auth.Token(
		ctx,
		idpConfig.IssuerUrl,
		idpConfig.ClientId,
		idpConfig.ClientSecret,
	)
	if err != nil {
		return "", err
	}

	proof := vctypes.Proof{
		Type:       "JWT",
		ProofValue: token,
	}

	client, err := s.nodeClientPrv.New(issuer.IdentityNodeURL)
	if err != nil {
		return "", err
	}

	md, err := client.GenerateID(ctx, &issuer.Issuer, &proof)
	if err != nil {
		return "", err
	}

	metadata := types.Metadata{
		ResolverMetadata: *md,
		IdpConfig:        idpConfig,
	}

	metadataId, err := s.metadataRepository.AddMetadata(vaultId, keyId, issuerId, &metadata)
	if err != nil {
		return "", err
	}

	return metadataId, nil
}

func (s *metadataService) GetAllMetadata(vaultId, keyId, issuerId string) ([]*types.Metadata, error) {
	metadata, err := s.metadataRepository.GetAllMetadata(vaultId, keyId, issuerId)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *metadataService) GetMetadata(vaultId, keyId, issuerId, metadataId string) (*types.Metadata, error) {
	metadata, err := s.metadataRepository.GetMetadata(vaultId, keyId, issuerId, metadataId)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (s *metadataService) ForgetMetadata(vaultId, keyId, issuerId, metadataId string) error {
	err := s.metadataRepository.RemoveMetadata(vaultId, keyId, issuerId, metadataId)
	if err != nil {
		return err
	}

	return nil
}
