// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"

	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/issuer/issuer/data"
	"github.com/agntcy/identity/internal/issuer/issuer/types"
	"github.com/agntcy/identity/internal/pkg/converters"
	"github.com/agntcy/identity/internal/pkg/nodeapi"
	"github.com/agntcy/identity/internal/pkg/oidc"
)

type IssuerService interface {
	RegisterIssuer(ctx context.Context, vaultId, keyId string, issuer *types.Issuer) (string, error)
	GetAllIssuers(vaultId, keyId string) ([]*types.Issuer, error)
	GetIssuer(vaultId, keyId, issuerId string) (*types.Issuer, error)
	ForgetIssuer(vaultId, keyId, issuerId string) error
}

type issuerService struct {
	issuerRepository data.IssuerRepository
	auth             oidc.Authenticator
	nodeClientPrv    nodeapi.ClientProvider
}

func NewIssuerService(
	issuerRepository data.IssuerRepository,
	auth oidc.Authenticator,
	nodeClientPrv nodeapi.ClientProvider,
) IssuerService {
	return &issuerService{
		issuerRepository: issuerRepository,
		auth:             auth,
		nodeClientPrv:    nodeClientPrv,
	}
}

func (s *issuerService) RegisterIssuer(
	ctx context.Context,
	vaultId, keyId string,
	issuer *types.Issuer,
) (string, error) {
	token, err := s.auth.Token(
		ctx,
		issuer.IdpConfig.IssuerUrl,
		issuer.IdpConfig.ClientId,
		issuer.IdpConfig.ClientSecret,
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

	err = client.RegisterIssuer(ctx, converters.Convert[issuertypes.Issuer](issuer.Issuer), &proof)
	if err != nil {
		return "", err
	}

	issuerId, err := s.issuerRepository.AddIssuer(vaultId, keyId, issuer)
	if err != nil {
		return "", err
	}

	return issuerId, nil
}

func (s *issuerService) GetAllIssuers(vaultId, keyId string) ([]*types.Issuer, error) {
	issuers, err := s.issuerRepository.GetAllIssuers(vaultId, keyId)
	if err != nil {
		return nil, err
	}

	return issuers, nil
}

func (s *issuerService) GetIssuer(vaultId, keyId, issuerId string) (*types.Issuer, error) {
	issuer, err := s.issuerRepository.GetIssuer(vaultId, keyId, issuerId)
	if err != nil {
		return nil, err
	}

	return issuer, nil
}

func (s *issuerService) ForgetIssuer(vaultId, keyId, issuerId string) error {
	err := s.issuerRepository.RemoveIssuer(vaultId, keyId, issuerId)
	if err != nil {
		return err
	}

	return nil
}
