package types

import (
	"github.com/decentralized-identity/web5-go/jwk"
)

// A PyramID Issuer
type Issuer struct {
	// The common name of the issuer
	// Example: isuser.com
	CommonName string `json:"commonName,omitempty"`

	// The keys of the issuer of the DID in JWK format
	// The public key is used to verify the signature of the DID
	PublicKey **jwk.JWK `json:"publicKey,omitempty"`

	// The private key of the issuer of the DID in JWK format
	PrivateKey *jwk.JWK `json:"privateKey,omitempty"`
}
