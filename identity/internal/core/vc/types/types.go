package types

// The Envelope Type of the Credential.
// Multiple envelope types can be supported: Embedded Proof, JOSE, COSE etc.
type CredentialEnvelopeType int

const (
	// Unspecified Envelope Type.
	CREDENTIAL_ENVELOPE_TYPE_UNSPECIFIED CredentialEnvelopeType = iota

	// Embedded Proof Envelope Type.
	CREDENTIAL_ENVELOPE_TYPE_EMBEDDED_PROOF

	// JOSE Envelope Type.
	CREDENTIAL_ENVELOPE_TYPE_JOSE
)

// The content of the Credential.
// Multiple content types can be supported: AgentBadge, etc.
type CredentialContentType int

const (
	// Unspecified Content Type.
	CREDENTIAL_CONTENT_TYPE_UNSPECIFIED CredentialContentType = iota

	// AgentBadge Content Type.
	// The Agent content representation following a defined schema
	// OASF: https://schema.oasf.agntcy.org/schema/objects/agent
	// Google A2A: https://github.com/google/A2A/blob/main/specification/json/a2a.json
	CREDENTIAL_CONTENT_TYPE_AGENT_BADGE
)

// EnvelopedCredential represents a Credential enveloped in a specific format.
type EnvelopedCredential struct {
	// EnvelopeType specifies the type of the envelope used to store the credential.
	EnvelopeType CredentialEnvelopeType `json:"envelope_type,omitempty"`

	// Value is the enveloped credential in the specified format.
	Value string `json:"value,omitempty"`
}

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

// A data integrity proof provides information about the proof mechanism,
// parameters required to verify that proof, and the proof value itself.
type Proof struct {
	// The type of the proof
	Type string `json:"type"`

	// The proof purpose
	ProofPurpose string `json:"proof_purpose"`

	// The proof value
	ProofValue string `json:"proof_value"`
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

	// https://w3id.org/security#proof
	Proof *Proof `json:"proof,omitempty"`
}

// DataModel represents the W3C Verifiable Presentation Data Model defined [here]
//
// [here]: https://www.w3.org/TR/vc-data-model/
type VerifiablePresentation struct {
	// https://www.w3.org/TR/vc-data-model/#contexts
	Context []string `json:"context"`

	// https://www.w3.org/TR/vc-data-model/#dfn-type
	Type []string `json:"type"`

	// https://www.w3.org/2018/credentials#verifiableCredential
	VerifiableCredential []VerifiableCredential `json:"verifiable_credential,omitempty"`

	// https://w3id.org/security#proof
	Proof *Proof `json:"proof,omitempty"`
}
