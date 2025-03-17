package types

import (
	"github.com/decentralized-identity/web5-go/dids/didcore"
)

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

// A PyramID Decentralized Identifier Document
type DidDocument struct {
	// The DID {ID}
	// A did could be represented as `did:agntcy:{ID}`
	// The metadata below is related as claims to the {ID}
	ID string `json:"id,omitempty"`

	// The node that was used to publish the document
	Node string `json:"node,omitempty"`

	// A DID Document is a JSON-LD document that
	// contains cryptographic information about the DID
	Content *didcore.Document `json:"content,omitempty"`
}
