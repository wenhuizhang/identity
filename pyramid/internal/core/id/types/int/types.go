package types

import "google.golang.org/protobuf/types/known/structpb"

// VerificationMethod expresses verification methods, such as cryptographic
// public keys, which can be used to authenticate or authorize interactions
// with the DID subject or associated parties. For example,
// a cryptographic public key can be used as a verification method with
// respect to a digital signature; in such usage, it verifies that the
// signer could use the associated cryptographic private key.
//
// Specification Reference: https://www.w3.org/TR/did-core/#verification-methods
type VerificationMethod struct {
	// A unique id of the verification method.
	ID string `json:"id"`

	// references exactly one verification method type. In order to maximize global
	// interoperability, the verification method type SHOULD be registered in the
	// DID Specification Registries: https://www.w3.org/TR/did-spec-registries/
	Type string `json:"type"`

	// a value that conforms to the rules in DID Syntax: https://www.w3.org/TR/did-core/#did-syntax
	Controller string `json:"controller"`

	// specification reference: https://www.w3.org/TR/did-core/#dfn-publickeyjwk
	PublicKeyJwk *JWK `json:"public_key_jwk,omitempty"`
}

// Service is used in DID documents to express ways of communicating with
// the DID subject or associated entities.
// A service can be any type of service the DID subject wants to advertise.
//
// Specification Reference: https://www.w3.org/TR/did-core/#services
type Service struct {
	// Id is the value of the id property and MUST be a URI conforming to RFC3986.
	// A conforming producer MUST NOT produce multiple service entries with
	// the same id. A conforming consumer MUST produce an error if it detects
	// multiple service entries with the same id.
	ID string `json:"id"`

	// Type is an example of registered types which can be found
	// here: https://www.w3.org/TR/did-spec-registries/#service-types
	Type string `json:"type"`

	// ServiceEndpoint is a network address, such as an HTTP URL, at which services
	// operate on behalf of a DID subject.
	ServiceEndpoint []string `json:"service_endpoint"`
}

// Did provides a way to parse and handle Decentralized Identifier (DID) URIs
// according to the W3C DID Core specification (https://www.w3.org/TR/did-core/).
type Did struct {
	// URI represents the complete Decentralized Identifier (DID) URI.
	// Spec: https://www.w3.org/TR/did-core/#did-syntax
	URI string `json:"uri"`

	// Method specifies the DID method in the URI, which indicates the underlying
	// method-specific identifier scheme (e.g., jwk, dht, key, etc.).
	// Spec: https://www.w3.org/TR/did-core/#method-schemes
	Method string `json:"method"`

	// ID is the method-specific identifier in the DID URI.
	// Spec: https://www.w3.org/TR/did-core/#method-specific-id
	ID string `json:"id"`
}

// DidDocument represents a set of data describing the DID subject including mechanisms such as:
//   - cryptographic public keys - used to authenticate itself and prove
//     association with the DID
//   - services - means of communicating or interacting with the DID subject or
//     associated entities via one or more service endpoints.
//     Examples include discovery services, agent services,
//     social networking services, file storage services,
//     and verifiable credential repository services.
//
// A DID Document can be retrieved by resolving a DID URI.
type DidDocument struct {
	// The DID {ID}
	// A did could be represented as `did:agntcy:{ID}`
	// The metadata below is related as claims to the {ID}
	ID string `json:"id,omitempty"`

	// The node that was used to publish the document
	Node string `json:"node,omitempty"`

	// Controller defines an entity that is authorized to make changes to a DID document.
	// The process of authorizing a DID controller is defined by the DID method.
	// It can be a string or a list of strings.
	Controller []string `json:"controller,omitempty"`

	// VerificationMethod is a list of cryptographic public keys, which can be used to authenticate or authorize
	// interactions with the DID subject or associated parties.
	VerificationMethod []VerificationMethod `json:"verification_method,omitempty"`

	// Service expresses ways of communicating with the DID subject or associated entities.
	// A service can be any type of service the DID subject wants to advertise.
	// spec reference: https://www.w3.org/TR/did-core/#verification-methods
	Service []Service `json:"service,omitempty"`

	// AssertionMethod is used to specify how the DID subject is expected to express claims,
	// such as for the purposes of issuing a Verifiable Credential.
	AssertionMethod []string `json:"assertion_method,omitempty"`
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
	Data *structpb.Struct `json:"data,omitempty"`
}
