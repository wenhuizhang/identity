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
	coretesting "github.com/agntcy/identity/internal/core/testing"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/node"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/stretchr/testify/assert"
)

func TestRegisterIssuerShouldNotRegisterSameIssuerTwice(t *testing.T) {
	verficationSrv := coretesting.NewFakeTruthyVerificationService()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	sut := node.NewIssuerService(issuerRepo, verficationSrv)

	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
		PublicKey:    generatePubKey(),
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

func generatePubKey() *idtypes.Jwk {
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil
	}
	pubkey, _ := jwk.PublicRawKeyOf(pk)
	key, _ := jwk.Import(pubkey)
	key.Set(jwk.AlgorithmKey, jwa.RS256())
	keyAsJson, _ := json.Marshal(key)

	var k idtypes.Jwk
	_ = json.Unmarshal(keyAsJson, &k)

	return &k
}
