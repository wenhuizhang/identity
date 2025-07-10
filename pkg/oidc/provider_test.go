// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

//nolint:testpackage // getOAuthWellKnownURL is not exported
package oidc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOAuthWellKnownURL_Should_Construct_URL(t *testing.T) {
	t.Parallel()

	issuer := "https://sso-tenant.sso.duosecurity.com/oauth/MY-ID"

	u := getOAuthWellKnownURL(issuer)

	assert.Equal(t, "https://sso-tenant.sso.duosecurity.com/.well-known/oauth-authorization-server/oauth/MY-ID", u)
}
