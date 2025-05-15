// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package jwkutil

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"math/big"

	"github.com/agntcy/identity/internal/core/id/types"
	"github.com/google/uuid"
	"github.com/open-quantum-safe/liboqs-go/oqs"
)

// GenerateJWK generates a new JWK with the specified algorithm, usage, and ID.
func GenerateJWK(alg string, use string, id string) (*types.Jwk, error) {
	// Generate a unique key ID if not provided
	if id == "" {
		id = uuid.NewString()
	}

	var jwk *types.Jwk

	switch alg {
	case "RS256", "RS384", "RS512":
		bits := map[string]int{
			"RS256": 2048,
			"RS384": 3072,
			"RS512": 4096,
		}[alg]

		privateKey, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return nil, err
		}

		// Populate the JWK fields for RSA
		jwk = &types.Jwk{
			KID: id,
			ALG: alg,
			KTY: "RSA",
			USE: use,
			N:   base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
			E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privateKey.PublicKey.E)).Bytes()),
			D:   base64.RawURLEncoding.EncodeToString(privateKey.D.Bytes()),
			P:   base64.RawURLEncoding.EncodeToString(privateKey.Primes[0].Bytes()),
			Q:   base64.RawURLEncoding.EncodeToString(privateKey.Primes[1].Bytes()),
			DP:  base64.RawURLEncoding.EncodeToString(privateKey.Precomputed.Dp.Bytes()),
			DQ:  base64.RawURLEncoding.EncodeToString(privateKey.Precomputed.Dq.Bytes()),
			QI:  base64.RawURLEncoding.EncodeToString(privateKey.Precomputed.Qinv.Bytes()),
		}

	case "ML-DSA-44", "ML-DSA-65", "ML-DSA-87":
		sig := oqs.Signature{}
		if err := sig.Init(alg, nil); err != nil {
			return nil, err
		}
		defer sig.Clean()

		publicKey, err := sig.GenerateKeyPair()
		if err != nil {
			return nil, err
		}

		privKeyBytes := sig.ExportSecretKey()

		// Populate the JWK fields for ML-DSA
		jwk = &types.Jwk{
			KID:  id,
			ALG:  alg,
			KTY:  "AKP",
			USE:  use,
			PUB:  base64.RawURLEncoding.EncodeToString(publicKey),
			PRIV: base64.RawURLEncoding.EncodeToString(privKeyBytes),
			SEED: "", // If a seed is used for key derivation, populate it here
		}

	default:
		return nil, errors.New("unsupported algorithm")
	}

	return jwk, nil
}
