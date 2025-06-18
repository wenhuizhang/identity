// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"encoding/json"
	"strings"
)

// JWK represents:
// - a JSON Web Key (JWK) with the respective fields specific to RSA algorithms.
// - a Quantum JSON Web Key (QJWK) with the respective fields specific to AKP algorithms.
type Jwk struct {
	// ALG represents the algorithm intended for use with the key.
	// Example algorithms for Post-Quantum ML-DSA family:
	// "ML-DSA-44", "ML-DSA-65", "ML-DSA-87".
	// Some example algorithms are "RS256", "RS384", "RS512" for RSA algorithms.
	ALG string `json:"alg,omitempty"`

	// KTY represents the key type parameter.
	// It specifies the family of quantum algorithms used with the key,
	// such as "AKP" for post quantum algorithms
	// or "RSA" for non quantum algorithms.
	KTY string `json:"kty,omitempty"`

	// Use represents the intended use of the key.
	// Some example values are "enc" and "sig".
	USE string `json:"use,omitempty"`

	// KID represents the key ID.
	// It is used to match a specific key.
	KID string `json:"kid,omitempty"`

	// The public key for the AKP kty.
	PUB string `json:"pub,omitempty"`

	// The private key for the AKP kty.
	PRIV string `json:"priv,omitempty"`

	// Seed used to derive keys for ML-DSA alg.
	SEED string `json:"seed,omitempty"`

	// The exponent for the RSA public key.
	E string `json:"e,omitempty"`

	// The modulus for the RSA public key.
	N string `json:"n,omitempty"`

	// The private exponent for the RSA kty.
	D string `json:"d,omitempty"`

	// The first prime factor for the RSA private key.
	P string `json:"p,omitempty"`

	// The second prime factor for the RSA private key.
	Q string `json:"q,omitempty"`

	// The first factor CRT exponent for the RSA private key.
	DP string `json:"dp,omitempty"`

	// The second factor CRT exponent for the RSA private key.
	DQ string `json:"dq,omitempty"`

	// The first CRT coefficient for the RSA private key.
	QI string `json:"qi,omitempty"`
}

// PublicKey returns a copy of the private Jwk containing only the public fields.
func (j *Jwk) PublicKey() *Jwk {
	pub := &Jwk{
		ALG: j.ALG,
		KTY: j.KTY,
		USE: j.USE,
		KID: j.KID,
	}

	switch strings.ToUpper(j.KTY) {
	case "RSA":
		pub.N = j.N
		pub.E = j.E
	case "AKP":
		pub.PUB = j.PUB
	}

	return pub
}

func (j *Jwk) ToJSON() []byte {
	raw, err := json.Marshal(j)
	if err != nil {
		return nil
	}

	return raw
}

// JWKS represents a set of JSON Web Keys (JWKs).
type Jwks struct {
	// Keys represents the list of JSON Web Keys.
	Keys []*Jwk `json:"keys"`
}

// Converts the Jwks to a JSON byte slice.
func (j *Jwks) Raw() []byte {
	rawJwks, err := json.Marshal(j)
	if err != nil {
		return nil
	}

	return rawJwks
}

// Convert a raw Jwks to string
func (j *Jwks) String() *string {
	rawJwks := j.Raw()
	if rawJwks == nil {
		return nil
	}

	stringJwks := strings.TrimSpace(string(rawJwks))

	return &stringJwks
}

// Converts a single Jwk to a Jwks object
func (j *Jwk) Jwks() *Jwks {
	return &Jwks{
		Keys: []*Jwk{j},
	}
}
