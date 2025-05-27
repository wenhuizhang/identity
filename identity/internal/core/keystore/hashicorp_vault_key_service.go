// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package keystore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"

	"github.com/agntcy/identity/internal/core/id/types"
	"github.com/hashicorp/vault/api"
)

type VaultKeyService struct {
	client      *api.Client
	mountPath   string
	keyBasePath string
}

type VaultStorageConfig struct {
	Address     string
	Token       string
	MountPath   string
	KeyBasePath string
	Namespace   string
}

func (s *VaultKeyService) SaveKey(ctx context.Context, id string, jwk *types.Jwk) error {
	if jwk == nil {
		return errors.New("jwk cannot be nil")
	}

	jsonData, err := json.Marshal(jwk)
	if err != nil {
		return fmt.Errorf("failed to marshal JWK: %w", err)
	}

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"jwk": string(jsonData),
		},
	}

	fullPath := buildKeyPath(s.mountPath, s.keyBasePath, id)

	_, err = s.client.Logical().WriteWithContext(ctx, fullPath, data)
	if err != nil {
		return fmt.Errorf("failed to write JWK to Vault: %w", err)
	}

	return nil
}

func (s *VaultKeyService) RetrievePubKey(ctx context.Context, id string) (*types.Jwk, error) {
	jwk, err := s.retrieveKey(ctx, id)
	if err != nil {
		return nil, err
	}

	return jwk.PublicKey(), nil
}

func (s *VaultKeyService) RetrievePrivKey(ctx context.Context, id string) (*types.Jwk, error) {
	return s.retrieveKey(ctx, id)
}

func (s *VaultKeyService) retrieveKey(ctx context.Context, id string) (*types.Jwk, error) {
	fullPath := buildKeyPath(s.mountPath, s.keyBasePath, id)

	secret, err := s.client.Logical().ReadWithContext(ctx, fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWK from Vault: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.New("key not found in Vault")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid data format in Vault")
	}

	jwkJSON, ok := data["jwk"].(string)
	if !ok {
		return nil, errors.New("JWK not found in Vault data")
	}

	var jwk types.Jwk
	if err := json.Unmarshal([]byte(jwkJSON), &jwk); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWK: %w", err)
	}

	return &jwk, nil
}

func buildKeyPath(mountPath, basePath, id string) string {
	return path.Join(mountPath, "data", basePath, id)
}
