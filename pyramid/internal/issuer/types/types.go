package types

import (
	idtypes "github.com/agntcy/pyramid/internal/core/id/types"
)

// A PyramID Issuer
type Issuer struct {
	// The organization of the issuer
	Organization string `json:"organization,omitempty"`

	// The sub organization of the issuer
	SubOrganization string `json:"sub_organization,omitempty"`

	// The common name of the issuer
	// Could be a FQDN or a FQDA
	CommonName string `json:"common_name,omitempty"`

	// The keys of the issuer in QJWK format
	// The public key is used to verify the signature of the different claims
	PublicKey *idtypes.Qjwk `json:"public_key,omitempty"`

	// The private key of the issuer in JWK format
	PrivateKey *idtypes.Qjwk `json:"private_key,omitempty"`
}
