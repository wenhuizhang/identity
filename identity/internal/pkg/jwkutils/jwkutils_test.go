// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package jwkutils

import (
	"testing"
)

func TestGenerateAndValidateRSAKeys(t *testing.T) {
	algorithms := []string{"RS256", "RS384", "RS512"}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			jwk, err := GenerateJWK(alg, "sig", "")
			if err != nil {
				t.Fatalf("GenerateJWK failed for %s: %v", alg, err)
			}

			if err := ValidatePubKey(jwk); err != nil {
				t.Errorf("ValidatePubKey failed for %s: %v", alg, err)
			}
			if err := ValidatePrivKey(jwk); err != nil {
				t.Errorf("ValidatePrivKey failed for %s: %v", alg, err)
			}
		})
	}
}

func TestGenerateAndValidateAKPKeys(t *testing.T) {
	algorithms := []string{"ML-DSA-44", "ML-DSA-65", "ML-DSA-87"}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			jwk, err := GenerateJWK(alg, "sig", "")
			if err != nil {
				t.Fatalf("GenerateJWK failed for %s: %v", alg, err)
			}

			if err := ValidatePubKey(jwk); err != nil {
				t.Errorf("ValidatePubKey failed for %s: %v", alg, err)
			}
			if err := ValidatePrivKey(jwk); err != nil {
				t.Errorf("ValidatePrivKey failed for %s: %v", alg, err)
			}
		})
	}
}

func TestGenerateJWK_UnsupportedAlgorithm(t *testing.T) {
	_, err := GenerateJWK("unsupported-alg", "sig", "")
	if err == nil {
		t.Error("Expected error for unsupported algorithm, got nil")
	}
}
