package types

import "fmt"

// Represents the reason for an error, providing a unique
// constant value for the error.
type ErrorReason int

const (
	// ERROR_REASON_UNSPECIFIED indicates that no specific error reason
	// has been specified.
	ERROR_REASON_UNSPECIFIED ErrorReason = iota

	// The Agent ID is invalid or not found
	ERROR_REASON_INVALID_ID

	// The credential envelope type is invalid. For valid values refer to
	// the enum CredentialEnvelopeType.
	ERROR_REASON_INVALID_CREDENTIAL_ENVELOPE_TYPE

	// The credential envelope value format does not correspond to the format
	// specified in envelope_type.
	ERROR_REASON_INVALID_CREDENTIAL_ENVELOPE_VALUE_FORMAT

	// The credential content type is invalid. For valid values refer to
	// the enum CredentialContentType.
	ERROR_REASON_INVALID_CREDENTIAL_CONTENT_TYPE

	// The credential content format is not a valid JSON-LD.
	ERROR_REASON_INVALID_CREDENTIAL_CONTENT_FORMAT

	// The issuer contains one or more invalid fields.
	ERROR_REASON_INVALID_ISSUER

	// The Verifiable Credential is invalid, this can be related to either
	// invalid format or unable to verify the Data Integrity proof.
	ERROR_REASON_INVALID_VERIFIABLE_CREDENTIAL

	// Unable to find an Agent Passport for an ID.
	ERROR_REASON_AGENT_PASSPORT_NOT_FOUND
)

// Describes the cause of the error with structured details.
type ErrorInfo struct {
	// The reason of the error, as defined by the ErrorReason enum.
	// This is a constant unique value that helps identify the cause of
	// the error.
	Reason ErrorReason `json:"reason,omitempty"`

	// The message describing the error in a human-readable way. This
	// field gives additional details about the error.
	Message string `json:"message,omitempty"`
}

func (err ErrorInfo) Error() string {
	return fmt.Sprintf("%s (reason: %d)", err.Message, err.Reason)
}
