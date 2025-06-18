// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package keystore

import (
	"context"

	"github.com/agntcy/identity/pkg/jwk"
)

// KeyService defines methods for generating, saving, and retrieving JWKs.
type KeyService interface {
	// SaveKey saves a JWK to the key storage. it supports local file, and hashicorp vault
	SaveKey(ctx context.Context, id string, jwk *jwk.Jwk) error

	// RetrieveKey retrieves a public JWK by its ID.
	RetrievePubKey(ctx context.Context, id string) (*jwk.Jwk, error)

	// RetrieveKey retrieves a private JWK by its ID.
	RetrievePrivKey(ctx context.Context, id string) (*jwk.Jwk, error)

	// DeleteKey deletes a JWK by its ID.
	DeleteKey(ctx context.Context, id string) error

	// ListKeys returns all available key IDs.
	ListKeys(ctx context.Context) ([]string, error)
}
