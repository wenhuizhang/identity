// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

//go:generate stringer -type=CredentialEnvelopeType

package types

import (
	"fmt"
	"slices"
	"strings"
)

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

	// McpBadge Content Type.
	// The MCP content representation following a defined schema
	// The schema is defined in the MCP specification as the MCPServer type
	CREDENTIAL_CONTENT_TYPE_MCP_BADGE
)

func (t CredentialContentType) String() string {
	switch t {
	case CREDENTIAL_CONTENT_TYPE_AGENT_BADGE:
		return "AgentBadge"
	case CREDENTIAL_CONTENT_TYPE_MCP_BADGE:
		return "MCPServerBadge"
	default:
		return ""
	}
}

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
	Type CredentialContentType `json:"content_type,omitempty" protobuf:"bytes,1,opt,name=content_type"`

	// The content representation in JSON-LD format.
	Content map[string]any `json:"content,omitempty" protobuf:"google.protobuf.Struct,2,opt,name=content"`
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

const (
	// ProofTypeJWT is the proof type for JWT
	proofTypeJWT = "JWT,JWTToken,Jwk,JwkToken,JwtToken,JwkToken,Jwt"
)

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

func (p *Proof) IsJWT() bool {
	return slices.Contains(strings.Split(proofTypeJWT, ","), p.Type)
}

// DataModel represents the W3C Verifiable Credential Data Model defined [here]
//
// [here]: https://www.w3.org/TR/vc-data-model/
type VerifiableCredential struct {
	// https://www.w3.org/TR/vc-data-model/#contexts
	Context []string `json:"context" protobuf:"bytes,1,opt,name=context"`

	// https://www.w3.org/TR/vc-data-model/#dfn-type
	Type []string `json:"type" protobuf:"bytes,2,opt,name=type"`

	// https://www.w3.org/TR/vc-data-model/#issuer
	Issuer string `json:"issuer" protobuf:"bytes,3,opt,name=issuer"`

	// https://www.w3.org/TR/vc-data-model/#credential-subject
	CredentialSubject map[string]any `json:"credential_subject" protobuf:"google.protobuf.Struct,4,opt,name=content"`

	// https://www.w3.org/TR/vc-data-model/#identifiers
	ID string `json:"id,omitempty" protobuf:"bytes,5,opt,name=id"`

	// https://www.w3.org/TR/vc-data-model/#issuance-date
	IssuanceDate string `json:"issuance_date" protobuf:"bytes,6,opt,name=issuance_date"`

	// https://www.w3.org/TR/vc-data-model/#expiration
	ExpirationDate string `json:"expiration_date,omitempty" protobuf:"bytes,7,opt,name=expiration_date"`

	// https://www.w3.org/TR/vc-data-model-2.0/#data-schemas
	CredentialSchema []*CredentialSchema `json:"credential_schema,omitempty" protobuf:"bytes,8,opt,name=credential_schema"`

	// https://w3id.org/security#proof
	Proof *Proof `json:"proof,omitempty" protobuf:"bytes,9,opt,name=proof"`
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

// BadgeClaims represents the content of a Badge VC defined [here]
//
// [here]: https://spec.identity.agntcy.org/docs/vc/intro/
type BadgeClaims struct {
	// The ID as defined [here]
	//
	// [here]: https://www.w3.org/TR/vc-data-model/#credential-subject
	ID string `json:"id"`

	// The content of the badge
	Badge string `json:"badge"`
}

func (c *BadgeClaims) ToMap() map[string]any {
	return map[string]any{
		"id":    c.ID,
		"badge": c.Badge,
	}
}

func (c *BadgeClaims) FromMap(src map[string]any) error {
	if id, ok := src["id"]; ok {
		c.ID = id.(string)
	} else {
		return fmt.Errorf("invalid badge claim: missing Resolver Metadata ID")
	}

	if b, ok := src["badge"]; ok {
		c.Badge = b.(string)
	} else {
		return fmt.Errorf("invalid badge claim: missing badge content")
	}

	return nil
}
