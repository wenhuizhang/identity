package types

// The Envelope Type of the Verifiable Credential.
// Multiple envelope types can be supported: JOSE, COSE etc.
type VerifiableCredentialEnvelopeType int

const (
	// Unspecified Envelope Type.
	VERIFIABLE_CREDENTIAL_ENVELOPE_TYPE_UNSPECIFIED VerifiableCredentialEnvelopeType = iota

	// JOSE Envelope Type.
	VERIFIABLE_CREDENTIAL_ENVELOPE_TYPE_JOSE
)

// The content of the Credential.
// Multiple content types can be supported: AgentPassport, etc.
type CredentialContentType int

const (
	// Unspecified Content Type.
	CREDENTIAL_CONTENT_TYPE_UNSPECIFIED CredentialContentType = iota

	// AgentPassport Content Type.
	// The Agent content representation following OASF schema
	// Specs: https://schema.oasf.agntcy.org/schema/objects/agent
	CREDENTIAL_CONTENT_TYPE_AGENT_PASSPORT
)

// CredentialContent represents the content of a Verifiable Credential.
type CredentialContent struct {
	// Type specifies the type of the content of the credential.
	Type CredentialContentType `json:"content_type,omitempty"`

	// The content representation in JSON-LD format.
	Content string `json:"content,omitempty"`
}

// CredentialSchema represents the credentialSchema property of a Verifiable Credential.
// more information can be found [here]
//
// [here]: https://www.w3.org/TR/vc-data-model-2.0/#data-schemas
type CredentialSchema struct {
	// Type specifies the type of the file
	Type string `json:"type"`

	// The URL identifying the schema file
	ID string `json:"id"`
}

// DataModel represents the W3C Verifiable Credential Data Model defined [here]
//
// [here]: https://www.w3.org/TR/vc-data-model/
type VerifiableCredential struct {
	// https://www.w3.org/TR/vc-data-model/#contexts
	Context []string `json:"context"`

	// https://www.w3.org/TR/vc-data-model/#dfn-type
	Type []string `json:"type"`

	// https://www.w3.org/TR/vc-data-model/#issuer
	Issuer string `json:"issuer"`

	// https://www.w3.org/TR/vc-data-model/#credential-subject
	CredentialSubject string `json:"credential_subject"`

	// https://www.w3.org/TR/vc-data-model/#identifiers
	ID string `json:"id,omitempty"`

	// https://www.w3.org/TR/vc-data-model/#issuance-date
	IssuanceDate string `json:"issuance_date"`

	// https://www.w3.org/TR/vc-data-model/#expiration
	ExpirationDate string `json:"expiration_date,omitempty"`

	// https://www.w3.org/TR/vc-data-model-2.0/#data-schemas
	CredentialSchema []CredentialSchema `json:"credential_schema,omitempty"`
}

// EnveloppedVerifiableCredential represents a Verifiable Credential envelopped in a specific format.
type EnveloppedVerifiableCredential struct {
	// EnvelopeType specifies the type of the envelope used to store the credential.
	EnvelopeType VerifiableCredentialEnvelopeType `json:"envelope_type,omitempty"`

	// Value is the envelopped credential in the specified format.
	Value string `json:"value,omitempty"`
}
