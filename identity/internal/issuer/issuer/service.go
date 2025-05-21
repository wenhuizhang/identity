// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"
	"log"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	nodeV1alpha "github.com/agntcy/identity/api/agntcy/identity/node/v1alpha1"
	"github.com/agntcy/identity/internal/issuer/issuer/data"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type IssuerService interface {
	RegisterIssuer(vaultId, identityNodeAddress string, idpConfig internalIssuerTypes.IdpConfig) (string, error)
	GetAllIssuers(vaultId string) ([]*internalIssuerTypes.Issuer, error)
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

	// Check connection to identity node
	// Check connection to idp
	// Check if idp is already created locally
	// Check if idp is already registered on the identity node
	// Register idp on the identity node
	issuer := coreV1alpha.Issuer{
		Organization:    func() *string { s := "AGNTCY1"; return &s }(),
		SubOrganization: func() *string { s := "AGNTCY2"; return &s }(),
		CommonName:      func() *string { s := "AGNTCY3"; return &s }(),
	}
	proof := coreV1alpha.Proof{
		Type:         func() *string { s := "RsaSignature2018"; return &s }(),
		ProofPurpose: func() *string { s := "assertionMethod"; return &s }(),
		ProofValue:   func() *string { s := "example-proof-value"; return &s }(),
	}

	registerIssuerRequest := nodeV1alpha.RegisterIssuerRequest{
		Issuer: &issuer,
		Proof:  &proof,
	}

	// Call the client to generate metadata
	log.Default().Println("Registering issuer with request: ", &registerIssuerRequest)

	issuerId, err := s.issuerRepository.AddIssuer(vaultId, identityNodeAddress, idpConfig, issuer)
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

func (s *issuerService) GetIssuer(vaultId, issuerId string) (*coreV1alpha.Issuer, error) {
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

func GetIdpToken(clientId, clientSecret, issuerUrl string) (*oauth2.Token, error) {
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
