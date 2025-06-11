// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"fmt"

	issuerData "github.com/agntcy/identity/internal/issuer/issuer/data"
	idptypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/pkg/jwtutil"
	"github.com/agntcy/identity/internal/pkg/oidc"
)

type Client interface {
	Token(
		ctx context.Context,
		vaultId, keyId, issuerId string,
		idpConfig *idptypes.IdpConfig,
	) (string, error)
}

type client struct {
	issuerRepository issuerData.IssuerRepository
	auth             oidc.Authenticator
	vaultSrv         vault.VaultService
}

func NewClient(
	issuerRepository issuerData.IssuerRepository,
	auth oidc.Authenticator,
	vaultSrv vault.VaultService,
) Client {
	return &client{
		issuerRepository: issuerRepository,
		auth:             auth,
		vaultSrv:         vaultSrv,
	}
}

func (s *client) Token(
	ctx context.Context,
	vaultId, keyId, issuerId string,
	idpConfig *idptypes.IdpConfig,
) (string, error) {
	issuer, err := s.issuerRepository.GetIssuer(vaultId, keyId, issuerId)
	if err != nil {
		return "", err
	}

	var auth string

	if idpConfig == nil {
		prvKey, keyErr := s.vaultSrv.RetrievePrivKey(
			ctx,
			vaultId,
			keyId,
		)
		if keyErr != nil {
			return "", fmt.Errorf("error retrieving public key: %w", err)
		}

		// If no IdpConfig is provided, we generate a JWT auth using the issuer's private key.
		auth, err = jwtutil.Jwt(
			issuer.CommonName,
			issuer.ID,
			prvKey,
		)
	} else {
		auth, err = s.auth.Token(
			ctx,
			issuer.IdpConfig.IssuerUrl,
			issuer.IdpConfig.ClientId,
			issuer.IdpConfig.ClientSecret,
		)
	}

	return auth, err
}
