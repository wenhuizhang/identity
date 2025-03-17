package types

// DidSubjectType defines the types of subjects that a Decentralized Identifier (DID) can represent.
// It categorizes the role or nature of the subject within a decentralized identity framework.
type DidSubjectType int

const (
	// Unspecified Function Type.
	DID_SUBJECT_TYPE_UNSPECIFIED DidSubjectType = iota

	// The DID subject is an agent
	DID_SUBJECT_TYPE_AGENT

	// The DID subject is an agent locator
	DID_SUBJECT_TYPE_AGENT_LOCATOR
)

// DidSubject represents a subject within a decentralized identity framework.
type DidSubject struct {
	// A local unique id of the subject.
	ID string `json:"id,omitempty"`

	// Type specifies the type of the subject, as defined by the DidSubjectType enum.
	// This indicates the role or nature of the subject in the decentralized identity system,
	// such as whether it is an agent or an agent locator.
	Type DidSubjectType `json:"type,omitempty"`
}

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

// A PyramID Decentralized Identifier Document represents a set of data describing the DID subject including mechanisms such as:
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

type Did struct{}
