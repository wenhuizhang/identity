// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package oidc

import (
	"context"
	"fmt"

	"github.com/agntcy/identity/internal/pkg/httputil"
	"github.com/agntcy/identity/pkg/log"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Claims struct {
	Issuer  string `json:"iss"`
	Subject string `json:"sub"`
}

// The Parser defines different methods for the PARSER standard
type Parser interface {
	ParseJwt(ctx context.Context, jwtString *string) (*Claims, error)
}

// The parser struct implements the Parser interface
type parser struct {
}

// NewParser creates a new instance of the Parser
func NewParser() Parser {
	return &parser{}
}

// ParseJwt parses the JWT, validates the signature and returns the claims
func (p *parser) ParseJwt(ctx context.Context, jwtString *string) (*Claims, error) {
	// Check if the JWT string is empty
	if jwtString == nil || *jwtString == "" {
		return nil, fmt.Errorf("JWT string is empty")
	}

	// Get the issuer from the JWT string
	jwtToken, err := jwt.Parse([]byte(*jwtString))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT")
	}

	// Get the JWKS from the issuer
	issuer, ok := jwtToken.Issuer()
	if !ok {
		return nil, fmt.Errorf("failed to decode JWT: missing 'iss' claim")
	}

	// Get subject from the JWT string
	subject, ok := jwtToken.Subject()
	if !ok {
		return nil, fmt.Errorf("failed to decode JWT: missing 'sub' claim")
	}

	// Get the JWKS URI from the issuer
	jwksURI, err := p.getJwksURI(ctx, issuer)
	if err != nil {
		return nil, err
	}

	// Validate the signature and get the issuer
	log.Debug("Validating JWT signature", jwksURI)

	return &Claims{
		Issuer:  issuer,
		Subject: subject,
	}, nil
}

func (p *parser) getWellKnownURL(issuer string) string {
	return issuer + "/.well-known/openid-configuration"
}

func (p *parser) getJwksURI(ctx context.Context, issuer string) (*string, error) {
	// Get the raw data from the issuer
	var metadata map[string]string

	// Get the metadata from the issuer
	httputil.GetData(ctx, p.getWellKnownURL(issuer), metadata)
	if metadata == nil {
		return nil, fmt.Errorf("failed to get metadata from issuer")
	}

	// Get the JWKS URI from the metadata
	jwksURI, ok := metadata["jwks_uri"]
	if !ok {
		return nil, fmt.Errorf("failed to decode JWT: missing 'jwks_uri' from metadata")
	}

	return &jwksURI, nil
}
