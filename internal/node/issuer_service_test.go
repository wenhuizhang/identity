// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"testing"

	idtypes "github.com/agntcy/identity/internal/core/id/types"
	issuertesting "github.com/agntcy/identity/internal/core/issuer/testing"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	verificationtesting "github.com/agntcy/identity/internal/core/issuer/verification/testing"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/node"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/stretchr/testify/assert"
)

func TestRegisterIssuer_Should_Not_Register_Same_Issuer_Twice(t *testing.T) {
	t.Parallel()

	verficationSrv := verificationtesting.NewFakeVerifiedVerificationService()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	sut := node.NewIssuerService(issuerRepo, verficationSrv)
	pubKey, _ := generatePubKey()

	issuer := &issuertypes.Issuer{
		CommonName:   verificationtesting.ValidProofIssuer,
		Organization: "Some Org",
		PublicKey:    pubKey,
	}

	proof := &vctypes.Proof{
		Type:       "JWT",
		ProofValue: "",
	}

	// Register once
	_, err := sut.Register(context.Background(), issuer, proof)
	assert.NoError(t, err)

	// Attempt to register the same issuer again
	_, err = sut.Register(context.Background(), issuer, proof)
	assert.Error(t, err)
}

func generatePubKey() (*idtypes.Jwk, error) {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	pubkey, err := jwk.PublicRawKeyOf(pk)
	if err != nil {
		return nil, err
	}

	key, err := jwk.Import(pubkey)
	if err != nil {
		return nil, err
	}

	err = key.Set(jwk.AlgorithmKey, jwa.RS256())
	if err != nil {
		return nil, err
	}

	keyAsJson, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}

	var k idtypes.Jwk

	err = json.Unmarshal(keyAsJson, &k)
	if err != nil {
		return nil, err
	}

	return &k, nil
}
