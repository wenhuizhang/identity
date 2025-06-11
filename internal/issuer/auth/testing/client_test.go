// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package auth_test

import (
	"context"
	"testing"

	"github.com/agntcy/identity/internal/issuer/auth"
	issuertesting "github.com/agntcy/identity/internal/issuer/issuer/data/testing"
	idptypes "github.com/agntcy/identity/internal/issuer/types"
	vaulttesting "github.com/agntcy/identity/internal/issuer/vault/testing"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/stretchr/testify/assert"
)

func TestToken_Should_Issue_A_Self_Issued_Token(t *testing.T) {
	t.Parallel()

	authClient := auth.NewClient(
		issuertesting.NewFakeIssuerRepository(),
		oidc.NewAuthenticator(),
		vaulttesting.NewFakeVaultService(),
	)

	_, err := authClient.Token(
		context.Background(),
		"vaultId",
		"keyId",
		"issuerId",
		nil,
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, "token", "Expected a token to be issued, but got an empty string")
}

func TestToken_Should_Issue_A_JWT_Signed_Token(t *testing.T) {
	t.Parallel()

	authClient := auth.NewClient(
		issuertesting.NewFakeIssuerRepository(),
		oidc.NewAuthenticator(),
		vaulttesting.NewFakeVaultService(),
	)

	_, err := authClient.Token(
		context.Background(),
		"vaultId",
		"keyId",
		"issuerId",
		&idptypes.IdpConfig{
			ClientId:     "client-id",
			ClientSecret: "client-secret",
			IssuerUrl:    "https://example.com",
		},
	)

	assert.Error(t, err, "Expected an error when issuing a JWT signed token without a private key")
}
