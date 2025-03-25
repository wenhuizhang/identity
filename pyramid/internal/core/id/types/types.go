package types

// QJWK represents a Quantum JSON Web Key (JWK) with the following fields specific to NTRU algorithms.
// This could be an extension of the JWK type.
type Qjwk struct {
	// ALG represents the algorithm intended for use with the key.
	// Some example algorithms are "Falcon" family: "Falcon-512", "Falcon-1024"
	ALG string `json:"alg,omitempty"`

	// KTY represents the key type parameter.
	// It specifies the family of quantum algorithms used with the key,
	// such as "NTRU"
	KTY string `json:"kty,omitempty"`

	// Use represents the intended use of the key.
	// Some example values are "enc" and "sig".
	USE string `json:"use,omitempty"`

	// KID represents the key ID.
	// It is used to match a specific key.
	KID string `json:"kid,omitempty"`

	// h represents the public key.
	H string `json:"h,omitempty"`

	// f represents the private key.
	F string `json:"f,omitempty"`

	// fp represents the private key.
	FP string `json:"fp,omitempty"`

	// g represents the public key.
	G string `json:"g,omitempty"`
}

// QJWKS represents a set of Quantum JSON Web Keys (JWKs).
type Qjwks struct {
	// Keys represents the list of Quantum JSON Web Keys.
	Keys []Qjwk `json:"keys"`
}

// The Envelope Type of the Agent Passport.
// Multiple envelope types can be supported: JOSE, COSE etc.
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

// The Agent content representation following OASF schema
// Specs: https://schema.oasf.agntcy.org/schema/objects/agent
type AgentContent struct {
	// The OASF in json format
	Value string `json:"value,omitempty"`
}
