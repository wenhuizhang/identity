// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"context"

	"github.com/agntcy/identity/internal/pkg/oidc"
)

type fakeParser struct {
	result *oidc.ParsedJWT
	err    error
}

func NewFakeParser(
	result *oidc.ParsedJWT,
	err error,
) oidc.Parser {
	return &fakeParser{
		result: result,
		err:    err,
	}
}

func (p *fakeParser) ParseJwt(
	ctx context.Context,
	jwtString *string,
	jwksString *string,
) (*oidc.ParsedJWT, error) {
	return p.result, p.err
}

func (p *fakeParser) GetClaims(
	ctx context.Context,
	jwtString *string,
) (*oidc.Claims, error) {
	return p.result.Claims, p.err
}
