package types

// JWK represents a JSON Web Key as per RFC7517 (https://tools.ietf.org/html/rfc7517)
// Note that this is a subset of the spec. There are a handful of properties that the
// spec allows for that are not represented here at the moment. This is because we
// only need a subset of the spec for our purposes.
type JWK struct {
	// ALG represents the algorithm intended for use with the key.
	ALG string `json:"alg,omitempty"`

	// KTY represents the key type parameter.
	// It specifies the family of cryptographic algorithms used with the key,
	// such as "RSA" or "EC" for elliptic curve keys.
	KTY string `json:"kty,omitempty"`

	// CRV represents the curve parameter for elliptic curve keys.
	// It specifies the cryptographic curve used with the key, such as "P-256" or "P-384".
	CRV string `json:"crv,omitempty"`

	// D represents the private key parameter.
	// This field is used to store the private key material for asymmetric keys.
	D string `json:"d,omitempty"`

	// X represents the x-coordinate for elliptic curve keys.
	// This field is part of the public key material for elliptic curve cryptography (ECC).
	X string `json:"x,omitempty"`

	// Y represents the y-coordinate for elliptic curve keys.
	// This field is part of the public key material for elliptic curve cryptography (ECC)
	Y string `json:"y,omitempty"`
}

type AgentPassportEnvelopeType int

const (
	// Unspecified Envelope Type.
	AGENT_PASSPORT_ENVELOPE_TYPE_UNSPECIFIED AgentPassportEnvelopeType = iota

	// JOSE Envelope Type.
	AGENT_PASSPORT_ENVELOPE_TYPE_JOSE
)

// AgentPassport represents an identity passport for an agent.
type AgentPassport struct {
	// EnvelopeType specifies the type of the envelope used to store the passport.
	EnvelopeType AgentPassportEnvelopeType `json:"envelope_type,omitempty"`

	// Value is the value of the passport.
	Value string `json:"value,omitempty"`
}

// The Agent definition following OASF schema
// Specs: https://schema.oasf.agntcy.org/schema/objects/agent
type AgentClaims struct {
	// The OASF in json format
	Value string `json:"value,omitempty"`
}
