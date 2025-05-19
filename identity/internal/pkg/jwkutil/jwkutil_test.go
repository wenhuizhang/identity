// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package jwkutil_test

import (
	"testing"

	"github.com/agntcy/identity/internal/core/id/types"
	"github.com/agntcy/identity/internal/pkg/jwkutil"
)

func TestGenerateAndValidateKeys(t *testing.T) {
	t.Parallel()

	algorithms := []string{"RS256", "RS384", "RS512"}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			t.Parallel()

			jwk, err := jwkutil.GenerateJWK(alg, "sig", "")
			if err != nil {
				t.Fatalf("GenerateJWK failed for %s: %v", alg, err)
			}

			// Validate public key (should pass)
			if err := jwkutil.ValidatePubKey(jwk.PublicKey()); err != nil {
				t.Errorf("ValidatePubKey failed for %s: %v", alg, err)
			}
			// Validate private key (should pass)
			if err := jwkutil.ValidatePrivKey(jwk); err != nil {
				t.Errorf("ValidatePrivKey failed for %s: %v", alg, err)
			}
			// Public key with private fields should fail
			if err := jwkutil.ValidatePubKey(jwk); err == nil {
				t.Errorf("ValidatePubKey should fail if private fields are present for %s", alg)
			}
		})
	}
}

func TestValidatePubKey_NilOrEmptyFields(t *testing.T) {
	t.Parallel()

	// Test nil Jwk pointer
	var nilJwk *types.Jwk
	if err := jwkutil.ValidatePubKey(nilJwk); err == nil {
		t.Error("ValidatePubKey should fail if Jwk is nil")
	}

	// Test Jwk with missing required fields for RSA
	jwk := &types.Jwk{KTY: "RSA"}
	if err := jwkutil.ValidatePubKey(jwk); err == nil {
		t.Error("ValidatePubKey should fail if required RSA fields are missing")
	}
}

func TestGenerateJWK_UnsupportedAlgorithm(t *testing.T) {
	t.Parallel()

	_, err := jwkutil.GenerateJWK("unsupported-alg", "sig", "")
	if err == nil {
		t.Error("Expected error for unsupported algorithm, got nil")
	}
}
