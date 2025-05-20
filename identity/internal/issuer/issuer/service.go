// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	"github.com/agntcy/identity/internal/issuer/issuer/data"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type IssuerService interface {
	RegisterIssuer(vaultId, identityNodeAddress string, idpConfig internalIssuerTypes.IdpConfig) (string, error)
	ListIssuerIds(vaultId string) ([]string, error)
	GetIssuer(vaultId, issuerId string) (*coreV1alpha.Issuer, error)
	ForgetIssuer(vaultId, issuerId string) error
}

type issuerService struct {
	issuerRepository data.IssuerRepository
}

func NewIssuerService(
	issuerRepository data.IssuerRepository,
) IssuerService {
	return &issuerService{
		issuerRepository: issuerRepository,
	}
}

func (s *issuerService) RegisterIssuer(
	vaultId, identityNodeAddress string, idpConfig internalIssuerTypes.IdpConfig,
) (string, error) {
	issuerId, err := s.issuerRepository.RegisterIssuer(vaultId, identityNodeAddress, idpConfig)
	if err != nil {
		return "", err
	}

	return issuerId, nil
}

func (s *issuerService) ListIssuerIds(vaultId string) ([]string, error) {
	issuerIds, err := s.issuerRepository.ListIssuerIds(vaultId)
	if err != nil {
		return nil, err
	}

	return issuerIds, nil
}

func (s *issuerService) GetIssuer(vaultId, issuerId string) (*coreV1alpha.Issuer, error) {
	issuer, err := s.issuerRepository.GetIssuer(vaultId, issuerId)
	if err != nil {
		return nil, err
	}

	return issuer, nil
}

func (s *issuerService) ForgetIssuer(vaultId, issuerId string) error {
	err := s.issuerRepository.ForgetIssuer(vaultId, issuerId)
	if err != nil {
		return err
	}

	return nil
}

func TestIdpConnection(clientId, clientSecret, issuerUrl string) (*oauth2.Token, error) {
	// Test the connection to the Identity Provider
	ctx := context.Background()

	// Discover OIDC provider config
	provider, err := oidc.NewProvider(ctx, issuerUrl)
	if err != nil {
		return nil, err
	}

	// Set up the OAuth2 client credentials config
	conf := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     provider.Endpoint().TokenURL,
		Scopes:       []string{},
	}

	// Retrieve a token
	token, err := conf.Token(ctx)
	if err != nil {
		return nil, err
	}

	return token, nil
}
