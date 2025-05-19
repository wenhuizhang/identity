// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package jwkutil

import (
	"testing"

	"github.com/agntcy/identity/internal/core/id/types"
)

func TestGenerateAndValidateKeys(t *testing.T) {
	algorithms := []string{"RS256", "RS384", "RS512"}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			jwk, err := GenerateJWK(alg, "sig", "")
			if err != nil {
				t.Fatalf("GenerateJWK failed for %s: %v", alg, err)
			}

			// Validate public key (should pass)
			if err := ValidatePubKey(jwk.PublicKey()); err != nil {
				t.Errorf("ValidatePubKey failed for %s: %v", alg, err)
			}
			// Validate private key (should pass)
			if err := ValidatePrivKey(jwk); err != nil {
				t.Errorf("ValidatePrivKey failed for %s: %v", alg, err)
			}
			// Public key with private fields should fail
			if err := ValidatePubKey(jwk); err == nil {
				t.Errorf("ValidatePubKey should fail if private fields are present for %s", alg)
			}
		})
	}
}

func TestValidatePubKey_NilOrEmptyFields(t *testing.T) {
	// Test nil Jwk pointer
	var nilJwk *types.Jwk
	if err := ValidatePubKey(nilJwk); err == nil {
		t.Error("ValidatePubKey should fail if Jwk is nil")
	}

	// Test Jwk with missing required fields for RSA
	jwk := &types.Jwk{KTY: "RSA"}
	if err := ValidatePubKey(jwk); err == nil {
		t.Error("ValidatePubKey should fail if required RSA fields are missing")
	}
}

func TestGenerateJWK_UnsupportedAlgorithm(t *testing.T) {
	_, err := GenerateJWK("unsupported-alg", "sig", "")
	if err == nil {
		t.Error("Expected error for unsupported algorithm, got nil")
	}
}
