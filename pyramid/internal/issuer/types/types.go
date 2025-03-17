package types

import (
	didtypes "github.com/agntcy/pyramid/internal/core/did/types"
)

// A PyramID Issuer
type Issuer struct {
	// The common name of the issuer
	// Example: isuser.com
	CommonName string `json:"common_name,omitempty"`

	// The keys of the issuer of the DID in JWK format
	// The public key is used to verify the signature of the DID
	PublicKey *didtypes.JWK `json:"public_key,omitempty"`

	// The private key of the issuer of the DID in JWK format
	PrivateKey *didtypes.JWK `json:"private_key,omitempty"`
}
