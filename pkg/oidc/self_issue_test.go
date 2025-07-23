// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package oidc_test

import (
	"testing"

	"github.com/agntcy/identity/pkg/joseutil"
	"github.com/agntcy/identity/pkg/oidc"
	"github.com/stretchr/testify/assert"
)

func TestSelfIssueJWT_Should_Generate_Unique_JWTs(t *testing.T) {
	t.Parallel()

	jwk, err := joseutil.GenerateJWK("RS256", "sig", "my-id")
	if err != nil {
		t.Error(err)
	}

	tokens := make([]*string, 0)

	for idx := 0; idx < 10; idx++ {
		token, err := oidc.SelfIssueJWT("issuer", "sub", jwk)

		assert.NoError(t, err)

		tokens = append(tokens, &token)
	}

	for idx, token := range tokens {
		assert.NotContains(t, tokens[:idx], token)
		assert.NotContains(t, tokens[idx+1:], token)
	}
}
