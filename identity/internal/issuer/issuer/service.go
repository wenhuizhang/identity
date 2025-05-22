// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"
	"log"

	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/issuer/issuer/data"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/pkg/converters"
	"github.com/agntcy/identity/internal/pkg/oidc"
)

type IssuerService interface {
	RegisterIssuer(ctx context.Context, vaultId string, issuer *internalIssuerTypes.Issuer) (string, error)
	GetAllIssuers(vaultId string) ([]*internalIssuerTypes.Issuer, error)
	GetIssuer(vaultId, issuerId string) (*internalIssuerTypes.Issuer, error)
	ForgetIssuer(vaultId, issuerId string) error
}

type issuerService struct {
	issuerRepository data.IssuerRepository
	auth             oidc.Authenticator
	nodeClientPrv    data.NodeClientProvider
}

func NewIssuerService(
	issuerRepository data.IssuerRepository,
	auth oidc.Authenticator,
	nodeClientPrv data.NodeClientProvider,
) IssuerService {
	return &issuerService{
		issuerRepository: issuerRepository,
		auth:             auth,
		nodeClientPrv:    nodeClientPrv,
	}
}

func (s *issuerService) RegisterIssuer(
	ctx context.Context,
	vaultId string,
	issuer *internalIssuerTypes.Issuer,
) (string, error) {
	// Check connection to identity node
	// Check connection to idp
	// Check if idp is already created locally
	// Check if idp is already registered on the identity node
	// Register idp on the identity node

	token, err := s.auth.Token(
		ctx,
		issuer.IdpConfig.IssuerUrl,
		issuer.IdpConfig.ClientId,
		issuer.IdpConfig.ClientSecret,
	)
	if err != nil {
		return "", err // TODO: return ErrInfo
	}

	proof := vctypes.Proof{
		Type:       "JWT",
		ProofValue: token,
	}

	log.Default().Printf("Registering issuer with request: %s\n", *issuer.Issuer.CommonName)

	client, err := s.nodeClientPrv.New(issuer.IdentityNodeConfig.IdentityNodeAddress)
	if err != nil {
		return "", err
	}

	err = client.RegisterIssuer(ctx, converters.Convert[issuertypes.Issuer](issuer.Issuer), &proof)
	if err != nil {
		return "", err
	}

	issuerId, err := s.issuerRepository.AddIssuer(vaultId, issuer)
	if err != nil {
		return "", err
	}

	return issuerId, nil
}

func (s *issuerService) GetAllIssuers(vaultId string) ([]*internalIssuerTypes.Issuer, error) {
	issuers, err := s.issuerRepository.GetAllIssuers(vaultId)
	if err != nil {
		return nil, err
	}

	return issuers, nil
}

func (s *issuerService) GetIssuer(vaultId, issuerId string) (*internalIssuerTypes.Issuer, error) {
	issuer, err := s.issuerRepository.GetIssuer(vaultId, issuerId)
	if err != nil {
		return nil, err
	}

	return issuer, nil
}

func (s *issuerService) ForgetIssuer(vaultId, issuerId string) error {
	err := s.issuerRepository.RemoveIssuer(vaultId, issuerId)
	if err != nil {
		return err
	}

	return nil
}
