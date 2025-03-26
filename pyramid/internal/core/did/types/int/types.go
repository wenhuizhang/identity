package types

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
