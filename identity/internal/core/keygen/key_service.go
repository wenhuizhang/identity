package keygen

import (
	"context"

	"github.com/agntcy/identity/internal/core/id/types"
)

// KeyService defines methods for generating, saving, and retrieving JWKs.
type KeyService interface {
	// GenerateKey generates a new JWK.
	GenerateKey(ctx context.Context, alg string, use string, id string) (*types.Jwk, error)

	// SaveKey saves a JWK to the key storage. it supports local file, 1Password, and OS keychain.
	SaveKey(ctx context.Context, id string, jwk *types.Jwk) error

	// RetrieveKey retrieves a public JWK by its ID.
	RetrievePubKey(ctx context.Context, id string) (*types.Jwk, error)

	// RetrieveKey retrieves a private JWK by its ID.
	RetrievePrivKey(ctx context.Context, id string) (*types.Jwk, error)
}
