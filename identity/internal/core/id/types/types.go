package types

// VerificationMethod expresses verification methods, such as cryptographic
// public keys, which can be used to authenticate or authorize interactions
// with the VCs associated with the ID. It is a part of the ID Document.
type VerificationMethod struct {
	// A unique id of the verification method.
	ID string `json:"id"`

	// The public key used for the verification method.
	PublicKeyJwk *Jwk `json:"public_key_jwk,omitempty"`
}

// IdDocument represents a set of data describing the ID including mechanisms such as:
//   - cryptographic public keys - used to authenticate itself and prove
//     association with the ID
//   - node - the node that was used to publish the document
//
// An ID Document can be retrieved by resolving an ID.
type IdDocument struct {
	// The ID
	// The metadata below is related as claims to the ID
	ID string `json:"id,omitempty"`

	// The node that was used to publish the document
	Node string `json:"node,omitempty"`

	// VerificationMethod is a list of cryptographic public keys, which can be used
	// to authenticate or authorize interactions with the VCs associated with the ID.
	VerificationMethod []VerificationMethod `json:"verification_method,omitempty"`
}

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

// JWKS represents a set of JSON Web Keys (JWKs).
type Jwks struct {
	// Keys represents the list of JSON Web Keys.
	Keys []Jwk `json:"keys"`
}
