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

type AlgorithmType int

const (
	// Unspecified Algorithm Type.
	ALGORITHM_TYPE_UNSPECIFIED AlgorithmType = iota

	// SHA-512 Algorithm Type.
	ALGORITHM_TYPE_SHA_512

	// SLH-DSA Algorithm Type.
	ALGORITHM_TYPE_SLH_DSA

	// Other Algorithm Type.
	ALGORITHM_TYPE_OTHER = 99
)

// The digest of the targeted content, conforming to the requirements. Retrieved
// content SHOULD be verified against this digest when consumed via untrusted
// sources. The digest property acts as a content identifier, enabling content
// addressability. It uniquely identifies content by taking a collision-resistant
// hash of the bytes. If the digest can be communicated in a secure manner, one
// can verify content from an insecure source by recalculating the digest
// independently, ensuring the content has not been modified. The value of
// the digest property is a string consisting of an algorithm portion and an
// encoded portion. The algorithm specifies the cryptographic hash function
// and encoding used for the digest; the encoded portion contains the encoded
// result of the hash function. A digest MUST be calculated for all properties
// except the digest itself which MUST be ignored during the calculation.
// The model SHOULD then be updated with the calculated digest.
type Digest struct {
	// The hash algorithm used to create the digital fingerprint, normalized
	// to the caption of algorithm_id. In the case of Other, it is defined by
	// the event source.
	Algorithm string `json:"algorithm,omitempty"`

	// The identifier of the normalized hash algorithm, which was used to create the digital fingerprint.
	AlgorithmID AlgorithmType `json:"algorithm_id,omitempty"`

	// The digital fingerprint value.
	Value string `json:"value,omitempty"`
}

// OASF Agent definition
// Specs: https://schema.oasf.agntcy.org/schema/objects/agent
type Agent struct {
	Digest *Digest `json:"digest,omitempty"`

	// Name of the agent.
	Name string `json:"name,omitempty"`

	// Version of the agent.
	Version string `json:"version,omitempty"`

	// List of agent’s authors in the form of `author-name <author-email>`.
	Authors []string `json:"authors,omitempty"`

	// Creation timestamp of the agent in the RFC3339 format.
	// Specs: https://www.rfc-editor.org/rfc/rfc3339.html
	CreatedAt string `json:"created_at,omitempty"`

	// Additional metadata associated with this agent.
	Annotations map[string]string `json:"annotations,omitempty"`

	// List of skills that this agent can perform.
	Skills []Skill `json:"skills,omitempty"`

	// List of source locators where this agent can be found or used from.
	Locators []Locator `json:"locators,omitempty"`

	// List of extensions that describe this agent and its capabilities
	// and constraints more in depth.
	Extensions []Extension `json:"extensions,omitempty"`
}

// OASF Skill definition
// Specs: https://schema.oasf.agntcy.org/schema/objects/skill
type Skill struct {
	// Schema/object version.
	Version string `json:"version,omitempty"`

	// UID of the category.
	CategoryUid string `json:"category_uid,omitempty"`

	// UID of the class.
	ClassUid string `json:"class_uid,omitempty"`

	// Additional metadata for this skill.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Optional human-readable name of the category.
	CategoryName *string `json:"category_name,omitempty"`

	// Optional human-readable name of the class.
	ClassName *string `json:"class_name,omitempty"`
}

// OASF Locator definition
// Specs: https://schema.oasf.agntcy.org/schema/objects/locator
type Locator struct {
	// Type of the locator. Can be custom or native LocatorType.
	Type string `json:"type,omitempty"`

	// Location URI where this source can be found/accessed.
	// Specs: https://datatracker.ietf.org/doc/html/rfc1738
	Url string `json:"url,omitempty"`

	// Metadata associated with this locator.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Size of the source in bytes pointed by the {url} property.
	Size *uint64 `json:"size,omitempty"`

	// Digest of the source pointed by the {url} property.
	// Specs: https://github.com/opencontainers/image-spec/blob/maindescriptor.md#digests
	Digest *string `json:"digest,omitempty"`
}

// OASF Extension definition
// Specs: https://schema.oasf.agntcy.org/schema/objects/extension
type Extension struct {
	// Name of the extension attached to an agent.
	Name string `json:"name,omitempty"`

	// Version of the extension attached to an agent.
	Version string `json:"version,omitempty"`

	// Metadata associated with this extension.
	Annotations map[string]string `json:"annotations,omitempty"`

	// Value of the data.
	// Data *structpb.Struct `json:"data,omitempty"`
}
