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

const (
	rsaBits256 = 2048
	rsaBits384 = 3072
	rsaBits512 = 4096
)

func GenerateJWK(alg, use, id string) (*types.Jwk, error) {
	if id == "" {
		id = uuid.NewString()
	}

	switch alg {
	case "RS256", "RS384", "RS512":
		return generateRSAJWK(alg, use, id)
	case "ML-DSA-44", "ML-DSA-65", "ML-DSA-87":
		return generateMLDSAJWK(alg, use, id)
	default:
		return nil, errors.New("unsupported algorithm")
	}
}

func generateRSAJWK(alg, use, id string) (*types.Jwk, error) {
	bits := map[string]int{
		"RS256": rsaBits256,
		"RS384": rsaBits384,
		"RS512": rsaBits512,
	}[alg]

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return &types.Jwk{
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
	}, nil
}

func generateMLDSAJWK(alg, use, id string) (*types.Jwk, error) {
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

	return &types.Jwk{
		KID:  id,
		ALG:  alg,
		KTY:  "AKP",
		USE:  use,
		PUB:  base64.RawURLEncoding.EncodeToString(publicKey),
		PRIV: base64.RawURLEncoding.EncodeToString(privKeyBytes),
		SEED: "",
	}, nil
}
