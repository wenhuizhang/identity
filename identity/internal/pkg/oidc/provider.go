// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package oidc

import (
	"context"
	"net/http"
	"strings"

	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/internal/pkg/httputil"
)

func (p *parser) detectProviderName(ctx context.Context, provider *providerMetadata) (ProviderName, error) {
	resp, err := httputil.Get(ctx, provider.JWKSURL, nil)
	if err != nil {
		return UnknownProviderName, err
	}

	resp.Body.Close()

	if isOkta(resp) {
		return OktaProviderName, nil
	} else if isDuo(resp) {
		return DuoProviderName, nil
	}

	return UnknownProviderName, errutil.Err(nil, "unable to detect provider name")
}

func isOkta(resp *http.Response) bool {
	return resp != nil && resp.Header.Get("X-Okta-Request-Id") != ""
}

func isDuo(resp *http.Response) bool {
	return resp != nil &&
		(strings.HasPrefix(strings.ToLower(resp.Header.Get("Server")), "duo") ||
			strings.HasPrefix(resp.Request.Host, "duosecurity.com"))
}
