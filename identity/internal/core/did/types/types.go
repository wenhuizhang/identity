package types

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
	PublicKeyJwk *Jwk `json:"public_key_jwk,omitempty"`
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

// JWK represents:
// - a JSON Web Key (JWK) with the respective fields specific to RSA algorithms.
// - a Quantum JSON Web Key (QJWK) with the respective fields specific to NTRU algorithms.
type Jwk struct {
	// ALG represents the algorithm intended for use with the key.
	// Some example algorithms are "Falcon" family: "Falcon-512", "Falcon-1024" for Quantum algorithms.
	// Some example algorithms are "RS256", "RS384", "RS512" for RSA algorithms.
	ALG string `json:"alg,omitempty"`

	// KTY represents the key type parameter.
	// It specifies the family of quantum algorithms used with the key,
	// such as "NTRU" or "RSA" for non quantum algorithms.
	KTY string `json:"kty,omitempty"`

	// Use represents the intended use of the key.
	// Some example values are "enc" and "sig".
	USE string `json:"use,omitempty"`

	// KID represents the key ID.
	// It is used to match a specific key.
	KID string `json:"kid,omitempty"`

	// The exponent for the RSA public key.
	E string `json:"e,omitempty"`

	// The modulus for the RSA public key.
	N string `json:"n,omitempty"`

	// The private exponent for the RSA private key.
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

	// The public key for the NTRU public key.
	H string `json:"h,omitempty"`

	// The polynomial for the NTRU private key.
	F string `json:"f,omitempty"`

	// The f inverse modulo p for the NTRU private key.
	FP string `json:"fp,omitempty"`

	// The polynomial for the NTRU private key.
	G string `json:"g,omitempty"`
}

// JWKS represents a set of JSON Web Keys (JWKs).
type Jwks struct {
	// Keys represents the list of JSON Web Keys.
	Keys []Jwk `json:"keys"`
}
