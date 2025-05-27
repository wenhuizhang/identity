// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package keystore_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/agntcy/identity/internal/core/keystore"
	"github.com/agntcy/identity/internal/pkg/joseutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultKeyService(t *testing.T) {
	t.Parallel()

	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping Vault integration test. Set INTEGRATION_TESTS=true to run")
	}

	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")

	require.NotEmpty(t, vaultAddr, "VAULT_ADDR environment variable must be set")
	require.NotEmpty(t, vaultToken, "VAULT_TOKEN environment variable must be set")

	config := keystore.VaultStorageConfig{
		Address:     vaultAddr,
		Token:       vaultToken,
		MountPath:   "secret",
		KeyBasePath: "test-jwks",
		Namespace:   "eticloud/apps/pyramid",
	}

	service, err := keystore.NewKeyService(keystore.VaultStorage, config)
	require.NoError(t, err, "Failed to create key service")

	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	keyID := "test-vault-key-" + uniqueSuffix
	priv, err := joseutil.GenerateJWK("RS256", "sig", keyID)
	require.NoError(t, err, "GenerateJWK failed")

	ctx := context.Background()

	err = service.SaveKey(ctx, priv.KID, priv)
	assert.NoError(t, err, "SaveKey failed")

	pub, err := service.RetrievePubKey(ctx, priv.KID)
	if !assert.NoError(t, err, "RetrievePubKey failed") {
		t.FailNow()
	}

	assert.Empty(t, pub.D, "PublicJWK should not contain private D field")
	assert.Empty(t, pub.P, "PublicJWK should not contain private P field")
	assert.Empty(t, pub.Q, "PublicJWK should not contain private Q field")

	assert.Equal(t, priv.N, pub.N, "PublicJWK should contain correct N field")
	assert.Equal(t, priv.E, pub.E, "PublicJWK should contain correct E field")

	gotPriv, err := service.RetrievePrivKey(ctx, priv.KID)
	assert.NoError(t, err, "RetrievePrivKey failed")

	assert.Equal(t, priv.D, gotPriv.D, "Private key should contain correct D field")
	assert.Equal(t, priv.P, gotPriv.P, "Private key should contain correct P field")
	assert.Equal(t, priv.Q, gotPriv.Q, "Private key should contain correct Q field")

	_, err = service.RetrievePubKey(ctx, "non-existent")
	assert.Error(t, err, "Should error for non-existent key")
}
