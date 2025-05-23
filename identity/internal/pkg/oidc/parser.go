// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package oidc

import (
	"context"
	"time"

	identitycache "github.com/agntcy/identity/internal/pkg/cache"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/internal/pkg/httputil"
	"github.com/agntcy/identity/pkg/log"
	freecache "github.com/coocood/freecache"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	freecache_store "github.com/eko/gocache/store/freecache/v4"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jws"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type Claims struct {
	Issuer  string `json:"iss"`
	Subject string `json:"sub"`
}

type ProviderName int

const (
	UnknownProviderName ProviderName = iota
	OktaProviderName
	DuoProviderName
)

type ParsedJWT struct {
	Claims   *Claims
	Provider ProviderName
}

type providerMetadata struct {
	Issuer   string `json:"issuer"`
	TokenURL string `json:"token_endpoint"`
	JWKSURL  string `json:"jwks_uri"`
}

const defaultCacheSize = 10 * 1024 * 1024 // 10MB
const defaultCacheExpiration = 24         // 24 hours

type CachedJwks struct {
	Jwks string
}

// The Parser defines different methods for the PARSER standard
type Parser interface {
	ParseJwt(ctx context.Context, jwtString *string) (*ParsedJWT, error)
}

// The parser struct implements the Parser interface
type parser struct {
	jwksCache *cache.Cache[[]byte]
}

// NewParser creates a new instance of the Parser
func NewParser() Parser {
	jwksCache := cache.New[[]byte](
		freecache_store.NewFreecache(
			freecache.NewCache(defaultCacheSize),
			store.WithExpiration(defaultCacheExpiration*time.Second),
		))

	return &parser{
		jwksCache,
	}
}

// ParseJwt parses the JWT, validates the signature and returns the claims
func (p *parser) ParseJwt(ctx context.Context, jwtString *string) (*ParsedJWT, error) {
	// Check if the JWT string is empty
	if jwtString == nil || *jwtString == "" {
		return nil, errutil.Err(nil, "JWT string is empty")
	}

	// Get the issuer from the JWT string
	jwtToken, err := jwt.Parse([]byte(*jwtString), jwt.WithVerify(false), jwt.WithValidate(true))
	if err != nil {
		log.Error(err)

		return nil, errutil.Err(err, "failed to parse JWT")
	}

	issuer, ok := jwtToken.Issuer()
	if !ok {
		return nil, errutil.Err(nil, "failed to decode JWT: missing 'iss' claim")
	}

	provider, err := p.getProviderMetadata(ctx, issuer)
	if err != nil {
		return nil, err
	}

	// Get the JWKS from the issuer
	jwks, err := p.getJwks(ctx, provider)
	if err != nil {
		return nil, err
	}

	providerName, err := p.detectProviderName(ctx, provider)
	if err != nil {
		return nil, err
	}

	// Verify the JWT signature
	_, err = jws.Verify([]byte(*jwtString), jws.WithKeySet(jwks))
	if err != nil {
		return nil, err
	}

	// Get subject from the JWT string
	subject, ok := jwtToken.Subject()
	if !ok {
		return nil, errutil.Err(nil, "failed to decode JWT: missing 'sub' claim")
	}

	return &ParsedJWT{
		Claims: &Claims{
			Issuer:  issuer,
			Subject: subject,
		},
		Provider: providerName,
	}, nil
}

func (p *parser) getWellKnownURL(issuer string) string {
	return issuer + "/.well-known/openid-configuration"
}

func (p *parser) getProviderMetadata(ctx context.Context, issuer string) (*providerMetadata, error) {
	// Get the raw data from the issuer
	var metadata providerMetadata

	// Get the well-known URL from the issuer
	wellKnownURL := p.getWellKnownURL(issuer)
	log.Debug("Getting metadata from issuer:", issuer, " with URL:", wellKnownURL)

	// Get the metadata from the issuer
	err := httputil.GetJSON(ctx, wellKnownURL, &metadata)
	if err != nil {
		return nil, errutil.Err(err, "failed to get metadata from issuer")
	}

	log.Debug("Got metadata from issuer:", metadata)

	return &metadata, nil
}

func (p *parser) getJwks(ctx context.Context, provider *providerMetadata) (jwk.Set, error) {
	// Try to get the cached JWKS
	cachedEntry, found := identitycache.GetFromCache[CachedJwks](ctx, p.jwksCache, provider.Issuer)
	if found {
		return p.parseJwks(&cachedEntry.Jwks)
	}

	// Get the raw data from the JWKS URI
	var jwksString string

	// Get the JWKs
	err := httputil.GetWithRawBody(ctx, provider.JWKSURL, nil, &jwksString)
	if err != nil {
		return nil, errutil.Err(err, "failed to get JWKS from issuer")
	}

	jwks, err := p.parseJwks(&jwksString)
	if err != nil {
		return nil, errutil.Err(err, "failed to parse JWKS")
	}

	// Cache the JWKS
	err = identitycache.AddToCache(ctx, p.jwksCache, provider.Issuer, &CachedJwks{Jwks: jwksString})
	if err != nil {
		log.Warn(err)
	}

	return jwks, nil
}

func (p *parser) parseJwks(jwksString *string) (jwk.Set, error) {
	if jwksString == nil || *jwksString == "" {
		return nil, errutil.Err(nil, "JWKS string is empty")
	}

	jwks, err := jwk.Parse([]byte(*jwksString))
	if err != nil {
		return nil, errutil.Err(err, "failed to parse JWKS")
	}

	return jwks, nil
}
